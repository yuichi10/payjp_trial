package dbase

import (
	_ "database/sql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)
/*
    insert into users (customer_card_id, name) values ("123abc1", "name_1");
    insert into users (customer_card_id, name) values ("123abc2", "name_2");
    insert into users (customer_card_id, name) values ("123abc3", "name_3");
*/

type User struct {
    gorm.Model
    CustomerCardID string `gorm:"column:customer_card_id;size:50"`
    Name string `gorm:"column:name;size:50"`

    Items []Item `gorm:"ForeignKey:UserID;"`
    Orders []Order `gorm:"ForeignKey:UserID;"`
}