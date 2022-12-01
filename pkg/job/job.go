package job

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"

	"infra-aws-operator/pkg/job/get_billcost"
)

func Job() {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(fmt.Errorf("Failed LoadLocation err: %s ", err.Error()))
	}
	cron1 := cron.New(cron.WithLocation(location), cron.WithSeconds())
	cron1.Start()

	// Aws Cost Explorer 一个请求，一美分，你准备好了嘛？
	cron1.AddJob("01 08 * * *", get_billcost.NewBillcostJob())
}
