package database

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type StringArray []string

func (a *StringArray) Scan(src interface{}) error {
	str, ok := src.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", src)
	}

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
