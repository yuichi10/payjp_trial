package dbase

import (
	_ "database/sql"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
    "time"
)

type Order struct {
	gorm.Model
	OrderChargeID     string      `gorm:"column:order_charge_id;size:50"`
	TransportAllocate int         `gorm:"column:transport_allocate;not null;default:0"`
	RentalFrom        *time.Time  `gorm:"column:rental_from;"`
	RentalTo		  *time.Time  `gorm:"column:rental_to;"`
	ItemID            int         
	UserID            int         
	DayPrice          int         `gorm:"column:day_price;not null"`
	AfterDayPrice     int         `gorm:"column:after_day_price;not null"`
	InsurancePrice    int         `gorm:"column:insurance_price;not null"`
	ManagementCharge  int         `gorm:"column:management_charge;not null"`
	Amount            int         `gorm:"column:amount;not null;"`
	CancelDate        *time.Time `gorm:"column:cancel_date"`
	CancelStatus      int         `gorm:"column:cancel_status;not null"`
	Status            int         `gorm:"column:status;not null;default:0"`
}
