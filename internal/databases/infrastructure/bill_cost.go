package infrastructure

import (
	"log"

	"gorm.io/gorm"

	"infra-aws-operator/pkg/utils/date"
)

// BillCostServiceDay 天账单
//go:generate gormgen -structs BillCostServiceDay -input .
type BillCostServiceDay struct {
	Id             int32   `json:"id" gorm:"column:id;primaryKey"`                                                                      // 主键
	DateStart      string  `json:"date_start" gorm:"column:date_start;comment:'cost start date'"`                                       // 聚合开始日期
	DateEnd        string  `json:"date_end" gorm:"column:date_end;comment:'cost end date'"`                                             // 聚合结束日期
	ServiceName    string  `json:"service_name" gorm:"column:service_name;comment:'service type name'"`                                 // 服务名称
	ServiceSubName string  `json:"service_sub_name" gorm:"column:service_sub_name;default:'aws';comment:'Service custom subtype name'"` // 服务备用自定义子名称
	BlendedCostUSD float64 `json:"blended_cost_usd" gorm:"column:blended_cost_usd;comment:'mixed cost USD'"`                            // 混合成本 USD
	UpdatedAt      int     `json:"updated_at" gorm:"column:updated_at"`                                                                 // 在创建时该字段值为零值或者在更新时，使用当前时间戳秒数填充
}

func WriteBillCostServiceDay(db *gorm.DB, dateStart, dateEnd, serviceName string, blendedCostUSD float64) error {
	res := db.Where(&BillCostServiceDay{
		DateStart:   dateStart,
		DateEnd:     dateEnd,
		ServiceName: serviceName}).Find(&[]BillCostServiceDay{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected >= 1 {
		log.Println("Info WriteBillCostServiceDay is exist not write:", dateStart, dateEnd, serviceName, blendedCostUSD)
		return nil
	}
	res = db.Select("DateStart", "DateEnd", "ServiceName", "BlendedCostUSD").Create(&BillCostServiceDay{
		DateStart:      dateStart,
		DateEnd:        dateEnd,
		ServiceName:    serviceName,
		BlendedCostUSD: blendedCostUSD,
	})
	return res.Error
}

// SelectBillCostServiceDay 查询天账单
func SelectBillCostServiceDay(db *gorm.DB, dateStart, dateEnd string) ([]BillCostServiceDay, error) {
	var d []BillCostServiceDay
	dates := date.GetDateRange(dateStart, dateEnd)
	res := db.Where("date_start IN ? AND service_sub_name = ?", dates, "aws").Find(&d)
	return d, res.Error
}

// BillCostServiceMonthly 月账单
//go:generate gormgen -structs BillCostServiceMonthly -input .
type BillCostServiceMonthly struct {
	Id             int32   `json:"Id" gorm:"column:id;primaryKey"`                                 // 主键
	DateStart      string  `json:"date_start" gorm:"column:date_start"`                            // 聚合开始日期
	DateEnd        string  `json:"date_end" gorm:"column:date_end"`                                // 聚合结束日期
	ServiceName    string  `json:"service_name" gorm:"column:service_name"`                        // 服务名称
	ServiceSubName string  `json:"service_sub_name" gorm:"column:service_sub_name;default:'aws';"` // 服务备用自定义子名称
	BlendedCostUSD float64 `json:"blended_cost_usd" gorm:"column:blended_cost_usd"`                // 混合成本 USD
	UpdatedAt      int     `json:"updated_at" gorm:"column:updated_at"`                            // 在创建时该字段值为零值或者在更新时，使用当前时间戳秒数填充
}

func WriteBillCostServiceMonthly(db *gorm.DB, dateStart, dateEnd, serviceName string, blendedCostUSD float64) error {
	res := db.Where(&BillCostServiceMonthly{
		DateStart:   dateStart,
		ServiceName: serviceName}).Find(&[]BillCostServiceMonthly{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected >= 1 {
		res.Updates(&BillCostServiceMonthly{
			DateStart:      dateStart,
			DateEnd:        dateEnd,
			ServiceName:    serviceName,
			BlendedCostUSD: blendedCostUSD,
		})
		return res.Error
	}
	res = db.Select("DateStart", "DateEnd", "ServiceName", "BlendedCostUSD").Create(&BillCostServiceMonthly{
		DateStart:      dateStart,
		DateEnd:        dateEnd,
		ServiceName:    serviceName,
		BlendedCostUSD: blendedCostUSD,
	})
	return res.Error
}

// SelectAllBillCostServiceMonthly 查询月账单
func SelectAllBillCostServiceMonthly(db *gorm.DB) ([]BillCostServiceMonthly, error) {
	var m []BillCostServiceMonthly
	res := db.Where("service_sub_name = ?", "aws").Find(&m)
	return m, res.Error
}
