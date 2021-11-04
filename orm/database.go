package orm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"regexp"
)

var db *gorm.DB

const PrefixRecordNotFound = "record not found"

type Model interface {
	GetByStringValue(field string, value string) (bool, error)
	GetByIntValue(field string, value int) (bool, error)
	Create() error
	Delete() error
}

func InitDB(host string, user string, password string, dbName string, port string) {
	var err error
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	db, err = gorm.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&User{}, &Client{}, &AuthToken{})
	db.AutoMigrate(&User{}, &Client{}, &AuthToken{})
	db.Model(&AuthToken{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&AuthToken{}).AddForeignKey("client_id", "clients(id)", "RESTRICT", "RESTRICT")
}

func CloseDB() {
	err := db.Close()
	if err != nil {
		panic(err)
	}
}

func IfEmpty(err string) bool {
	matched, _ := regexp.MatchString(PrefixRecordNotFound, err)
	if matched {
		return true
	}
	return false
}
