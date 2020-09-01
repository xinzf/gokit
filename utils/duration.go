package utils

import (
	"fmt"
	"github.com/hako/durafmt"
	"strings"
	"time"
)

type Duration struct {
	d time.Duration
}

func ParseDuration(d time.Duration) Duration {
	return Duration{d: d}
}

func (static Duration) String(limit int) string {
	if static.d == 0 {
		return ""
	}

	if static.d < time.Minute {
		return "小于1分钟"
	}

	duration := durafmt.Parse(static.d)
	str := duration.LimitFirstN(limit).String()
	fmt.Println(str)

	durationMp := make([]map[string]string, 0)
	{
		durationMp = append(durationMp, map[string]string{
			"microseconds": "微秒",
		}, map[string]string{
			"milliseconds": "毫秒",
		}, map[string]string{
			"seconds": "秒",
		}, map[string]string{
			"minutes": "分钟",
		}, map[string]string{
			"hours": "小时",
		}, map[string]string{
			"days": "天",
		}, map[string]string{
			"weeks": "星期",
		}, map[string]string{
			"months": "月",
		}, map[string]string{
			"years": "年",
		}, map[string]string{
			"hour": "小时",
		}, map[string]string{
			"day": "天",
		}, map[string]string{
			"week": "星期",
		}, map[string]string{
			"month": "月",
		}, map[string]string{
			"second": "秒",
		}, map[string]string{
			"microsecond": "微秒",
		}, map[string]string{
			"millisecond": "毫秒",
		})
	}

	//str = strings.Replace(str, "s", "", -1)

	for _, mp := range durationMp {
		for src, dst := range mp {
			str = strings.Replace(str, src, dst, -1)
		}
	}
	str = strings.Replace(str, " ", "", -1)

	return str
}

func (static Duration) Duration() time.Duration {
	return static.d
}

func (static Duration) Seconds() int64 {
	return int64(static.d / 1000 / 1000 / 1000)
}
