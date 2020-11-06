package utils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
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

type JSON json.RawMessage

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}
