package account

import (
	_ "database/sql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

/*
   insert into items (user_id, name) values (1, "p_name_1");
   insert into items (user_id, name) values (2, "p_name_2");
   insert into items (user_id, name) values (3, "p_name_3");
   insert into items (user_id, name) values (4, "p_name_4");

   alter table items add constraint user_item_fk foreign key (user_id) references users(id) on delete cascade on update cascade;
*/
const (
	ItemID = "item_id"
)
type Item struct {
	gorm.Model
	UserID uint
	Name   string `gorm:"column:name;size:50"`
	BasePrice int	`gorm:"column:base_price;not null"`
	DailyCharge int	`gorm:"column:daily_charge;not null"`
	DepositFee int	`gorm:"column:deposit_fee;not null"`
	Orders []Order `gorm:"ForeignKey:ItemID;"`
}

func getItemInfo(itemID uint64, db *gorm.DB) *Item {
	item := new(Item)
	db.Where("id=?", itemID).First(&item)
	if item.ID == 0 {
		return nil
	}
	return item
}
