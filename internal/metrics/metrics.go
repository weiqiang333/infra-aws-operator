// Package metrics author: weiqiang; date: 2022-12
package metrics

import (
	"fmt"
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"infra-aws-operator/internal/databases/infrastructure"
	"infra-aws-operator/pkg/aws/connect"
	"infra-aws-operator/pkg/aws/get_lightsail"
	"infra-aws-operator/pkg/utils/date"
)

const (
	namespace = "aws"
)

type Exporter struct {
	DBC *gorm.DB
	// BillCost
	BillCostServiceByDay     *prometheus.GaugeVec
	BillCostServiceByMonthly *prometheus.GaugeVec
	// LightsailInstances
	LightsailInstances                            *prometheus.GaugeVec
	LightsailInstancesGbPerMonthAllocatedTransfer *prometheus.GaugeVec
	LightsailInstancesGbMonthNetworkIn            *prometheus.GaugeVec
	LightsailInstancesGbMonthNetworkOut           *prometheus.GaugeVec
	LightsailInstancesGbUseMonthNetwork           *prometheus.GaugeVec
	LightsailInstancesGbRemainMonthNetwork        *prometheus.GaugeVec
}

func NewExporter(dbC *gorm.DB) *Exporter {
	return &Exporter{
		DBC: dbC,
		BillCostServiceByDay: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "cost",
			Name:      "billcost_service_by_day",
			Help:      "按天计费服务",
		}, []string{"date_start", "date_end", "service_name"}),
		BillCostServiceByMonthly: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "cost",
			Name:      "billcost_service_by_month",
			Help:      "按月计费服务",
		}, []string{"date_start", "date_end", "service_name"}),
		LightsailInstances: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "lightsail",
			Name:      "instances",
			Help:      "光帆实例信息状态",
		}, []string{"name", "regions", "region_name", "region_name_abbreviation", "blueprint_name", "created_at", "cpu_count", "ram_size_in_gb", "private_ip_address", "public_ip_address", "gb_per_month_allocated_transfer"}),
		LightsailInstancesGbPerMonthAllocatedTransfer: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "lightsail",
			Name:      "gb_per_month_allocated_transfer",
			Help:      "光帆实例流量包免费流量 GB",
		}, []string{"name", "regions", "region_name", "region_name_abbreviation", "public_ip_address"}),
		LightsailInstancesGbMonthNetworkIn: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "lightsail",
			Name:      "gb_month_network_in",
			Help:      "光帆实例当月已使用传入流量 GB",
		}, []string{"name", "regions", "region_name", "region_name_abbreviation", "public_ip_address"}),
		LightsailInstancesGbMonthNetworkOut: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "lightsail",
			Name:      "gb_month_network_out",
			Help:      "光帆实例当月已使用传出流量 GB",
		}, []string{"name", "regions", "region_name", "region_name_abbreviation", "public_ip_address"}),
		LightsailInstancesGbUseMonthNetwork: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "lightsail",
			Name:      "gb_use_month_network",
			Help:      "光帆实例当月流量使用 sum GB",
		}, []string{"name", "regions", "region_name", "region_name_abbreviation", "public_ip_address"}),
		LightsailInstancesGbRemainMonthNetwork: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "lightsail",
			Name:      "gb_remain_month_network",
			Help:      "光帆实例当月流量套餐剩余 sum GB",
		}, []string{"name", "regions", "region_name", "region_name_abbreviation", "public_ip_address"}),
	}
}
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.BillCostServiceByDay.Describe(ch)
	e.BillCostServiceByMonthly.Describe(ch)
	e.LightsailInstances.Describe(ch)
	e.LightsailInstancesGbPerMonthAllocatedTransfer.Describe(ch)
	e.LightsailInstancesGbMonthNetworkIn.Describe(ch)
	e.LightsailInstancesGbMonthNetworkOut.Describe(ch)
	e.LightsailInstancesGbUseMonthNetwork.Describe(ch)
	e.LightsailInstancesGbRemainMonthNetwork.Describe(ch)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.BillCostServiceByDay.Reset()
	e.BillCostServiceByMonthly.Reset()
	e.LightsailInstances.Reset()
	e.LightsailInstancesGbPerMonthAllocatedTransfer.Reset()
	e.LightsailInstancesGbMonthNetworkIn.Reset()
	e.LightsailInstancesGbMonthNetworkOut.Reset()
	e.LightsailInstancesGbUseMonthNetwork.Reset()
	e.LightsailInstancesGbRemainMonthNetwork.Reset()

	for _, res := range getBillCostServiceDay(e.DBC) {
		e.BillCostServiceByDay.WithLabelValues(res.DateStart, res.DateEnd, res.ServiceName).Set(res.BlendedCostUSD)
	}
	for _, res := range getAllBillCostServiceMonthly(e.DBC) {
		e.BillCostServiceByMonthly.WithLabelValues(res.DateStart, res.DateEnd, res.ServiceName).Set(res.BlendedCostUSD)
	}
	instances, err := getLightsailInstances()
	if err != nil {
		log.Println("Warn metrics in getLightsailInstances error:", err.Error())
	}
	if err == nil {
		for _, i := range instances {
			e.LightsailInstances.WithLabelValues(
				i.Name, i.Regions, i.RegionName, i.RegionNameAbbreviation, i.BlueprintName, i.CreatedAt.String(), strconv.FormatInt(i.CpuCount, 10),
				fmt.Sprintf("%f", i.RamSizeInGb), i.PrivateIpAddress, i.PublicIpAddress, fmt.Sprintf("%f", i.GbPerMonthAllocatedTransfer),
			).Set(i.State)
			e.LightsailInstancesGbPerMonthAllocatedTransfer.WithLabelValues(i.Name, i.Regions, i.RegionName, i.RegionNameAbbreviation, i.PublicIpAddress).Set(i.GbPerMonthAllocatedTransfer)
			e.LightsailInstancesGbMonthNetworkIn.WithLabelValues(i.Name, i.Regions, i.RegionName, i.RegionNameAbbreviation, i.PublicIpAddress).Set(i.GbMonthNetworkIn)
			e.LightsailInstancesGbMonthNetworkOut.WithLabelValues(i.Name, i.Regions, i.RegionName, i.RegionNameAbbreviation, i.PublicIpAddress).Set(i.GbMonthNetworkOut)
			e.LightsailInstancesGbUseMonthNetwork.WithLabelValues(i.Name, i.Regions, i.RegionName, i.RegionNameAbbreviation, i.PublicIpAddress).Set(i.GbUseMonthNetwork)
			e.LightsailInstancesGbRemainMonthNetwork.WithLabelValues(i.Name, i.Regions, i.RegionName, i.RegionNameAbbreviation, i.PublicIpAddress).Set(i.GbRemainMonthNetwork)
		}
	}

	e.BillCostServiceByDay.Collect(ch)
	e.BillCostServiceByMonthly.Collect(ch)
	e.LightsailInstances.Collect(ch)
	e.LightsailInstancesGbPerMonthAllocatedTransfer.Collect(ch)
	e.LightsailInstancesGbMonthNetworkIn.Collect(ch)
	e.LightsailInstancesGbMonthNetworkOut.Collect(ch)
	e.LightsailInstancesGbUseMonthNetwork.Collect(ch)
	e.LightsailInstancesGbRemainMonthNetwork.Collect(ch)
}

