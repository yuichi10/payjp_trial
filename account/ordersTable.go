package account

import (
	_ "database/sql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
	"fmt"
	"D"
)

/*
	alter table orders add constraint user_order_fk foreign key (user_id) references users(id) on delete cascade on update cascade;
	alter table orders add constraint item_order_fk foreign key (user_id) references users(id) on delete cascade on update cascade;
*/

const (
	OrderID = "order_id"
	RentalFrom = "rental_from"
	RentalTo = "rental_to"
	Status = "status"
)

type Order struct {
	gorm.Model
	OrderChargeID     string     `gorm:"column:order_charge_id;size:50" json:"order_charge_id"`
	TransportAllocate int        `gorm:"column:transport_allocate;not null;default:0" json:"transport_allocate"`
	RentalFrom        *time.Time `gorm:"column:rental_from;" json:"rental_from"`
	RentalTo          *time.Time `gorm:"column:rental_to;" json:"rental_to"`
	ItemID            uint		 `json:"item_id"`
	UserID            uint		 `json:"user_id"`
	BasePrice         int        `gorm:"column:base_price;not null" json:"base_price"`
	DailyCharge       int        `gorm:"column:daily_charge;not null" json:"daily_charge"`
	InsurancePrice    int        `gorm:"column:insurance_price;not null" json:"insurance_price"`
	ManagementCharge  int        `gorm:"column:management_charge;not null" json:"management_charge"`
	DepositFee        int        `gorm:"column:deposit_fee;not null;default:0" json:"deposit_fee"`
	Amount            int        `gorm:"column:amount;not null;" json:"amount"`
	CancelDate        *time.Time `gorm:"column:cancel_date" json:"cancel_date"`
	CancelStatus      int        `gorm:"column:cancel_status;not null" json:"cancel_status"`
	Status            int        `gorm:"column:status;not null;default:0" json:"status"`
}

//ダブルブッキングの数を返す
func getDoubleBookingNum(rFrom, rTo *time.Time, itemID uint, state int, db *gorm.DB) int {
	count := 0
	count += countOtherOverlapBook(rFrom, rTo, itemID, state, db)
	count += countOtherCoverBook(rFrom, rTo, itemID, state, db)
	return count
}
//期間の一部もしくは全部がかぶっている予約されている期間があるかどうかの検索
func countOtherOverlapBook(rFrom, rTo *time.Time, itemID uint, state int, db *gorm.DB) int {
	count := 0
	rMerginFrom := rFrom.AddDate(0,0,-D.BOOK_MARGIN_DAYS)
	rMerginTo   := rTo.AddDate(0,0,D.BOOK_MARGIN_DAYS)
	baseWhere := fmt.Sprintf("(%v=? AND %v=?) AND (? BETWEEN %v AND %v OR ? BETWEEN %v AND %v)", ItemID, Status, RentalFrom, RentalTo, RentalFrom, RentalTo)
	db.Model(&Order{}).Where(baseWhere, itemID, state, &rMerginFrom, &rMerginTo).Count(&count)
	return count
}
//すべての期間を含む場合
func countOtherCoverBook(rFrom, rTo *time.Time, itemID uint, state int, db *gorm.DB) int {
	count := 0
	rMerginFrom := rFrom.AddDate(0,0,-D.BOOK_MARGIN_DAYS)
	rMerginTo   := rTo.AddDate(0,0,D.BOOK_MARGIN_DAYS)
	baseWhere := fmt.Sprintf("(%v=? AND %v=?) AND (? < %v AND ? > %v)", ItemID, Status, RentalFrom, RentalTo)
	db.Model(&Order{}).Where(baseWhere, itemID, state, &rMerginFrom, &rMerginTo).Count(&count)
	return count
}

func (order *Order) calcPureRentalPrice() int {
	price := 0
	period := calcSubDate(order.RentalFrom, order.RentalTo)
	if period < 0 {
		period = 0
	}
	price = order.BasePrice + period * order.DailyCharge
	return price
}