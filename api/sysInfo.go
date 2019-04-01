package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jimmy/server/models"
	"github.com/jimmy/server/ws"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

type sysInfoDatabase interface {
	GetCpuInfoBetween(start time.Time, end time.Time) []models.CPUInfo
	GetSysInfoBetween(start time.Time, end time.Time, t interface{}) interface{}
}

const (
	cpuPercent		= "cpuPercent"
	cpuTemp			= "cpuTemp"
	swapMemInfo		= "swapMemInfo"
	virtualMemStat	= "virtualMemStat"
	processStat		= "processStat"
	netIOCounter	= "netIOCounter"
	connectionsStat	= "connectionsStat"
	hostInfo		= "hostInfo"
)

var InfoTypeEnum = map[string]int{
	cpuPercent:			0,
	cpuTemp:			1,
	swapMemInfo:		2,
	virtualMemStat:		3,
	processStat:		4,
	netIOCounter:		5,
	connectionsStat:	6,
	hostInfo:			7,
}

type SysInfoAPI struct {
	DB sysInfoDatabase
	SIR *sysInfoService
}

func NewSysInfoApi(db sysInfoDatabase) SysInfoAPI {
	sir := NewSysInfoService()
	return SysInfoAPI{DB: db, SIR: sir}
}

func (d *SysInfoAPI) GetCpuInfoBetween(ctx *gin.Context) {
	//start, _ := time.Parse(time.RFC3339, "2019-03-22T02:00:00+08:00")
	//end, _ := time.Parse(time.RFC3339, "2019-03-22T03:00:00+08:00")
	//fmt.Println(ctx.Param("startDate"))
	//fmt.Println(ctx.Param("endDate"))

	start, end := parseTime(ctx)

	x := d.DB.GetCpuInfoBetween(start, end)
	ctx.JSON(http.StatusOK, &x)
}

func (d *SysInfoAPI) GetHostInfo(ctx *gin.Context) {
	hostInfo, _ := host.Info()
	s, _ := host.SensorsTemperatures()
	fmt.Println(s)

	ctx.JSON(http.StatusOK, &hostInfo)
}

func(d *SysInfoAPI) GetSysInfo(ctx *gin.Context) {
	var info string
	start, end := parseTime(ctx)
	if d, ok := ctx.Params.Get("info"); ok {
		info = d
	}

	switch info {
	case cpuTemp:
		t := d.DB.GetSysInfoBetween(start, end, &models.CPUTemp{})
		cpuTemps, err := t.([]models.CPUTemp)
		if !err {
			log.Fatal("Type assertion if fail: ", err)
		}
		ctx.JSON(http.StatusOK, &cpuTemps)

	case cpuPercent:
		fmt.Println("Get cpu percentage")
		t := d.DB.GetSysInfoBetween(start, end, &models.CPUInfo{})
		cpuInfos, err := t.([]models.CPUInfo)
		if !err {
			log.Fatal("Type assertion if fail: ", err)
		}
		ctx.JSON(http.StatusOK, &cpuInfos)

	}

}



func (d *SysInfoAPI) GetSysInfoWS(ctx *gin.Context) {
	var services []string
	conn := ws.NewWS(ctx) // Upgrade the connection from GET to WebSocket.
	defer conn.Close()

	c := &client{
		conn: conn,
		sendingMessage: make(chan *message, 8),
		service: make([]bool, len(InfoTypeEnum)),
	}

	defer func(){
		fmt.Println("Client offline")
		c.closed = true
		for _, s := range services {
			d.SIR.RWMux.Lock()
			d.SIR.serviceClientsCount[s] -= 1
			d.SIR.RWMux.Unlock()
		}
	}()


	go c.ClientRoutine()
	d.SIR.AppendClient(c)

	for {
		var recv map[string]string
		var subscribeType string
		if err := conn.ReadJSON(&recv); err != nil {
			log.Println("ws readjson fail: ", err)
			return
		} else {
			fmt.Println(recv)
			d.SIR.RWMux.Lock()
			subscribeType = recv["subscribe"]
			closed, ok := d.SIR.serviceClosed[subscribeType]
			d.SIR.RWMux.Unlock()
			if !ok {
				fmt.Println("Wrong parameter: ", subscribeType)
				continue
			}
			if closed {
				d.SIR.RWMux.Lock()
				d.SIR.serviceClosed[subscribeType] = false
				d.SIR.RWMux.Unlock()
				fmt.Println(subscribeType)
				d.SIR.StartService(subscribeType)
			}
			d.SIR.RWMux.Lock()
			d.SIR.serviceClientsCount[subscribeType] += 1
			d.SIR.RWMux.Unlock()
			services = append(services, subscribeType)

			c.service[InfoTypeEnum[subscribeType]] = true

		}
	}

}

