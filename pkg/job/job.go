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

	cron1.AddJob("@every 20s", get_billcost.NewBillcostJob())
}
