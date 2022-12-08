package billcost

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"

	"infra-aws-operator/pkg/utils/date"
)

type Billcost struct {
	Sess *session.Session
	Res  []*costexplorer.ResultByTime
}

func NewBillcost(sess *session.Session) *Billcost {
	return &Billcost{
		Sess: sess,
		Res:  nil,
	}
}

// GetDailyData 获取之前30天的数据; 聚合以 ${granularity}, 类型为 Service
func (b *Billcost) GetDailyData(granularity string) error {
	cli := costexplorer.New(b.Sess)

	startTime := date.GetLastMonth1stDay()
	endTime := date.GetNowDay()
	cu, err := cli.GetCostAndUsage(&costexplorer.GetCostAndUsageInput{
		Granularity: aws.String(granularity),
		Metrics:     aws.StringSlice([]string{"BlendedCost"}),
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(startTime),
			End:   aws.String(endTime),
		},
		GroupBy: []*costexplorer.GroupDefinition{&costexplorer.GroupDefinition{
			Type: aws.String("DIMENSION"),
			Key:  aws.String("SERVICE"),
		}},
	})
	if err != nil {
		msg := fmt.Sprintf("Failed aws api GetDailyData err: %s", err.Error())
		log.Println(msg)
		return err
	}
	b.Res = cu.ResultsByTime
	return nil
}
