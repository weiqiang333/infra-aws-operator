package infrastructure

import (
	"log"

	"gorm.io/gorm"
)

func CreateTables(db *gorm.DB) error {
	if err := db.AutoMigrate(&BillCostServiceDay{}, &BillCostServiceMonthly{}); err != nil {
		log.Println("Failed db init CreateTables error:", err.Error())
		return err
	}
	return nil
}
