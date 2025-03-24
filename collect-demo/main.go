package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fatih/structs"
	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"golang.org/x/sync/errgroup"

	"github.com/shirou/gopsutil/v4/cpu"
)

var (
	lastNetIOStatTimeStamp int64
	lastNetInfoMap         map[string]*NetInfo
)

func main() {
	// 连接InfluxDB
	client := connInflux()
	// 组织名称
	org := "lbw"
	// 数据库bucket名称
	bucket := "test"
	// 获取写API
	writeAPI := client.WriteAPIBlocking(org, bucket)

	for {
		// 创建一个错误组
		// 每秒执行一次
		// 并发执行三个函数
		<-time.Tick(time.Second)
		// 加载CPU信息
		g := errgroup.Group{}
		// 获取点
		g.Go(func() error {
			cpuInfo := loadCPUInfo()
			point, err := GetPoint("cpu_percent", cpuInfo)
			if err != nil {
				return err
			}

			log.Printf("write %s info success", cpuInfo.Type())

			// 加载内存信息
			return writeAPI.WritePoint(context.Background(), point)
			// 获取点
		})

		g.Go(func() error {
			memInfo := loadMemInfo()
			point, err := GetPoint("mem", memInfo)
			if err != nil {
				return err
			}

			log.Printf("write %s info success", memInfo.Type())

			// 加载网络信息
			return writeAPI.WritePoint(context.Background(), point)
			// 获取点
		})

		g.Go(func() error {
			netInfo := loadNetInfo()
			points := make([]*write.Point, 0, len(netInfo))
			for _, v := range netInfo {
				point, err := GetPoint("net", v)
				if err != nil {
					log.Printf("get %s info point failed, net name: %v, err: %v", v.Type(), v.NetName, err)
				}
				points = append(points, point)
			}

			log.Printf("write net info success")
			return writeAPI.WritePoint(context.Background(), points...)
		})

		g.Go(func() error {
			diskInfo := loadDiskInfo()
			points := make([]*write.Point, 0, len(diskInfo))
			for _, v := range diskInfo {
				point, err := GetPoint("disk", v)
				if err != nil {
					log.Printf("get %s info point failed, net name: %v, err: %v", v.Type(), v.MountPoint, err)
				}
				points = append(points, point)
			}

			log.Printf("write disk info success")
			return writeAPI.WritePoint(context.Background(), points...)
		})

		err := g.Wait()
		if err != nil {
			log.Println(err)
		}

	}
}

func GetPoint(measurement string, s Sys) (*write.Point, error) {
	var (
		tags map[string]string
	)

	fields := structs.Map(s)
	switch s.Type() {
	case CPU:
		tags = map[string]string{
			"cpu": "cpu",
		}
	case Mem:
		tags = map[string]string{
			"mem": "mem",
		}

	case Disk:
		diskInfo := s.(*DiskInfo)
		tags = map[string]string{
			"disk": diskInfo.MountPoint,
		}
		fields = map[string]interface{}{
			"total":               diskInfo.UsageStat.Total,
			"free":                diskInfo.UsageStat.Free,
			"used":                diskInfo.UsageStat.Used,
			"used_percent":        diskInfo.UsageStat.UsedPercent,
			"inodes_total":        diskInfo.UsageStat.InodesTotal,
			"inodes_used":         diskInfo.UsageStat.InodesUsed,
			"inodes_free":         diskInfo.UsageStat.InodesFree,
			"inodes_used_percent": diskInfo.UsageStat.InodesUsedPercent,
		}
	case Net:
		netInfo := s.(*NetInfo)
		tags = map[string]string{
			"net": netInfo.NetName,
		}
		fields = map[string]interface{}{
			"bytes_sent_rate":   netInfo.BytesSentRate,
			"bytes_recv_rate":   netInfo.BytesRecvRate,
			"packets_sent_rate": netInfo.PacketsSentRate,
			"packets_recv_rate": netInfo.PacketsRecvRate,
		}
	default:
		return nil, fmt.Errorf("type %s is not valid", s.Type())
	}

	return write.NewPoint(measurement, tags, fields, time.Now()), nil
}

func connInflux() influxdb.Client {
	return influxdb.NewClient("http://localhost:8086", "YuE8PZ5uxwXFa4rlmXZWFKOPdiPnmYaE_G3sm9Y5_iajb9oCTQWc8d7c34TorpA1oVuvuMWPyxn7Nk3ijJvIlw==")
}

func loadCPUInfo() *CPUInfo {
	percent, _ := cpu.Percent(time.Second, false)
	return &CPUInfo{Percent: percent[0]}
}

func loadMemInfo() *MemInfo {
	info, err := mem.VirtualMemory()
	if err != nil {
		log.Println(err)
		return nil
	}

	return &MemInfo{
		Total:       info.Total,
		Used:        info.Used,
		Available:   info.Available,
		UsedPercent: info.UsedPercent,
		Active:      info.Active,
		Inactive:    info.Inactive,
		Wired:       info.Wired,
	}
}

func loadNetInfo() map[string]*NetInfo {
	res := make(map[string]*NetInfo)

	currentTimeStamp := time.Now().Unix()

	netIOs, err := net.IOCounters(true)
	if err != nil {
		log.Printf("get net io counters failed, err:%v", err)
		return nil
	}

	for _, netIO := range netIOs {
		ioStat := new(NetInfo)

		ioStat.NetName = netIO.Name
		ioStat.BytesSent = netIO.BytesSent
		ioStat.BytesRecv = netIO.BytesRecv
		ioStat.PacketsSent = netIO.PacketsSent
		ioStat.PacketsRecv = netIO.PacketsRecv

		res[netIO.Name] = ioStat

		// 开始计算网卡相关速率
		if lastNetIOStatTimeStamp == 0 || lastNetInfoMap == nil ||
			lastNetInfoMap[netIO.Name] == nil {
			continue
		}

		interval := float64(currentTimeStamp - lastNetIOStatTimeStamp)
		ioStat.BytesSentRate = float64(ioStat.BytesSent-lastNetInfoMap[netIO.Name].BytesSent) / interval
		ioStat.BytesRecvRate = float64(ioStat.BytesRecv-lastNetInfoMap[netIO.Name].BytesRecv) / interval
		ioStat.PacketsSentRate = float64(ioStat.PacketsSent-lastNetInfoMap[netIO.Name].PacketsSent) / interval
		ioStat.PacketsRecvRate = float64(ioStat.PacketsRecv-lastNetInfoMap[netIO.Name].PacketsRecv) / interval

	}

	lastNetIOStatTimeStamp = currentTimeStamp
	lastNetInfoMap = res

	return res
}

func loadDiskInfo() map[string]*DiskInfo {
	res := make(map[string]*DiskInfo, 16)
	parts, err := disk.Partitions(true)
	if err != nil {
		log.Printf("get disk parts info failed, err:%v", err)
		return nil
	}

	for _, part := range parts {
		usageStat, err := disk.Usage(part.Mountpoint)
		if err != nil {
			log.Printf("get usage stat failed, err: %v", err)
		}

		diskInfo := &DiskInfo{
			MountPoint: part.Mountpoint,
			UsageStat:  usageStat,
		}

		res[part.Mountpoint] = diskInfo
	}

	return res
}
