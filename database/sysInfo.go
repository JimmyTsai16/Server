package database

import (
	"github.com/jimmy/server/models"
	"log"
	"time"
)

func (d *GormDatabase) SaveCpuInfo(cpuInfo *models.CPUInfo) {
	//cpu := models.CPUInfo{
	//	CpuPercentage: cpuPercentage,
	//}
	d.DB.Create(&cpuInfo)
}

func (d *GormDatabase) SaveCpuTemperature(cpuTemp *models.CPUTemp) {
	d.DB.Create(&cpuTemp)
}

func (d *GormDatabase) SaveSwapMemStat(memStat *models.SwapMemoryStat) {
	d.DB.Create(& memStat)
}

func (d *GormDatabase) SaveVirtualMemStat(memStat *models.VirtualMemoryStat) {
	d.DB.Create(&memStat)
}

func (d *GormDatabase) SaveConnectionStat(connStat *models.ConnectionStat) {
	d.DB.Create(&connStat)
}

func (d *GormDatabase) SaveNetIOCounterStat(iocStat *models.NetIOCountersStat) {
	d.DB.Create(&iocStat)
}

func (d *GormDatabase) SaveProcessStat(processStat *models.ProcessStat) {
	d.DB.Create(&processStat)
}

func (d *GormDatabase) GetSysInfoBetween(start time.Time, end time.Time, t interface{}) interface{} {
	switch t.(type) {
	case *models.CPUInfo:
		var cpuInfos []models.CPUInfo
		d.DB.Where("Created_at BETWEEN ? AND ?", start, end).Find(&cpuInfos)
		return cpuInfos

	case *models.CPUTemp:
		var cpuTemps []models.CPUTemp
		d.DB.Where("Created_at BETWEEN ? AND ?", start, end).Find(&cpuTemps)
		return cpuTemps

	case *models.VirtualMemoryStat:
		var VMSs []models.VirtualMemoryStat
		d.DB.Where("Created_at BETWEEN ? AND ?", start, end).Find(&VMSs)
		return VMSs

	case *models.SwapMemoryStat:
		var SMSs []models.SwapMemoryStat
		d.DB.Where("Created_at BETWEEN ? AND ?", start, end).Find(&SMSs)
		return SMSs

	case *models.ConnectionStat:
		var CSs []models.ConnectionStat
		d.DB.Where("Created_at BETWEEN ? AND ?", start, end).Find(&CSs)
		return CSs

	case *models.NetIOCountersStat:
		var NIOCSs []models.NetIOCountersStat
		d.DB.Where("Created_at BETWEEN ? AND ?", start, end).Find(&NIOCSs)
		return NIOCSs

	default:
		log.Println("Get sys info function Nothing to match")
		return nil
	}
}

func (d *GormDatabase) GetCpuInfoBetween(start time.Time, end time.Time) []models.CPUInfo {
	var cpuInfo []models.CPUInfo
	// fmt.Println(start, end)
	d.DB.Where("Created_at BETWEEN ? AND ?", start, end).Find(&cpuInfo)
	return cpuInfo
}
