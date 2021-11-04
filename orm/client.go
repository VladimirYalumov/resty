package orm

import (
	"fmt"
	"time"
)

type Client struct {
	Id         int         `gorm:"primary_key" json:"id"`
	Name       string      `gorm:"type:varchar(100);unique;not null" json:"name"`
	Key        string      `gorm:"type:varchar(255);unique;not null" json:"key"`
	AuthTokens []AuthToken `gorm:"references:UserId"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (model *Client) GetByStringValue(field string, value string) (bool, error) {
	err := db.First(&model, fmt.Sprintf("%s  = ?", field), value)
	if err.Error != nil {
		if IfEmpty(err.Error.Error()) {
			return false, nil
		}
		return false, err.Error
	}
	return true, nil
}

func (model *Client) GetByIntValue(field string, value int) (bool, error) {
	err := db.First(&model, fmt.Sprintf("%s  = ?", field), value)
	if err.Error != nil {
		if IfEmpty(err.Error.Error()) {
			return false, nil
		}
		return false, err.Error
	}
	return true, nil
}

func (model *Client) Create() error {
	err := db.Create(&model)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (model *Client) Delete() error {
	err := db.Delete(&model)
	if err.Error != nil {
		return err.Error
	}
	return nil
}