type sysInfoService struct {
	serviceClosed  map[string]bool
	Clients        []*client
	sendingMessage chan *message
	RWMux          sync.RWMutex
	serviceClientsCount map[string]uint32
}

type client struct {
	conn 			*websocket.Conn
	sendingMessage	chan *message
	service			[]bool
	closed			bool
	RWMux 			sync.RWMutex
}

type message struct {
	InfoType	string		`json:"infoType"`
	Data		interface{}	`json:"data"`
}

func NewSysInfoService() *sysInfoService {
	var s sysInfoService
	s.serviceClosed = make(map[string]bool)
	s.serviceClientsCount = make(map[string]uint32)
	// s.Clients = make(map[string][]*client)
	s.sendingMessage = make(chan *message, 16)

	s.serviceClosed[cpuPercent] = true
	s.serviceClosed[cpuTemp] = true
	s.serviceClosed[swapMemInfo] = true
	s.serviceClosed[virtualMemStat] = true
	s.serviceClosed[processStat] = true
	s.serviceClosed[netIOCounter] = true
	s.serviceClosed[connectionsStat] = true
	s.serviceClosed[hostInfo] = true

	for i := range InfoTypeEnum {
		s.serviceClientsCount[i] = 0
	}

	go func() {
		for {
			select {
			case <- time.After(time.Second * 5):
				for i, d := range s.serviceClientsCount {
					if d == 0 {
						s.serviceClosed[i] = true
					}
				}
			}
		}
	}()

	// Use channel to avoid websocket concurrent writing.
	go func() {
		for {
			select {
			case d := <-s.sendingMessage:

				for i := 0; i < len(s.Clients); i++ {
					c := s.Clients[i]
					if c.closed {
						c.Close()
						s.RWMux.Lock()
						s.Clients = append(s.Clients[:i], s.Clients[i + 1:]...)
						s.RWMux.Unlock()
						i--
						continue
					}
					if c.service[InfoTypeEnum[d.InfoType]] {
						c.sendingMessage <- d
					}
				}
			}
		}
	}()

	return &s
}

func (s *sysInfoService)StartService(subscribe string) {
	s.RWMux.Lock()
	defer s.RWMux.Unlock()
	switch subscribe {
	case cpuPercent:
		if closed, ok := s.serviceClosed[cpuPercent]; ok {
			if closed {
				log.Println(cpuPercent, " Service has running.")
				return
			}
		} else {
			log.Println(cpuPercent, " Service is not initialized.")
			return
		}
		go s.SaveCpuPercentRoutine(1)

	case cpuTemp:
		go s.SaveTemperatureRoutine(1)

	case swapMemInfo:
		go s.SaveSwapMemInfoRoutine(1)

	case virtualMemStat:
		go s.SaveVirtualMemStatRoutine(1)

	case connectionsStat:
		go s.SaveConnectionStatRoutine(1)

	case netIOCounter:
		go s.SaveNetIOCounterStatRoutine(1)

	case processStat:
		go s.SaveProcessStatRoutine(5)
	}
}

func (s *sysInfoService) AppendClient(c *client){
	s.RWMux.Lock()
	defer s.RWMux.Unlock()

	s.Clients = append(s.Clients, c)
}

func (c *client) ClientRoutine() {
	for {

		d, ok := <- c.sendingMessage
		if !ok {
			return
		}


		if err := c.conn.WriteJSON(d); err != nil {
			// log.Println("ws writeJson error: ", err)
			return
		}
		//select {
		//case d := <-c.sendingMessage:
		//	if err := c.conn.WriteJSON(d); err != nil {
		//		// log.Println("ws writeJson error: ", err)
		//		c.Close()
		//		return
		//	}
		//default:
		//
		//}
	}
}

func (c *client) Close(){
	close(c.sendingMessage)
}


