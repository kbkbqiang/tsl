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
	"time"
	"github.com/Dark86Chen/tsl/log"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"strconv"
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
	client, err := GetPool()
	if err != nil {
		log.Logger.Errorf("connect mongo error --> %s", err.Error())
		return
	} else {
		log.Logger.Infof("reconnect mongo success!")
	}

	PoolCon.Client = client
	return
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
	keepLiveTime := os.Getenv("keepLiveSecond")
	keepLiveSecond, _ := strconv.Atoi(keepLiveTime)
	if keepLiveSecond == 0 {
		keepLiveSecond = 20
	}
	go keepLive(time.Duration(keepLiveSecond) * time.Second)

	return nil
}

func (p Pool)ping() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	status := make(chan bool, 0)
	go func() {
		defer close(status)
		if err := p.Client.Ping(ctx, readpref.Primary()); err != nil {
			log.Logger.Error("mongo ping error --> ", err.Error())
			//return false
			status <- false
			return
		}
		status <- true
	}()

	for {
		select {
		case <- ctx.Done():
			return false
		case v,_ := <- status:
			return v
		}
	}
	//return true
}
