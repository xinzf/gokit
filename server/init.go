package server

import (
    "os"
)

var (
	Server *server
	//Config   Cfg
	localIP  string
	hostName string
	//debug        = kingpin.Flag("debug", "Enable debug mode.").Bool()
	//confFilePath = kingpin.Flag("config", "Provide a valid configuration path").Short('c').Default("./conf/").ExistingFileOrDir()
)

func init() {
	Server = new(server)
	//ips, err := nettools.IntranetIP()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//if len(ips) == 0 {
	//	log.Fatalln("cant't get local ip")
	//}
	//localIP = ips[0]

	hostName, _ = os.Hostname()
}
