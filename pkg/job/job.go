package job

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"infra-aws-operator/pkg/job/get_billcost"
)

func Job(dbCRUD *gorm.DB) {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(fmt.Errorf("Failed LoadLocation err: %s ", err.Error()))
	}
	cron1 := cron.New(cron.WithLocation(location), cron.WithSeconds())
	cron1.Start()

	billcostTime := viper.GetString("job.get_billcost")
	if len(billcostTime) == 0 {
		billcostTime = "01 01 01 * * *"
	}
	// Aws Cost Explorer 一个请求，一美分，你准备好了嘛？
	id, err := cron1.AddJob(billcostTime, get_billcost.NewBillcostJob(dbCRUD))
	if err != nil {
		panic(fmt.Errorf("Failed Job AddJob get_billcost %v error: %s ", id, err.Error()))
	}

	costCallTime := viper.GetString("job.cost_call")
	if len(costCallTime) == 0 {
		costCallTime = "01 01 09 * * *"
	}
	if costCallTime != "off" {
		id, err = cron1.AddJob(costCallTime, get_billcost.NewBCostCallJob(dbCRUD))
		if err != nil {
			panic(fmt.Errorf("Failed Job AddJob get_billcost %v error: %s ", id, err.Error()))
		}
	}
}
