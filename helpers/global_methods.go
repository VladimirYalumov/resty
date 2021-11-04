package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"reflect"
)

// IsEmpty Check that object is empty
func IsEmpty(object interface{}) bool {
	//First check normal definitions of empty
	if object == nil {
		return true
	} else if object == "" {
		return true
	} else if object == false {
		return true
	}

	//Then see if it's a struct
	if reflect.ValueOf(object).Kind() == reflect.Struct {
		// and create an empty copy of the struct object to compare against
		empty := reflect.New(reflect.TypeOf(object)).Elem().Interface()
		if reflect.DeepEqual(object, empty) {
			return true
		}
	}
	return false
}

func EncryptPassword(password string) string {
	encryptedPassword := md5.Sum([]byte(password))
	return hex.EncodeToString(encryptedPassword[:])
}
