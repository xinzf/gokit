package registry

import "fmt"

type Service struct {
	ID      string
	Name    string
	Address string
	Port    int
	Metas   map[string]string
}

func (this *Service) Addr() string {
	return fmt.Sprintf("http://%s:%d", this.Address, this.Port)
}

func (this *Service) Method(hdlName, methodName string) string {
	return fmt.Sprintf("%s?_service=%s&_method=%s", this.Addr(), hdlName, methodName)
}
