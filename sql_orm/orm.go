package sql_orm

import (
	"github.com/go-xorm/xorm"
	"time"
	"github.com/Dark86Chen/tsl/log"
	_ "github.com/go-sql-driver/mysql"
)


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
	engine.TZLocation, err = time.LoadLocation(e.Location)

	e.State = true

	if err != nil {
		log.Logger.Warning("set orm engine location err --> ", err.Error())
	}

	return engine, nil
}