// Package metrics author: weiqiang; date: 2022-12
package metrics

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"

	"infra-aws-operator/internal/databases/infrastructure"
	"infra-aws-operator/pkg/utils/date"
)

const (
	namespace = "aws"
)

type Exporter struct {
	DBC                      *gorm.DB
	BillCostServiceByDay     *prometheus.GaugeVec
	BillCostServiceByMonthly *prometheus.GaugeVec
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
	}
}
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.BillCostServiceByDay.Describe(ch)
	e.BillCostServiceByMonthly.Describe(ch)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.BillCostServiceByDay.Reset()
	e.BillCostServiceByMonthly.Reset()

	for _, res := range getBillCostServiceDay(e.DBC) {
		e.BillCostServiceByDay.WithLabelValues(res.DateStart, res.DateEnd, res.ServiceName).Set(res.BlendedCostUSD)
	}
	for _, res := range getAllBillCostServiceMonthly(e.DBC) {
		e.BillCostServiceByMonthly.WithLabelValues(res.DateStart, res.DateEnd, res.ServiceName).Set(res.BlendedCostUSD)
	}

	e.BillCostServiceByDay.Collect(ch)
	e.BillCostServiceByMonthly.Collect(ch)
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
