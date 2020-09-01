# GOKIT

1. The first need tod install gokit
```shell script
go get -u github.com/xinzf/gokit
```

2. Import it in your code:
```go
import "github.com/xinzf/gokit"
```

3. Tutorial:
```go
package main
import (
    "context"
    "fmt"
    "github.com/xinzf/gokit/logger"
    "github.com/xinzf/gokit/server"
    "github.com/xinzf/gokit/storage"
)
// Define project name 
const PROJECT_NAME = "test_gokit"

type Request struct{
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type Response struct{
    FullName string `json:"fullname"`
}

type Handler struct{
    
}

func (this *Handler) GetFullName(ctx context.Context,req *Request,rsp *Response) error {
    rsp.FullName = fmt.Sprintf("%s_%d",req.ID,req.Name)
    return nil
}

func main() {
    logger.DefaultLogger.Init(logger.ProjectName(PROJECT_NAME))
    server.New(
        server.Name(PROJECT_NAME),
        server.Port(6060),
    )

    server.Register("Handler", new(Handler))

    // You can set up several callback functions before the server starting.
    server.Run(context.Background(), func() error {
        return storage.DB.Init(
            storage.DbLogger(logger.DefaultLogger),
            storage.DbConfig("127.0.0.1:3306", "db_user", "db_pswd", "db_name"),
        )
    },func() error {
        fmt.Println("execute before server start.")
    
        // If an error occurs, the server will stop starting..
        return nil
    })
}
```