func getBillCostServiceDay(db *gorm.DB) []infrastructure.BillCostServiceDay {
	res, err := infrastructure.SelectBillCostServiceDay(db, date.GetBeforeDay(-7), date.GetNowDay())
	if err != nil {
		log.Println("Failed Metrics getBillCostServiceDay error: ", err)
		return nil
	}
	return res
}

// getAllBillCostServiceMonthly 获取所有月账单记录
func getAllBillCostServiceMonthly(db *gorm.DB) []infrastructure.BillCostServiceMonthly {
	res, err := infrastructure.SelectAllBillCostServiceMonthly(db)
	if err != nil {
		log.Println("Failed Metrics getAllBillCostServiceMonthly error: ", err)
		return nil
	}
	return res
}

// getLightsailInstances 获取光帆节点信息
func getLightsailInstances() ([]get_lightsail.Instances, error) {
	var instances []get_lightsail.Instances
	lightsailProfileName := viper.GetString("aws.credentials.lightsail")
	conn := connect.NewAwsConnect()
	if err1 := conn.Conn(lightsailProfileName, ""); err1 != nil {
		log.Println("getLightsailInstances run aws conn error: ", err1.Error())
		return instances, err1
	}
	lightsailCli := get_lightsail.NewGetLightsail(conn.Sess)
	regionsNames, err1 := lightsailCli.GetRegionsName()
	if err1 != nil {
		return instances, err1
	}
	for _, rn := range regionsNames {
		i, err := lightsailCli.GetInstances(rn)
		if err != nil {
			log.Printf("Warn getLightsailInstances in GetInstances error (regionsNames: %s): %s\n", rn, err.Error())
			continue
		}
		if len(i) == 0 {
			continue
		}
		instances = append(instances, i...)
		log.Printf("Info getLightsailInstances in GetInstances (regionsNames: %s) Successfully added %v instances", rn, len(i))
	}
	return instances, nil
}
