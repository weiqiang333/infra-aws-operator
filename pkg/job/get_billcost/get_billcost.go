package get_billcost

import (
	"fmt"
	"log"

	"infra-aws-operator/pkg/aws/billcost"
	"infra-aws-operator/pkg/aws/connect"
)

type BillcostJob struct {
}

func NewBillcostJob() *BillcostJob {
	return &BillcostJob{}
}

func (j *BillcostJob) Run() {
	log.Println("BillcostJob run start")
	conn := connect.NewAwsConnect()
	if err := conn.Conn(); err != nil {
		log.Println("BillcostJob run aws conn error: ", err.Error())
		return
	}
	billCli := billcost.NewBillcost(conn.Sess)
	if err := billCli.GetDailyData(); err != nil {
		log.Println("BillcostJob run GetDailyData error: ", err.Error())
		return
	}
	for _, res := range billCli.Res {
		fmt.Println(res.TimePeriod, res.Total, res.Groups)
		if !*res.Estimated {
			log.Println("Warn BillcostJob GetDailyData TimePeriod no is true, ", *res.TimePeriod)
			continue
		}
		for _, v := range res.Groups {
			fmt.Println(*res.TimePeriod.Start, *res.TimePeriod.End, v.Keys, *v.Metrics["BlendedCost"].Amount)
		}
	}

	if err := billCli.GetMonthlyData(); err != nil {
		log.Println("BillcostJob run GetMonthlyData error: ", err.Error())
		return
	}

	log.Println("BillcostJob run done")
}
