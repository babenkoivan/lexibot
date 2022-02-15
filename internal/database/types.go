package database

import (
	"database/sql/driver"
	"strings"
)

type StringArray []string

func (a *StringArray) Scan(src interface{}) error {
	str, _ := src.(string)
	// todo error handling
	*a = strings.Split(str, ",")
	return nil
}

func (a StringArray) Value() (driver.Value, error) {
	if a == nil || len(a) == 0 {
		return nil, nil
	}

	return strings.Join(a, ","), nil
}

func (StringArray) GormDataType() string {
	return "text"
}
