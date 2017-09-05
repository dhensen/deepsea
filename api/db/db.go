package db

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open("mysql", "root:deepsea@tcp(localhost:3307)/test_deepsea?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Panicln(err)
	}
}
