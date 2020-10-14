package registry

import (
	"fmt"
	"strings"
)

type Service struct {
	ID      string
	Name    string
	Address string
	Port    int
	Metas   map[string]string
}

func (this *Service) Domain() string {
	return fmt.Sprintf("http://%s:%d", this.Address, this.Port)
}

func (this *Service) URL(hdlName string) string {
	var (
		hdl, method string
	)
	strs := strings.Split(hdlName, ".")
	if len(strs) > 0 {
		method = strs[len(strs)-1:][0]
		hdl = strings.Join(strs[:len(strs)-1], ".")
	}
	return fmt.Sprintf("%s?_service=%s&_method=%s", this.Domain(), hdl, method)
}
