package orm

import (
	"fmt"
	"time"
)

type User struct {
	Id         int         `gorm:"primary_key" json:"id"`
	Name       string      `gorm:"type:varchar(15);not null" json:"name"`
	Email      string      `gorm:"type:varchar(100);unique;not null" json:"email"`
	Phone      string      `gorm:"type:varchar(15)" json:"phone"`
	Password   string      `gorm:"type:varchar(255)" json:"password"`
	Active     bool        `gorm:"type:boolean;default:false" json:"active"`
	AuthTokens []AuthToken `gorm:"references:user_id"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func VerifyUser(user *User, verify bool) error {
	err := db.Model(&user).Updates(User{Active: verify})
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (model *User) GetByStringValue(field string, value string) (bool, error) {
	err := db.First(&model, fmt.Sprintf("%s  = ?", field), value)
	if err.Error != nil {
		if IfEmpty(err.Error.Error()) {
			return false, nil
		}
		return false, err.Error
	}
	return true, nil
}

func (model *User) GetByIntValue(field string, value int) (bool, error) {
	err := db.First(&model, fmt.Sprintf("%s  = ?", field), value)
	if err.Error != nil {
		if IfEmpty(err.Error.Error()) {
			return false, nil
		}
		return false, err.Error
	}
	return true, nil
}

func (model User) Create() error {
	err := db.Create(&model)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (model *User) Delete() error {
	err := db.Delete(&model)
	if err.Error != nil {
		return err.Error
	}
	return nil
}
