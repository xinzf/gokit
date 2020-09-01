package gokit

//
//import (
//	"errors"
//	"fmt"
//	"github.com/spf13/viper"
//	"reflect"
//)
//
//type Cfg struct {
//	Server   ServerConfig   `mapstructure:"server"`
//	Registry RegisterConfig `mapstructure:"register"`
//	Db       *dbConfig      `mapstructure:"db"`
//	Redis    *redisConfig   `mapstructure:"redis"`
//	Mongo    *mongoConfig   `mapstructure:"mongo"`
//	Logger   LoggerConfig   `mapstructure:"logger"`
//	Mode     string         `mapstrucure:"mod"`
//}
//
//type ServerConfig struct {
//	Name string `mapstructure:"name"`
//	Host string `mapstructure:"host"`
//	Port int    `mapstructure:"port"`
//}
//
//type RegisterConfig struct {
//	Addr string `mapstructure:"addr"`
//}
//
//func (s ServerConfig) String() string {
//	return fmt.Sprintf("%s:%d", s.Host, s.Port)
//}
//
//func (s ServerConfig) Addr() string {
//	return fmt.Sprintf("http://%s:%d", s.Host, s.Port)
//}
//
//
//
//
//
//func ReadConfig(fileName string, dstConf interface{}, initFn ...func(interface{}) error) error {
//	viper.SetConfigName(fileName)
//	viper.AddConfigPath(*confFilePath)
//	viper.SetConfigType("yaml")
//
//	if err := viper.ReadInConfig(); err != nil {
//		return err
//	}
//
//	refType := reflect.TypeOf(dstConf)
//	if refType.Kind() != reflect.Ptr {
//		return errors.New("config's argument: dstConf is not a pointer")
//	}
//
//	if err := viper.Unmarshal(dstConf); err != nil {
//		return err
//	}
//
//	for _, fn := range initFn {
//		if err := fn(dstConf); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
