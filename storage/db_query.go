package storage

//
//import (
//	"github.com/jinzhu/gorm"
//	"github.com/xinzf/gokit/utils"
//	"math"
//	"reflect"
//	"strings"
//)
//
//type Preload struct {
//	Name      string
//	Condition string
//	OrderBy   string
//	Limit     int
//}
//
//type DBQuery struct {
//	*gorm.DB
//	preloads        []Preload
//	scopes          []func(db *gorm.DB) *gorm.DB
//	hasDeleted      bool
//	defaultPageSize int
//	defaultPageNo   int
//
//	_mdl     interface{}
//	_tblName string
//}
//
//func NewDBQuery(db *gorm.DB, mdl interface{}) *DBQuery {
//
//	snakeString := func(s string) string {
//		data := make([]byte, 0, len(s)*2)
//		j := false
//		num := len(s)
//		for i := 0; i < num; i++ {
//			d := s[i]
//			if i > 0 && d >= 'A' && d <= 'Z' && j {
//				data = append(data, '_')
//			}
//			if d != '_' {
//				j = true
//			}
//			data = append(data, d)
//		}
//		return strings.ToLower(string(data[:]))
//	}
//
//	var tableName string
//	ref := reflect.ValueOf(mdl)
//	fn := ref.MethodByName("TableName")
//	if fn.Kind() == reflect.Func {
//		vals := fn.Call(nil)
//		if len(vals) > 0 {
//			tableName = vals[0].Interface().(string)
//		}
//	}
//
//	if tableName == "" {
//		strs := strings.Split(reflect.TypeOf(mdl).String(), ".")
//		tableName = snakeString(strs[len(strs)-1:][0])
//	}
//
//	query := new(DBQuery)
//    query._tblName = tableName
//	query.UpdateDB(db)
//	query._mdl = mdl
//	query.reset()
//	return query
//}
//
//func (this *DBQuery) reset() {
//	this.scopes = make([]func(db *gorm.DB) *gorm.DB, 0)
//	this.preloads = make([]Preload, 0)
//	this.hasDeleted = false
//	this.defaultPageNo = 1
//	this.defaultPageSize = 15
//}
//
//func (this *DBQuery) GetDB() *gorm.DB {
//	return this.DB
//}
//
//func (this *DBQuery) UpdateDB(db *gorm.DB) *DBQuery {
//	this.DB = db.Table(this._tblName)
//	return this
//}
//
//func (this *DBQuery) AddScope(scope ...func(db *gorm.DB) *gorm.DB) *DBQuery {
//    this.DB.Scopes(scope...)
//	return this
//}
//
//func (this *DBQuery) HasDeleted() *DBQuery {
//	this.hasDeleted = true
//	return this
//}
//
//func (this *DBQuery) AddPreload(preloads ...Preload) *DBQuery {
//    for _,p:=range preloads{
//        this.DB.Preload(p.Name, func(db *gorm.DB)*gorm.DB {
//            return db.Where(p.Condition)
//        })
//    }
//	//this.preloads = append(this.preloads, preloads...)
//	return this
//}
//
//func (this *DBQuery) prepare() *gorm.DB {
//    db:=this.DB.Begin()
//    db.Table(this._tblName)
//    if this.hasDeleted {
//        db.Unscoped()
//    }
//
//    db.Scopes(this.scopes...)
//    if len(this.preloads) > 0 {
//        for _,preload:=range this.preloads{
//            db.Preload(preload.Name, func(db2 *gorm.DB)*gorm.DB {
//                return db2.Where(preload.Condition).Order(preload.OrderBy)
//            })
//        }
//    }
//
//    return db
//}
//
//func (this *DBQuery) getScopes() []func(db *gorm.DB) *gorm.DB {
//    if len(this.preloads) > 0 {
//        for _, preload := range this.preloads {
//            this.scopes = append(this.scopes, func(db *gorm.DB) *gorm.DB {
//                return db.Preload(preload.Name, func(db2 *gorm.DB) *gorm.DB {
//                    if preload.Limit > 0 {
//                        return db2.Where(preload.Condition).Order(preload.OrderBy).Limit(preload.Limit)
//                    } else {
//                        return db2.Where(preload.Condition).Order(preload.OrderBy)
//                    }
//                })
//            })
//        }
//    }
//
//	if this.hasDeleted {
//		this.scopes = append(this.scopes, func(db *gorm.DB) *gorm.DB {
//			return db.Unscoped()
//		})
//	}
//
//	return this.scopes
//}
//
//func (this *DBQuery) One(ret interface{}) (found bool, err error) {
//	defer this.reset()
//
//	obj := this.getNewModel()
//	err = this.DB.First(obj).Error
//	//err = this.prepare().First(obj).Error
//	if err != nil {
//		if gorm.IsRecordNotFoundError(err) {
//			err = nil
//		}
//		return
//	}
//
//	err = utils.NewConvert(obj).Bind(ret)
//	if err != nil {
//		return
//	}
//
//	found = true
//	return
//}
//
//func (this *DBQuery) All(ret interface{}) (err error) {
//	defer this.reset()
//	list := this.getNewModelSlice()
//	err = this.DB.Table(this._tblName).Scopes(this.getScopes()...).Find(list).Error
//	if err != nil {
//		return
//	}
//
//	err = utils.NewConvert(list).Bind(ret)
//	return
//}
//
//func (this *DBQuery) Count() (int, error) {
//	defer this.reset()
//	return this.count()
//}
//
//func (this *DBQuery) count() (int, error) {
//	var totalCount int
//	err := this.DB.Table(this._tblName).Scopes(this.getScopes()...).Count(&totalCount).Error
//	return totalCount, err
//}
//
//func (this *DBQuery) AllByPage(ret interface{}, pageSize, pageNo int) (totalCount int, totalPage int, err error) {
//	defer this.reset()
//
//	totalCount, err = this.count()
//	if err != nil || totalCount == 0 {
//		return
//	}
//
//	if pageSize <= 0 {
//		pageSize = this.defaultPageSize
//	}
//	if pageNo <= 0 {
//		pageNo = this.defaultPageNo
//	}
//
//	offset := (pageNo - 1) * pageSize
//
//	totalPage = int(math.Ceil(float64(totalCount) / float64(pageSize)))
//	if pageNo > totalPage {
//		return
//	}
//
//	this.scopes = append(this.scopes, func(db *gorm.DB) *gorm.DB {
//		return db.Limit(pageSize).Offset(offset)
//	})
//
//	list := this.getNewModelSlice()
//
//	if err = this.DB.Table(this._tblName).Scopes(this.getScopes()...).Find(list).Error; err != nil {
//		return
//	}
//
//	err = utils.NewConvert(list).Bind(ret)
//	return
//}
//
//// 获取新的struct，返回值 *struct{}
//func (this *DBQuery) getNewModel() interface{} {
//	t := reflect.TypeOf(this._mdl)
//	m := t.Elem()
//	return reflect.Indirect(reflect.New(m)).Addr().Interface()
//}
//
//// 获取新的struct切片，返回值 *[]*struct{}
//func (this *DBQuery) getNewModelSlice() interface{} {
//	t := reflect.TypeOf(this._mdl)
//	list := reflect.New(reflect.SliceOf(t)).Elem()
//	list.Set(reflect.MakeSlice(list.Type(), 0, 0))
//	return reflect.Indirect(list).Addr().Interface()
//}
