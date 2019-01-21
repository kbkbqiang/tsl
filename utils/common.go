package utils

import (
	"github.com/satori/go.uuid"
	"strconv"
	"strings"
	"time"
	"github.com/Dark86Chen/tsl/utils/EAS"
)

// 时间转换 设置时区 东巴区
var cstZone = time.FixedZone("CST", 8*3600)

// 获取当前时间戳到秒
func GetNowTimeStamp() int {
	t := time.Now().In(cstZone)
	nowTime := strconv.FormatInt(t.UTC().UnixNano(), 10)
	nowTime = nowTime[:10]
	timeStamp, _ := strconv.Atoi(nowTime)
	return timeStamp
}

func GetNowTimeDate() string{
	t := time.Now().In(cstZone)
	return t.Format("2006-01-02 15:04:05")
}

// 获取当前时间戳到毫秒
func GetNowMillisecondTimeStamp() int64 {
	t := time.Now().In(cstZone)
	nowTime := strconv.FormatInt(t.UTC().UnixNano(), 10)
	nowTime = nowTime[:13]
	timeStamp, _ := strconv.Atoi(nowTime)
	return int64(timeStamp)
}

// 获取当前时间戳到分钟
func GetNowMinutTimeStamp() int64 {
	t := time.Now().In(cstZone)
	nowTime := strconv.FormatInt(t.UTC().UnixNano(), 10)
	nowTime = nowTime[:7]
	timeStamp, _ := strconv.Atoi(nowTime)
	return int64(timeStamp)
}


// 生成uuid
func GetUuid() string {
	u := uuid.NewV1()
	uid := u.String()
	return strings.Replace(uid, "-", "", -1)
}

// 生成token
func GenerateToken(tokenByte []byte) (token string, err error) {
	token, err = EAS.Encrypt(tokenByte)
	if err != nil {
		return token, err
	}
	return token, nil
}