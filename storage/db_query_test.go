package storage

import (
	"github.com/jinzhu/gorm"
	"reflect"
	"testing"
)

type UserInfo struct {
}

func TestNewDBQuery(t *testing.T) {
	type args struct {
		db  *gorm.DB
		mdl interface{}
	}
	tests := []struct {
		name string
		args args
		want *DBQuery
	}{
		{
			name: "测试",
			args: args{
				db:  new(gorm.DB),
				mdl: new(UserInfo),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDBQuery(tt.args.db, tt.args.mdl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDBQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
