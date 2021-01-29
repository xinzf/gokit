package storage

import (
	"fmt"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strings"
)

//type Column interface {
//	Name() string
//	DatabaseTypeName() string
//	Length() (length int64, ok bool)
//	DecimalSize() (precision int64, scale int64, ok bool)
//	Nullable() (nullable bool, ok bool)
//}

type Column struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Length     int    `json:"length"`
	Decimal    int    `json:"decimal"`
	Nullable   bool   `json:"nullable"`
	PrimaryKey bool   `json:"primary_key"`
	Comment    string `json:"comment"`
}

type Columns []*Column

type DbTable struct {
	tableName string
	db        *gorm.DB
	columns   Columns
}

func NewDbTable(tableName string, db *gorm.DB) *DbTable {
	return &DbTable{tableName: tableName, db: db}
}

func (this *DbTable) GetColumns() (Columns, error) {
	if this.columns != nil {
		return this.columns, nil
	}

	this.columns = Columns{}

	sql := fmt.Sprintf("SHOW FULL FIELDS FROM %s", this.tableName)
	mp := make([]map[string]interface{}, 0)
	err := this.db.Raw(sql).Find(&mp).Error
	if err != nil {
		return this.columns, err
	}

	for _, v := range mp {
		typ := strings.Split(cast.ToString(v["Type"]), "(")[0]
		typ = strings.Split(typ, " ")[0]
		label := cast.ToString(v["Comment"])
		if label == "" {
			label = cast.ToString(v["Field"])
		}

		this.columns = append(this.columns, &Column{
			Name:       cast.ToString(v["Field"]),
			Type:       typ,
			Length:     0,
			Decimal:    0,
			Nullable:   false,
			PrimaryKey: false,
			Comment:    label,
		})
	}
	return this.columns, nil
}
