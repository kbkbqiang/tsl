package sql_orm

import (
	"fmt"
	"os"
	"strconv"
	"github.com/go-xorm/xorm"
)

type Engine struct {
	Engine *xorm.Engine
	MaxOpenConns int 	// 最大打开连接数
	MaxIdleConns int 	// 连接池的空闲数大小
	Location	 string // 时区
	State        bool   // 链接状态
}

var (
	DataSourceName string
	DriverName	string
	EngineCon = Engine{MaxOpenConns: 100, MaxIdleConns: 100, Location: "Asia/Shanghai"}
	MaxOpenConns int
	MaxIdleConns int
	Location  	 string
	Err			error
)

func init()  {
	DataSourceName = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"),
		os.Getenv("DB_POINT"), os.Getenv("DB_NAME"), os.Getenv("DB_CHARSET"))
	DriverName = os.Getenv("Driver_Name")

	MaxOpenConnsStr := os.Getenv("MaxOpenConns")
	if MaxOpenConnsStr != "" {
		if MaxOpenConns,Err = strconv.Atoi(MaxOpenConnsStr); Err == nil {
			EngineCon.MaxOpenConns = MaxOpenConns
		}
	}

	MaxIdleConnsStr := os.Getenv("MaxIdleConns")
	if MaxIdleConnsStr != "" {
		if MaxIdleConns,Err = strconv.Atoi(MaxIdleConnsStr); Err == nil {
			EngineCon.MaxOpenConns = MaxIdleConns
		}
	}

	Location = os.Getenv("Location")
	if Location != "" {
		EngineCon.Location = Location
	}

}
