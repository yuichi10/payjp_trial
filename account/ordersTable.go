package account

import (
	_ "database/sql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

/*
	alter table orders add constraint user_order_fk foreign key (user_id) references users(id) on delete cascade on update cascade;
	alter table orders add constraint item_order_fk foreign key (user_id) references users(id) on delete cascade on update cascade;
*/

const (
	RentalFrom = "rental_from"
	RentalTo = "rental_to"
)

type Order struct {
	gorm.Model
	OrderChargeID     string     `gorm:"column:order_charge_id;size:50"`
	TransportAllocate int        `gorm:"column:transport_allocate;not null;default:0"`
	RentalFrom        *time.Time `gorm:"column:rental_from;"`
	RentalTo          *time.Time `gorm:"column:rental_to;"`
	ItemID            uint
	UserID            uint
	DayPrice          int        `gorm:"column:day_price;not null"`
	AfterDayPrice     int        `gorm:"column:after_day_price;not null"`
	InsurancePrice    int        `gorm:"column:insurance_price;not null"`
	ManagementCharge  int        `gorm:"column:management_charge;not null"`
	DepositPrice      int        `gorm:"column:deposit_price;not null;default:0"`
	Amount            int        `gorm:"column:amount;not null;"`
	CancelDate        *time.Time `gorm:"column:cancel_date"`
	CancelStatus      int        `gorm:"column:cancel_status;not null"`
	Status            int        `gorm:"column:status;not null;default:0"`
}

//期間がまるごと予約されている期間があるかどうかの検索
func countOtherOverlapBook(db *gorm.DB) {
	count := 0
	db.Model(&Order{}).Where("(item_id=? AND status=?) AND (? BETWEEN rental_from AND rental_to OR ? BETWEEN rental_from AND rental_to)").Count(&count)
}