func (s *sysInfoService) SaveCpuPercentRoutine(sec uint32) {
	defer s.Close(cpuPercent)
	for {

		if s.serviceClosed[cpuPercent] {
			log.Println(cpuPercent, "Routine have closed")
			return
		}

		var sum float64
		for i := uint32(0); i < sec; i++ {
			p, err := cpu.Percent(time.Second, false)
			if err != nil {
				log.Println("Get cpu information fail: ", err)
				return
			}
			sum += p[0]
		}
		mean := sum / float64(sec)
		s.sendingMessage <- &message{
			InfoType: cpuPercent,
			Data:     fmt.Sprintf("%f", mean),
		}

	}
}

func (s *sysInfoService)SaveTemperatureRoutine(sec uint32) {
	defer s.Close(cpuTemp)
	for {
		if s.serviceClosed[cpuTemp] {
			log.Println(cpuTemp, "Routine have closed")
			return
		}

		coreVolts := exec.Command("vcgencmd", "measure_temp")

		stdout, _ := coreVolts.StdoutPipe()
		if err := coreVolts.Start(); err != nil {
			log.Println("Command exec error: ", err)
			return
		}

		var buf []byte
		buf = make([]byte, 16)
		if _, err := stdout.Read(buf); err != nil {
			log.Println("Stdout read error: ", err)
		}

		s.sendingMessage <- &message{
			InfoType: cpuTemp,
			Data: string(buf[5:9]),
		}

		time.Sleep(time.Second * time.Duration(sec))
	}
}

func (s *sysInfoService)SaveSwapMemInfoRoutine(sec uint64) {
	defer s.Close(swapMemInfo)
	for {

		if s.serviceClosed[swapMemInfo] {
			log.Println(swapMemInfo, "Routine have closed")
			return
		}

		if swapMemStat, err := mem.SwapMemory(); err != nil {
			log.Println("Get swap memory stat fail: ", err)
			return
		} else {
			s.sendingMessage <- &message{
				InfoType: swapMemInfo,
				Data: models.SwapMemoryStat{
					Total: swapMemStat.Total,
					Used: swapMemStat.Used,
					Free: swapMemStat.Free,
					UsedPercent: swapMemStat.UsedPercent,
				},
			}
		}
		time.Sleep(time.Second * time.Duration(sec))
	}
}

func (s *sysInfoService)SaveVirtualMemStatRoutine(sec uint64) {
	defer s.Close(virtualMemStat)
	for {

		if s.serviceClosed[virtualMemStat] {
			log.Println(virtualMemStat, "Routine have closed")
			return
		}

		if VMS, err := mem.VirtualMemory(); err != nil {
			log.Println("Get virtual memory stat fail: ", err)
			return
		} else {
			s.sendingMessage <- &message{
				InfoType: virtualMemStat,
				Data: models.VirtualMemoryStat{
					Total:       VMS.Total,
					Used:        VMS.Used,
					Available:   VMS.Available,
					UsedPercent: VMS.UsedPercent,
				},
			}
		}
		time.Sleep(time.Second * time.Duration(sec))
	}
}

func (s *sysInfoService)SaveConnectionStatRoutine(sec uint64) {
	defer s.Close(connectionsStat)
	for {

		if s.serviceClosed[connectionsStat] {
			log.Println(connectionsStat, "Routine have closed")
			return
		}

		if cs, err := net.Connections("all"); err != nil {
			log.Println("Get connections fail: ", err)
			return
		} else {
			s.sendingMessage <- &message{
			InfoType: connectionsStat,
			Data: cs,
		}
		}



		time.Sleep(time.Second * time.Duration(sec))
	}
}

func (s *sysInfoService)SaveNetIOCounterStatRoutine(sec uint64 ) {
	defer s.Close(netIOCounter)
	for {

		if s.serviceClosed[netIOCounter] {
			log.Println(netIOCounter, "Routine have closed")
			return
		}

		iocStart, err := net.IOCounters(true)
		if err != nil {
			log.Println("Get net IOCounter start fail: ", err)
			return
		}

		time.Sleep(time.Second * time.Duration(sec))

		iocEnd, err := net.IOCounters(true)
		if err != nil {
			log.Println("Get net IOCounter end fail: ", err)
		}
		if len(iocStart) != len(iocEnd) {
			log.Println("Number of net IOCounter both pic interface isn't match")
		}

		var diffNetIOC []models.NetIOCountersStat

		for i := range iocStart {
			diffNetIOC = append(diffNetIOC, models.NetIOCountersStat{
				Name:        iocStart[i].Name,
				BytesSent:   iocEnd[i].BytesSent - iocStart[i].BytesSent,
				BytesRecv:   iocEnd[i].BytesRecv - iocStart[i].BytesRecv,
				PacketsSent: iocEnd[i].PacketsSent - iocStart[i].PacketsSent,
				PacketsRecv: iocEnd[i].PacketsRecv - iocStart[i].PacketsRecv,
				Errin:       iocEnd[i].Errin - iocStart[i].Errin,
				Errout:      iocEnd[i].Errout - iocStart[i].Errout,
				Dropin:      iocEnd[i].Dropin - iocStart[i].Dropin,
				Dropout:     iocEnd[i].Dropout - iocStart[i].Dropout,
				Fifoin:      iocEnd[i].Fifoin - iocStart[i].Fifoin,
				Fifoout:     iocEnd[i].Fifoout - iocStart[i].Fifoout,
			})
		}
		s.sendingMessage <- &message{
			InfoType: netIOCounter,
			Data: diffNetIOC,
		}
	}
}

