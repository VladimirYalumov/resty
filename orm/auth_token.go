package orm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"math/rand"
	"time"
)

type AuthToken struct {
	ClientId  int    `gorm:"primary_key;auto_increment:false" json:"client_id"`
	UserId    int    `gorm:"primary_key;auto_increment:false" json:"user_id"`
	Token     string `gorm:"type:varchar(255);not null" json:"key"`
	PushToken string `gorm:"type:varchar(255)" json:"push_token"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetAuthToken(userId int, clientId int) (AuthToken, error) {
	var authToken AuthToken

	err := db.First(&authToken, "user_id = ? AND client_id = ?", userId, clientId)
	authTokenExist := true

	if err.Error != nil {
		if IfEmpty(err.Error.Error()) {
			authTokenExist = false
		} else {
			return authToken, err.Error
		}
	}

	if !authTokenExist {
		authToken = AuthToken{
			ClientId: clientId,
			UserId:   userId,
			Token:    tokenGenerator(),
		}
		err = db.Create(&authToken)
		if err.Error != nil {
			return AuthToken{}, err.Error
		}

		return authToken, nil
	}

	return authToken, nil
}

func tokenGenerator() string {
	b := make([]byte, 20)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (model *AuthToken) GetByStringValue(field string, value string) (bool, error) {
	err := db.First(&model, fmt.Sprintf("%s  = ?", field), value)
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err.Error
	}
	return true, nil
}

func (model *AuthToken) GetByIntValue(field string, value int) (bool, error) {
	err := db.First(&model, fmt.Sprintf("%s  = ?", field), value)
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err.Error
	}
	return true, nil
}

func (model *AuthToken) Create() error {
	err := db.Create(&model)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (model *AuthToken) Delete() error {
	err := db.Delete(&model)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func CheckAuth(userId int, token string, clientKey string) (bool, error) {
	var authToken AuthToken
	var client Client
	_, clientErr := client.GetByStringValue("key", clientKey)

	if clientErr != nil {
		return false, clientErr
	}

	err := db.First(&authToken, "user_id = ? AND token = ? AND client_id = ?", userId, token, client.Id)

	if err.Error != nil {
		if IfEmpty(err.Error.Error()) {
			return false, nil
		} else {
			return false, err.Error
		}
	}

	return true, nil
}

func DeleteToken(userId int, clientKey string) error {
	var authToken AuthToken
	var client Client
	_, clientErr := client.GetByStringValue("key", clientKey)
	if clientErr != nil {
		return clientErr
	}

	err := db.Delete(&authToken, "user_id = ? AND client_id = ?", userId, client.Id)
	if err.Error != nil {
		return err.Error
	}

	return nil
}
