package get_billcost

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/viper"
	"gorm.io/gorm"

	"infra-aws-operator/internal/databases/infrastructure"
	"infra-aws-operator/pkg/telegram"
	"infra-aws-operator/pkg/utils/date"
)

type CostCallJob struct {
	DBCRUD *gorm.DB
}

func NewBCostCallJob(dbCRUD *gorm.DB) *CostCallJob {
	return &CostCallJob{
		DBCRUD: dbCRUD,
	}
}

func (j *CostCallJob) Run() {
	log.Println("BillcostJob run start")
	res := getYesterdayBillCostServiceByDay(j.DBCRUD)
	totalCost := geTotalBillCost(res)
	botToken := viper.GetString("telegram.cost_call.bot_token")
	chatId := viper.GetString("telegram.cost_call.chat_id")
	callUser := viper.GetString("telegram.cost_call.call_user")
	var content bytes.Buffer
	content.WriteString("*成本账单* " + callUser + "\n" +
		"  昨日账单成本: $" + fmt.Sprintf("%v", totalCost) + "\n" +
		"------\n")
	for i, v := range res {
		content.WriteString(strconv.Itoa(i+1) + ". " + v.ServiceName + fmt.Sprintf(" $%v", v.BlendedCostUSD))
		if i+1 != len(res) {
			content.WriteString("\n")
		}
	}
	if err := telegram.SendMessage(botToken, chatId, content.String()); err != nil {
		log.Println("BillcostJob run SendMessage error: ", err.Error())
		log.Println("BillcostJob run done - fail")
		return
	}
	log.Println("BillcostJob run done - success")
}

// getYesterdayBillCostServiceByDay 获取昨天的成本
func getYesterdayBillCostServiceByDay(db *gorm.DB) []infrastructure.BillCostServiceDay {
	res, err := infrastructure.SelectBillCostServiceDay(db, date.GetBeforeDay(-1), date.GetNowDay())
	if err != nil {
		log.Println("Failed CostCallJob getYesterdayBillCostServiceByDay error: ", err)
		return nil
	}
	return res
}

// get total cost of yesterday bill
func geTotalBillCost(res []infrastructure.BillCostServiceDay) float64 {
	var t float64
	for _, c := range res {
		t = t + c.BlendedCostUSD
	}
	return t
}
