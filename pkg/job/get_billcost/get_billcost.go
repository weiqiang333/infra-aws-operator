package get_billcost

import (
	"fmt"
	"log"
	"strconv"

	"gorm.io/gorm"

	"infra-aws-operator/internal/databases/infrastructure"
	"infra-aws-operator/pkg/aws/billcost"
	"infra-aws-operator/pkg/aws/connect"
)

type BillcostJob struct {
	DBCRUD *gorm.DB
}

func NewBillcostJob(dbCRUD *gorm.DB) *BillcostJob {
	return &BillcostJob{
		DBCRUD: dbCRUD,
	}
}

func (j *BillcostJob) Run() {
	log.Println("BillcostJob run start")
	conn := connect.NewAwsConnect()
	if err := conn.Conn(); err != nil {
		log.Println("BillcostJob run aws conn error: ", err.Error())
		log.Println("BillcostJob run done - fail")
		return
	}
	billCli := billcost.NewBillcost(conn.Sess)
	if err := j.writeData(billCli, "DAILY"); err != nil {
		log.Println("BillcostJob run done - fail")
		return
	}
	if err := j.writeData(billCli, "MONTHLY"); err != nil {
		log.Println("BillcostJob run done - fail")
		return
	}
	log.Println("BillcostJob run done - success")
}

func (j *BillcostJob) writeData(billCli *billcost.Billcost, granularity string) error {
	if err := billCli.GetDailyData(granularity); err != nil {
		log.Println("Failed BillcostJob run GetDailyData error: ", err.Error())
		return fmt.Errorf("Failed BillcostJob run GetDailyData error: %s ", err.Error())
	}
	for _, res := range billCli.Res {
		dateStart := *res.TimePeriod.Start
		dateEnd := *res.TimePeriod.End
		if !*res.Estimated {
			log.Println("Warn BillcostJob GetDailyData TimePeriod no is true, ", dateStart, dateEnd)
			//continue
		}
		for _, v := range res.Groups {
			serviceName := *v.Keys[0]
			blendedCostUSD := *v.Metrics["BlendedCost"].Amount
			blendedCostUSD64, _ := strconv.ParseFloat(blendedCostUSD, 64)
			if blendedCostUSD64 < 0.001 {
				log.Println("Info job run blendedCostUSD64 is 0 not write:", dateStart, dateEnd, serviceName, blendedCostUSD64)
				continue
			}
			if granularity == "DAILY" {
				if err := infrastructure.WriteBillCostServiceDay(j.DBCRUD, dateStart, dateEnd, serviceName, blendedCostUSD64); err != nil {
					log.Println("Failed job run in WriteBillCostServiceDay error:", err.Error())
					continue
				}
			}
			if granularity == "MONTHLY" {
				if err := infrastructure.WriteBillCostServiceMonthly(j.DBCRUD, dateStart, dateEnd, serviceName, blendedCostUSD64); err != nil {
					log.Println("Failed job run in WriteBillCostServiceMonthly error:", err.Error())
					continue
				}
			}

		}
	}
	return nil
}
