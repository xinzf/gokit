package main

import (
	"github.com/xinzf/gokit/logger"
	"github.com/xinzf/gokit/registry"
)

func main() {
	watcher := registry.NewConsulWatcher("127.0.0.1:8500", "gateway")
	//watcher.Add("ding.api.litudai.com", func(result *registry.Result) {
	//	logger.DefaultLogger.Debug(result.Path, result.Event, result.String())
	//})
	watcher.Add("ding.api.litudai.com/paas.enums.list/enums", func(result *registry.Result) {
		//if result.Event == "delete" {
		//watcher.Remove(result.Path)
		//}else{
		logger.DefaultLogger.Debug(result.Path, result.Event, result.String())
		//}
	})
	watcher.Run()
}
