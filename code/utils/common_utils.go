package utils

import (
	"time"
)

func GetCurrentDateAsString() string {
	return time.Now().Format("2006-01-02")

	// 本地测试可以用这个。将1天缩短到10秒。
	//return strconv.Itoa((time.Now().Second() + 100000) / 10)
}
