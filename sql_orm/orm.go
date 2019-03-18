package sql_orm

import (
	"github.com/go-xorm/xorm"
	"time"
	"github.com/Dark86Chen/tsl/log"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"github.com/pkg/errors"
)

var cstZone = time.FixedZone("CST", 8*3600)

func (e *Engine)GetOrmEngine() (engine *xorm.Engine, err error) {
	if EngineCon.Engine != nil {
		if err := e.Engine.Ping(); err != nil {
			// 关闭原来的链接
			e.Engine.Close()

			engine, err := e.createEngine()

			if err != nil {
				log.Logger.Error("create engine err --> ", err.Error())
				return nil, err
			}
			e.Engine = engine
		}
	} else {
		engine, err := e.createEngine()

		if err != nil {
			log.Logger.Error("create init engine err --> ", err.Error())
			return nil, err
		}
		e.Engine = engine
	}

	return e.Engine, nil
}


func (e *Engine)createEngine() (engine *xorm.Engine, err error) {
	engine, err = xorm.NewEngine(DriverName, DataSourceName)
	if err != nil {
		return nil, err
	}
	engine.ShowSQL(true)
	engine.SetMaxOpenConns(10)

	engine.SetMaxOpenConns(e.MaxOpenConns)
	engine.SetMaxIdleConns(e.MaxIdleConns)

	// 设置时区
	engine.TZLocation = cstZone //time.LoadLocation(e.Location)

	e.State = true

	if err != nil {
		log.Logger.Warning("set orm engine location err --> ", err.Error())
	}

	return engine, nil
}


func (s *ShortEngine)GetShortEngine() (engine *xorm.Engine, err error) {
	ShortDataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		s.User, s.Pwd, s.Host,
		s.Port, s.DbName, s.Charset)
	engine, err = xorm.NewEngine(s.DriverName, ShortDataSourceName)
	if err != nil {
		return nil, err
	}

	pingState := make(chan bool)

	go func() {
		if err := engine.Ping(); err != nil {
			log.Logger.Error("connection db error --> ", err.Error())
		}
		pingState <- true
	}()

	time.AfterFunc(5 * time.Second, func() {
		pingState <- false
	})

	select {
		case state := <-pingState:
			if state == false {
				return nil, errors.New("connection db error")
			} else {
				goto END
			}
	}

	END:
	engine.ShowSQL(true)
	// 设置时区
	engine.TZLocation = cstZone //time.LoadLocation(e.Location)

	if err != nil {
		log.Logger.Warning("set orm engine location err --> ", err.Error())
	}

	return engine, nil
}