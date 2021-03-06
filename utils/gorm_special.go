package utils

import (
	"database/sql/driver"
	"github.com/json-iterator/go"
)

var (
	GORM_JSON_FIELD bool = false
)

type GormStrings []string

func (static GormStrings) Value() (driver.Value, error) {
	if GORM_JSON_FIELD {
		return jsoniter.MarshalToString(static)
	} else {
		return jsoniter.Marshal(static)
	}
}

func (this *GormStrings) Scan(v interface{}) error {
	var strs []string
	if err := jsoniter.Unmarshal(v.([]byte), &strs); err != nil {
		return err
	}

	*this = strs
	return nil
}

type GormInts []int

func (static GormInts) Value() (driver.Value, error) {
	if GORM_JSON_FIELD {
		return jsoniter.MarshalToString(static)
	} else {
		return jsoniter.Marshal(static)
	}
}

func (this *GormInts) Scan(v interface{}) error {
	var strs []int
	if err := jsoniter.Unmarshal(v.([]byte), &strs); err != nil {
		return err
	}

	*this = strs
	return nil
}

type GormInt64s []int64

func (static GormInt64s) Value() (driver.Value, error) {
	if GORM_JSON_FIELD {
		return jsoniter.MarshalToString(static)
	} else {
		return jsoniter.Marshal(static)
	}
}

func (this *GormInt64s) Scan(v interface{}) error {
	var strs []int64
	if err := jsoniter.Unmarshal(v.([]byte), &strs); err != nil {
		return err
	}

	*this = strs
	return nil
}

type GormFloat64s []float64

func (static GormFloat64s) Value() (driver.Value, error) {
	if GORM_JSON_FIELD {
		return jsoniter.MarshalToString(static)
	} else {
		return jsoniter.Marshal(static)
	}
}

func (this *GormFloat64s) Scan(v interface{}) error {
	var strs []float64
	if err := jsoniter.Unmarshal(v.([]byte), &strs); err != nil {
		return err
	}

	*this = strs
	return nil
}

type GormMap map[string]interface{}

func (static GormMap) Value() (driver.Value, error) {
	if GORM_JSON_FIELD {
		return jsoniter.MarshalToString(static)
	} else {
		return jsoniter.Marshal(static)
	}
}

func (this *GormMap) Scan(v interface{}) error {
	var strs map[string]interface{}
	if err := jsoniter.Unmarshal(v.([]byte), &strs); err != nil {
		return err
	}

	*this = strs
	return nil
}

type GormMapString map[string]string

func (static GormMapString) Value() (driver.Value, error) {
	if GORM_JSON_FIELD {
		return jsoniter.MarshalToString(static)
	} else {
		return jsoniter.Marshal(static)
	}
}

func (this *GormMapString) Scan(v interface{}) error {
	var strs map[string]string
	if err := jsoniter.Unmarshal(v.([]byte), &strs); err != nil {
		return err
	}

	*this = strs
	return nil
}
