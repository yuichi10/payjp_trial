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

type Item struct {
	gorm.Model
	UserID int
	Name   string `gorm:"column:name;size:50"`

	Orders []Order `gorm:"ForeignKey:ItemID;"`
}
