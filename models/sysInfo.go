package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type CPUInfo struct {
	gorm.Model
	CpuPercentage float64
}

type CPUTemp struct {
	gorm.Model
	Temperature float64
}

type SwapMemoryStat struct {
	gorm.Model
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
}

type VirtualMemoryStat struct {
	gorm.Model
	// Total amount of RAM on this system
	Total uint64 `json:"total"`

	// RAM available for programs to allocate
	//
	// This value is computed from the kernel specific values.
	Available uint64 `json:"available"`

	// RAM used by programs
	//
	// This value is computed from the kernel specific values.
	Used uint64 `json:"used"`

	// Percentage of RAM used by programs
	//
	// This value is computed from the kernel specific values.
	UsedPercent float64 `json:"usedPercent"`


	// Linux specific numbers
	// https://www.centos.org/docs/5/html/5.1/Deployment_Guide/s2-proc-meminfo.html
	// https://www.kernel.org/doc/Documentation/filesystems/proc.txt
	// https://www.kernel.org/doc/Documentation/vm/overcommit-accounting
	//Buffers        uint64 `json:"buffers"`
	//Cached         uint64 `json:"cached"`
	//Writeback      uint64 `json:"writeback"`
	//Dirty          uint64 `json:"dirty"`
	//WritebackTmp   uint64 `json:"writebacktmp"`
	//Shared         uint64 `json:"shared"`
	//Slab           uint64 `json:"slab"`
	//SReclaimable   uint64 `json:"sreclaimable"`
	//PageTables     uint64 `json:"pagetables"`
	//SwapCached     uint64 `json:"swapcached"`
	//CommitLimit    uint64 `json:"commitlimit"`
	//CommittedAS    uint64 `json:"committedas"`
	//HighTotal      uint64 `json:"hightotal"`
	//HighFree       uint64 `json:"highfree"`
	//LowTotal       uint64 `json:"lowtotal"`
	//LowFree        uint64 `json:"lowfree"`
	//SwapTotal      uint64 `json:"swaptotal"`
	//SwapFree       uint64 `json:"swapfree"`
	//Mapped         uint64 `json:"mapped"`
	//VMallocTotal   uint64 `json:"vmalloctotal"`
	//VMallocUsed    uint64 `json:"vmallocused"`
	//VMallocChunk   uint64 `json:"vmallocchunk"`
	//HugePagesTotal uint64 `json:"hugepagestotal"`
	//HugePagesFree  uint64 `json:"hugepagesfree"`
	//HugePageSize   uint64 `json:"hugepagesize"`
}

type ConnectionStat struct {
	gorm.Model
	Fd     uint32  `json:"fd"`
	Family uint32  `json:"family"`
	Type   uint32  `json:"type"`
	LAddr  string  `gorm:"type:varchar(255)";json:"localaddr"` // type: Addr
	RAddr  string  `json:"remoteaddr"` // type: Addr
	Status string  `json:"status"`
	UIDs   string  `gorm:"type:longtext";json:"uids"` // type []int32
	Pid    int32   `json:"pid"`
}


type Addr struct {
	IP   string `json:"ip"`
	Port uint32 `json:"port"`
}

type NetIOCountersStat struct {
	gorm.Model
	Name        string `json:"name"`        // interface name
	BytesSent   uint64 `json:"bytesSent"`   // number of bytes sent
	BytesRecv   uint64 `json:"bytesRecv"`   // number of bytes received
	PacketsSent uint64 `json:"packetsSent"` // number of packets sent
	PacketsRecv uint64 `json:"packetsRecv"` // number of packets received
	Errin       uint64 `json:"errin"`       // total number of errors while receiving
	Errout      uint64 `json:"errout"`      // total number of errors while sending
	Dropin      uint64 `json:"dropin"`      // total number of incoming packets which were dropped
	Dropout     uint64 `json:"dropout"`     // total number of outgoing packets which were dropped (always 0 on OSX and BSD)
	Fifoin      uint64 `json:"fifoin"`      // total number of FIFO buffers errors while receiving
	Fifoout     uint64 `json:"fifoout"`     // total number of FIFO buffers errors while sending

}

type ProcessStat struct {
	gorm.Model
	Pid				int32		`json:"pid"`
	Name			string		`json:"name"`
	CmdLine			string		`json:"cmdLine"`
	Connections		string		`json:"connections"`
	CreateTime		time.Time	`json:"createTime"`
	ExePath			string		`json:"exePath"`
	Username		string		`json:"username"`
	IOCounter		string		`json:"ioCounter"`
	MemoryInfo		string		`json:"memoryInfo"`
	CWD				string		`json:"cwd"`
	CPUPercentage 	float64		`json:"cpuPercentage"`
	NetIOCounter	string		`json:"netIOCounter"`

}