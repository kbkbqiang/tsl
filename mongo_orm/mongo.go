/*
 @Time : 2020/8/15 1:19 PM
 @Author : chenye
 @File : mongo
 @Software : GoLand
 @Remark : make mongo client pool
*/

package mongo_orm

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"context"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
	"github.com/Dark86Chen/tsl/log"
)

func GetPool() (client *mongo.Client, err error) {
	if err := PoolCon.conn(PoolCon.PoolClientOptions); err != nil {
		return  nil, err
	}

	return PoolCon.Client, nil
}

func keepLive(num time.Duration) {
	PoolCon.KeepLiveTicker = time.NewTicker(num)
	defer PoolCon.KeepLiveTicker.Stop()

	for _ = range PoolCon.KeepLiveTicker.C {
		if PoolCon.ClientStatus == false {
			log.Logger.Warning("continue mongo connect")
			break
		}
		if status := PoolCon.ping(); !status {
			log.Logger.Warning("mongo is coon lost")
			break
		} else {
			log.Logger.Info("mongo status is live ")
		}
	}
	GetPool()
}

func (p *Pool)conn(options *options.ClientOptions) (err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.ClientStatus == true {
		err = p.Client.Disconnect(PoolCon.Ctx)
		if err != nil {
			log.Logger.Error("disconnect mongo error --> ", err.Error())
			p.ClientStatus = false
			go keepLive(2 * time.Second)
			return err
		}
		p.ClientStatus = false
	} else {
		p.Ctx = context.Background()
	}

	p.Client, err = mongo.Connect(p.Ctx, options)
	if err != nil {
		p.ClientStatus = false
		return err
	}

	p.ClientStatus = true

	// TODO keep live
	go keepLive(2 * time.Second)

	return nil
}

func (p Pool)ping() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()

	if err := p.Client.Ping(ctx, readpref.Primary()); err != nil {
		log.Logger.Error("mongo ping error --> ", err.Error())
		return false
	}

	return true
}
