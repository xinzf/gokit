package client

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/cstockton/go-conv"
	"github.com/json-iterator/go"
	reg "github.com/xinzf/gokit/registry"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	callServiceErr = func(svcName, msg string) error {
		return fmt.Errorf("call service: '%s' failed, err: %s", svcName, msg)
	}

	registry reg.Registry
	client   *http.Client
	once     sync.Once
)

func New(r reg.Registry) {
	once.Do(func() {

		r.Watch()
		registry = r

		var transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: func(network, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(network, addr, 2*time.Second)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(3 * time.Second))
				return conn, nil
			},
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
		runtime.SetFinalizer(&transport, func(tr **http.Transport) {
			(*tr).CloseIdleConnections()
		})

		client = &http.Client{
			Transport: transport,
		}

	})
}

func Call(serviceName string, body interface{}, rsp interface{}, headers ...map[string]interface{}) (err error) {
	if registry == nil {
		return errors.New("client have not registry")
	}

	strs := strings.Split(serviceName, ".")
	if len(strs) < 3 {
		return callServiceErr(serviceName, "service name is wrong")
	}

	svcName := strs[0]
	hdlName := strings.Join(strs[1:len(strs)-1], ".")
	methodName := strs[len(strs)-1:][0]

	svc, found := registry.Service(svcName)
	if !found {
		return fmt.Errorf("Not found service: %s\n", serviceName)
	}

	address := svc.Method(hdlName, methodName)
	data, err := jsoniter.Marshal(body)
	if err != nil {
		return callServiceErr(serviceName, err.Error())
	}

	var reqBody io.Reader
	if len(data) > 0 {
		reqBody = bytes.NewBuffer(data)
	}
	req, _ := http.NewRequest("post", address, reqBody)

	mp := make(map[string]string)
	if len(headers) > 0 {
		for k, v := range headers[0] {
			vv, _ := conv.String(v)
			if vv != "" {
				mp[k] = vv
			}
		}
	}

	hds := parseHeaders(mp)
	for k, v := range hds {
		req.Header.Add(k, v[0])
		req.Header.Set(k, v[0])
	}

	response, err := client.Do(req)
	if err != nil {
		return callServiceErr(serviceName, err.Error())
	}

	if response.StatusCode != 200 {
		return callServiceErr(serviceName, fmt.Sprintf("status code: %d", response.StatusCode))
	}

	defer response.Body.Close()
	respBody, _ := ioutil.ReadAll(response.Body)
	if err = jsoniter.Unmarshal(respBody, rsp); err != nil {
		return callServiceErr(serviceName, err.Error())
	}
	return nil
}

// 解析请求头
func parseHeaders(header map[string]string) http.Header {
	h := http.Header{}
	h.Set("Accept-Encoding", "gzip")

	if header != nil {

		for key, values := range header {
			h.Add(key, values)
		}
	}

	h.Add("Connection", "Close")

	_, hasAccept := h["Accept"]
	if !hasAccept {
		h.Add("Accept", "*/*")
	}
	_, hasAgent := h["User-Agent"]
	if !hasAgent {
		h.Add("User-Agent", "gokit-client/v1.0")
	}
	_, hasContentType := h["Content-Type"]
	if !hasContentType {
		h.Add("Content-Type", "application/json")
	}
	return h
}
