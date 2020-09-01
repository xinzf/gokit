package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"time"
	// "os"
	"strconv"
	"strings"

	"github.com/cstockton/go-conv"
	"github.com/fatih/structs"
	"github.com/satori/go.uuid"
	"unicode"
)

func GetIP() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return ips, err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}

	return ips, nil
}

func GetPort() int {
	l, _ := net.Listen("tcp", ":0")
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

func Ip2num(ip string) int {
	canSplit := func(c rune) bool { return c == '.' }
	lisit := strings.FieldsFunc(ip, canSplit) //[58 215 20 30]
	//fmt.Println(lisit)
	ip1_str_int, _ := strconv.Atoi(lisit[0])
	ip2_str_int, _ := strconv.Atoi(lisit[1])
	ip3_str_int, _ := strconv.Atoi(lisit[2])
	ip4_str_int, _ := strconv.Atoi(lisit[3])
	return ip1_str_int<<24 | ip2_str_int<<16 | ip3_str_int<<8 | ip4_str_int
}

func Num2ip(num int) string {
	ip1_int := (num & 0xff000000) >> 24
	ip2_int := (num & 0x00ff0000) >> 16
	ip3_int := (num & 0x0000ff00) >> 8
	ip4_int := num & 0x000000ff
	//fmt.Println(ip1_int)
	data := fmt.Sprintf("%d.%d.%d.%d", ip1_int, ip2_int, ip3_int, ip4_int)
	return data
}

func NowTime() string {
	return TimeFormat(time.Now())
}

func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func UUID() string {
	return strings.Replace(uuid.NewV4().String(), "-", "", -1)
}

// StructToMap 转换struct为map
func StructToMap(s interface{}) map[string]interface{} {
	return structs.Map(s)
}

// StringToInt 字符串转数值
func StringToInt(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func UnixMilli(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

func UnixMilliToTime(mse int64) time.Time {
	str, _ := conv.String(mse)
	var (
		s  int64
		ms int64
	)
	s, _ = conv.Int64(str[:10])
	if len(str) >= 13 {
		ms, _ = conv.Int64(str[10:])
	}
	return time.Unix(s, ms)
}

func IsChinese(str string) bool {
	var count int
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			count++
			break
		}
	}
	return count > 0
}

func MD5(str string, toUpper ...bool) string {
	ret := md5.Sum([]byte(str))
	s := hex.EncodeToString(ret[:])

	var upper bool
	if len(toUpper) > 0 {
		upper = toUpper[0]
	}
	if upper {
		s = strings.ToUpper(s)
	}
	return s
}
