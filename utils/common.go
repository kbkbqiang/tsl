package utils

import (
	"github.com/satori/go.uuid"
	"strconv"
	"strings"
	"time"
	"github.com/Dark86Chen/tsl/utils/EAS"
	"log"
	"fmt"
)

// 时间转换 设置时区 东巴区
var cstZone = time.FixedZone("CST", 8*3600)

func init()  {
	time.Local = cstZone
}


// 获取当前时间戳到秒
func GetNowTimeStamp() int {
	t := time.Now().In(cstZone)
	nowTime := strconv.FormatInt(t.UTC().UnixNano(), 10)
	nowTime = nowTime[:10]
	timeStamp, _ := strconv.Atoi(nowTime)
	return timeStamp
}

// 获取当前日期格式
func GetNowTimeDate() string{
	t := time.Now().In(cstZone)
	return t.Format("2006-01-02 15:04:05")
}

func GetTimeStampByDate(date string) int64 {
	t,err := time.ParseInLocation("2006-01-02 15:04:05", date, cstZone)
	if err != nil {
		return 0
	}
	return t.Unix()
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

// 时间戳转换日期
func GetDateByTimeStamp(timeStamp int64) (date string, err error){
	secondTimeStamp := strconv.FormatInt(timeStamp, 10)
	i, err := strconv.ParseInt(secondTimeStamp[:10],10, 64)
	if err != nil {
		return date, err
	}
	t := time.Unix(i, 0).In(cstZone)
	return t.Format("2006-01-02 15:04:05"), nil
}

// 生成uuid
func GetUuid() string {
	u,_ := uuid.NewV1()
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

// 10进制转16进制
func DecHex(n int64) string {
	if n < 0 {
		log.Println("Decimal to hexadecimal error: the argument must be greater than zero.")
		return ""
	}
	if n == 0 {
		return "0"
	}
	hex := map[int64]int64{10: 65, 11: 66, 12: 67, 13: 68, 14: 69, 15: 70}
	s := ""
	for q := n; q > 0; q = q / 16 {
		m := q % 16
		if m > 9 && m < 16 {
			m = hex[m]
			s = fmt.Sprintf("%v%v", string(m), s)
			continue
		}
		s = fmt.Sprintf("%v%v", m, s)
	}
	return s
}

func GetDateFormat(timeStamp int64,formatString string) (date string, err error) {
	secondTimeStamp := strconv.FormatInt(timeStamp, 10)
	i, err := strconv.ParseInt(secondTimeStamp[:10],10, 64)
	if err != nil {
		return date, err
	}
	t := time.Unix(i, 0).In(cstZone)
	switch formatString {
	case "YYYYMMDD":
		return t.Format("20060102"), nil
	}
	return t.Format("2006-01-02 15:04:05"), nil
}