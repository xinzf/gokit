package main

import (
	"context"
	"fmt"
	"github.com/xinzf/gokit/logger"
	"github.com/xinzf/gokit/registry"
	"github.com/xinzf/gokit/server"
)

// Define project name
const PROJECT_NAME = "Test"

type Request struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Response struct {
	Err      string `json:"err"`
	FileName string `json:"file_name"`
}

type Handler struct {
}

func (this *Handler) Image(ctx context.Context, req *server.UploadFile, rsp *Response) error {
	f, err := req.Get("upfile")
	if err != nil {
		return err
	}
	if err := f.SaveFile("./haha.jpg"); err != nil {
		return err
	}
	rsp.FileName = f.FileName()
	return nil
}

func main() {
	logger.DefaultLogger.Init(logger.ProjectName(PROJECT_NAME))

	server.New(
		server.Name(PROJECT_NAME),
		server.Port(6060),
		server.Registry(registry.NewConsul(registry.Addrs("127.0.0.1:8500"))),
	)

	server.Register("uploader", new(Handler))

	// You can set up several callback functions before the server starting.
	server.Run(context.Background(), func() error {
		fmt.Println("execute before server start.")
		return nil
	})
}
