package get_lightsail

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/aws/aws-sdk-go/service/lightsail"

	"infra-aws-operator/pkg/utils/date"
)

type GetLightsail struct {
	Sess *session.Session
	Res  []*costexplorer.ResultByTime
}

func NewGetLightsail(sess *session.Session) *GetLightsail {
	return &GetLightsail{
		Sess: sess,
		Res:  nil,
	}
}

type Instances struct {
	Name                        string    `json:"name"`           // 节点名称
	State                       float64   `json:"state"`          //节点状态 0 异常 1 running
	RegionsNames                string    `json:"regions_names"`  // 区域名称
	BlueprintName               string    `json:"blueprint_name"` // 系统类型名称
	CreatedAt                   time.Time `json:"created_at"`
	CpuCount                    int64     `json:"cpu_count"`
	RamSizeInGb                 float64   `json:"ram_size_in_gb"`                  // 内存大小
	GbPerMonthAllocatedTransfer float64   `json:"gb_per_month_allocated_transfer"` // 当月流量套餐 sum GB
	PrivateIpAddress            string    `json:"private_ip_address"`
	PublicIpAddress             string    `json:"public_ip_address"`
	GbMonthNetworkIn            float64   `json:"gb_month_network_in"`     // 当月流量流入 sum
	GbMonthNetworkOut           float64   `json:"gb_month_network_out"`    // 当月流量流出 sum
	GbUseMonthNetwork           float64   `json:"gb_use_month_network"`    // 当月流量使用 sum
	GbRemainMonthNetwork        float64   `json:"gb_remain_month_network"` // 当月流量套餐剩余 sum
}

// GetRegionsName 获取所有区域名称
func (l *GetLightsail) GetRegionsName() ([]string, error) {
	var regionsName []string
	cli := lightsail.New(l.Sess)
	out, err := cli.GetRegions(&lightsail.GetRegionsInput{})
	if err != nil {
		msg := fmt.Sprintf("Failed aws api GetLightsail in GetRegionsName GetRegions err: %s", err.Error())
		log.Println(msg)
		return regionsName, err
	}
	for _, r := range out.Regions {
		regionsName = append(regionsName, *r.Name)
	}
	return regionsName, nil
}

// GetInstances 获取当前区域所有节点统计信息
func (l *GetLightsail) GetInstances(region string) ([]Instances, error) {
	var instances []Instances
	l.Sess.Config.Region = aws.String(region)
	cli := lightsail.New(l.Sess)
	out, err := cli.GetInstances(&lightsail.GetInstancesInput{})
	if err != nil {
		msg := fmt.Sprintf("Failed aws api GetLightsail in GetInstances err: %s", err.Error())
		log.Println(msg)
		return instances, err
	}
	for _, i := range out.Instances {
		var state float64
		if *i.State.Name == "running" {
			state = 1
		}
		sumNetworkInGb, sumNetworkOutGb, sumNetworkGb, err := l.GetInstanceMetricData(cli, *i.Name)
		if err != nil {
			msg := fmt.Sprintf("Failed aws api GetLightsail in GetInstances GetInstanceMetricData err: %s", err.Error())
			log.Println(msg)
		}
		gbRemainMonthNetwork := float64(*i.Networking.MonthlyTransfer.GbPerMonthAllocated) - sumNetworkGb
		instances = append(instances, Instances{
			Name:                        *i.Name,
			State:                       state,
			RegionsNames:                region,
			BlueprintName:               *i.BlueprintName,
			CreatedAt:                   *i.CreatedAt,
			CpuCount:                    *i.Hardware.CpuCount,
			RamSizeInGb:                 *i.Hardware.RamSizeInGb,
			GbPerMonthAllocatedTransfer: float64(*i.Networking.MonthlyTransfer.GbPerMonthAllocated),
			PrivateIpAddress:            *i.PrivateIpAddress,
			PublicIpAddress:             *i.PublicIpAddress,
			GbMonthNetworkIn:            sumNetworkInGb,
			GbMonthNetworkOut:           sumNetworkOutGb,
			GbUseMonthNetwork:           sumNetworkGb,
			GbRemainMonthNetwork:        gbRemainMonthNetwork,
		})
	}
	return instances, nil
}

// GetInstanceMetricData 获取节点 metrics 数据，返回流量数据统计: sumNetworkInGb, sumNetworkOutGb, sumNetworkGb, error
func (l *GetLightsail) GetInstanceMetricData(cli *lightsail.Lightsail, instanceName string) (float64, float64, float64, error) {
	var sumNetworkIn float64 = 0
	var sumNetworkOut float64 = 0
	var sumNetwork float64 = 0
	startTime := date.GetNowMonth1stDay()
	endTime := date.GetNextMonth1stDay()
	outNetworkIn, err := cli.GetInstanceMetricData(&lightsail.GetInstanceMetricDataInput{
		InstanceName: aws.String(instanceName),
		MetricName:   aws.String("NetworkIn"),
		Period:       aws.Int64(86400),
		StartTime:    aws.Time(startTime),
		EndTime:      aws.Time(endTime),
		Statistics:   aws.StringSlice([]string{"Sum"}),
		Unit:         aws.String("Bytes"),
	})
	if err != nil {
		msg := fmt.Sprintf("Failed aws api GetLightsail in GetInstanceMetricData GetInstanceMetricData err (instanceName: %s, NetworkIn): %s", instanceName, err.Error())
		log.Println(msg)
		return sumNetworkIn, sumNetworkOut, sumNetwork, err
	}
	for _, i := range outNetworkIn.MetricData {
		sumNetworkIn = sumNetworkIn + *i.Sum
	}

	outNetworkOut, err := cli.GetInstanceMetricData(&lightsail.GetInstanceMetricDataInput{
		InstanceName: aws.String(instanceName),
		MetricName:   aws.String("NetworkOut"),
		Period:       aws.Int64(86400),
		StartTime:    aws.Time(startTime),
		EndTime:      aws.Time(endTime),
		Statistics:   aws.StringSlice([]string{"Sum"}),
		Unit:         aws.String("Bytes"),
	})
	if err != nil {
		msg := fmt.Sprintf("Failed aws api GetLightsail in GetInstanceMetricData GetInstanceMetricData err (instanceName: %s, NetworkOut): %s", instanceName, err.Error())
		log.Println(msg)
		return sumNetworkIn, sumNetworkOut, sumNetwork, err
	}
	for _, i := range outNetworkOut.MetricData {
		sumNetworkOut = sumNetworkOut + *i.Sum
	}

	sumNetwork = sumNetworkIn + sumNetworkOut
	return sumNetworkIn / 1024 / 1024 / 1024, sumNetworkOut / 1024 / 1024 / 1024, sumNetwork / 1024 / 1024 / 1024, nil
}
