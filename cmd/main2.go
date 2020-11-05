package main

import (
	"fmt"
	"strings"
)

func main() {
	str := "Paas.form.layout.Save"
	strs := strings.Split(str, ".")
	if len(strs) == 0 {
		return
	}

	fmt.Println(strs[1:])
	method := strs[len(strs)-1:][0]
	hdl := strings.Join(strs[:len(strs)-1], ".")
	fmt.Println(hdl, "  ", method)
	//consulKV := kv.NewConfig(
	//	kv.WithPrefix("gateway"),
	//	kv.WithAddress("127.0.0.1:8500"),
	//)
	//if err := consulKV.Init(); err != nil {
	//	panic(err)
	//}
	//
	//consulKV.Watch("ding.api.litudai.com", func(result *kv.Result) {
	//	logger.DefaultLogger.Debug(result.Key(), result.String())
	//})
	//consulKV.Watch("ding.api.litudai.com/paas.enums.list", func(result *kv.Result) {
	//    logger.DefaultLogger.Debug(result.Key(), result.String())
	//})
	//
	//select {}
}
