package server

import (
	"context"
	"fmt"
	"github.com/json-iterator/go"
	"reflect"
	"strings"
)

type Method struct {
	hdlName    string
	methodName string
	hdl        reflect.Value
	in         struct {
		num    int
		ctxTpy reflect.Type
		reqTpy reflect.Type
		rspTpy reflect.Type
	}
	out reflect.Type
}

func (this *Method) Structure() string {
	//return this.hdl.Type().String()
	structure := strings.TrimLeft(this.hdl.Type().String(), "func")
	return fmt.Sprintf("func (*%s) %s%s", this.hdlName, this.methodName, structure)
}

func (this *Method) Name() string {
	return this.methodName
}

func newMethod(hdlName, methodName string, methodVal reflect.Value) (*Method, error) {
	m := &Method{methodName: methodName, hdlName: hdlName}

	numIn := methodVal.Type().NumIn()
	numOut := methodVal.Type().NumOut()

	if numIn > 3 || numIn < 2 {
		return &Method{}, fmt.Errorf("%s.%s args num error", hdlName, methodName)
	}

	if numOut != 1 {
		return &Method{}, fmt.Errorf("%s.%s return num error", hdlName, methodName)
	}

	m.in.num = numIn
	if m.in.num == 2 {
		m.in.reqTpy = methodVal.Type().In(0)
		m.in.rspTpy = methodVal.Type().In(1)
	} else if m.in.num == 3 {
		m.in.ctxTpy = methodVal.Type().In(0)
		m.in.reqTpy = methodVal.Type().In(1)
		m.in.rspTpy = methodVal.Type().In(2)
	}

	m.out = methodVal.Type().Out(0)

	if m.in.reqTpy.Kind() != reflect.Interface && m.in.reqTpy.Kind() != reflect.Ptr {
		return &Method{}, fmt.Errorf("the request argument of %s.%s is not a pointer", hdlName, methodName)
	}

	if m.in.rspTpy.Kind() != reflect.Ptr {
		return &Method{}, fmt.Errorf("the response argument of %s.%s is not a pointer", hdlName, methodName)
	}

	if m.in.num == 3 && m.in.ctxTpy.String() != "context.Context" {
		return &Method{}, fmt.Errorf("the context argument of %s.%s is not context.Context", hdlName, methodName)
	}

	if m.out.String() != "error" {
		return &Method{}, fmt.Errorf("the return argument of %s.%s is not error type", hdlName, methodName)
	}

	m.hdl = methodVal

	return m, nil
}

func (this *Method) Call(ctx context.Context, rawData []byte) (interface{}, error) {

	var reqVal reflect.Value
	if this.in.reqTpy.Kind() == reflect.Interface {
		reqVal = reflect.New(this.in.reqTpy)
	} else {
		reqVal = reflect.New(this.in.reqTpy.Elem())
	}

	if len(rawData) > 0 {
		if err := jsoniter.Unmarshal(rawData, reqVal.Interface()); err != nil {
			return nil, err
		}
	}

	var rspVal reflect.Value
	//if this.in.rspTpy.Kind() == reflect.Interface {
	//    rspVal = reflect.New(this.in.rspTpy)
	//}else{
	rspVal = reflect.New(this.in.rspTpy.Elem())
	//}
	values := make([]reflect.Value, 0)
	if this.in.num == 2 {
		values = this.hdl.Call([]reflect.Value{
			reqVal,
			rspVal,
		})
	} else {
		values = this.hdl.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reqVal,
			rspVal,
		})
	}

	er, _ := values[0].Interface().(error)
	rsp := rspVal.Convert(this.in.rspTpy)
	return rsp.Interface(), er
}