func (s *sysInfoService)SaveProcessStatRoutine(sec uint64) {
	defer s.Close(processStat)
	for {

		if s.serviceClosed[processStat] {
			log.Println(processStat, "Routine have closed")
			return
		}

		// numCpu, _ := cpu.Counts(false)

		ps, err := process.Processes()
		if err != nil {
			log.Println("Get processes pid is fail: ", err)
			return
		}

		var pStats []models.ProcessStat
		for _, p := range ps {
			name, _ := p.Name()
			cmdLine, _ := p.Cmdline()
			connections, _ := p.Connections()
			connectionsString, _ := json.Marshal(&connections)
			createTime, _ := p.CreateTime()
			exePath, _ := p.Exe()
			username, _ := p.Username()
			ioCounters, _ := p.IOCounters()
			ioCountersString, _ := json.Marshal(ioCounters)
			memoryInfo, _ := p.MemoryInfo()
			memoryInfoString, _ := json.Marshal(&memoryInfo)
			cwd, _ := p.Cwd()
			netIOCounters, _ := p.NetIOCounters(true)
			netIOCountersString, _ := json.Marshal(netIOCounters)
			cpuPercentage, _ := p.CPUPercent()
			//cpuPercentage, _ := p.Percent(time.Millisecond * 100)
			//if cpuPercentage >= 0.005 {
			//	cpuPercentage = cpuPercentage / float64(numCpu)
			//}

			pStats = append(pStats, models.ProcessStat{
				Pid: p.Pid,
				Name: name,
				CmdLine: cmdLine,
				Connections: string(connectionsString),
				CreateTime: time.Unix(createTime / 1000.0, 0),
				ExePath: exePath,
				Username: username,
				IOCounter: string(ioCountersString),
				MemoryInfo: string(memoryInfoString),
				CWD: cwd,
				NetIOCounter: string(netIOCountersString),
				CPUPercentage: cpuPercentage,
			})
		}

		s.sendingMessage <- &message{
			InfoType: processStat,
			Data: pStats,
		}
		time.Sleep(time.Second * time.Duration(sec))
	}
}

func (s *sysInfoService) Close(serviceType string) {
	s.RWMux.Lock()
	defer s.RWMux.Unlock()
	s.serviceClosed[serviceType] = true
}

func (s *sysInfoService) SaveHostInfoRoutine(sec uint32) {
	defer s.Close(hostInfo)
	for {

		if s.serviceClosed[hostInfo] {
			log.Println(hostInfo, "Routine have closed")
			return
		}

		var sum float64
		for i := uint32(0); i < sec; i++ {
			p, err := cpu.Percent(time.Second, false)
			if err != nil {
				log.Println("Get cpu information fail: ", err)
				return
			}
			sum += p[0]
		}
		mean := sum / float64(sec)
		s.sendingMessage <- &message{
			InfoType: hostInfo,
			Data:     fmt.Sprintf("%f", mean),
		}

	}
}

func parseTime(ctx *gin.Context) (start time.Time, end time.Time) {
	// var start, end time.Time

	if s, ok := ctx.Params.Get("startDate");ok {
		t, _ := strconv.Atoi(s)
		t = t / 1000
		start = time.Unix(int64(t), 0)
		// fmt.Println(start)
	}else{
		log.Println("Get startDate Failed.")
	}

	if e, ok := ctx.Params.Get("endDate");ok {
		t, _ := strconv.Atoi(e)
		t = t / 1000
		end = time.Unix(int64(t), 0)
		// fmt.Println(end)
	}else{
		log.Println("Get endDate Failed.")
	}

	return
}