/*
 @Time : 2020/8/15 1:22 PM
 @Author : chenye
 @File : init
 @Software : GoLand
 @Remark : 
*/

package mongo_orm

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
	"strconv"
	"sync"
	"context"
)

var (
	PoolCon  = new(Pool)
)

type Pool struct {
	Client		 *mongo.Client
	ClientStatus bool
	Ctx     context.Context
	lock    sync.Mutex
	PoolAuth options.Credential
	PoolClientOptions *options.ClientOptions
	KeepLiveTicker *time.Ticker
	KeepLiveTime  time.Duration
}

func init()  {
	PoolCon.PoolClientOptions = options.Client()
	// init config
	PoolCon.PoolAuth.Username = os.Getenv("MONGO_USER")
	PoolCon.PoolAuth.Password = os.Getenv("MONGO_PWD")

	if os.Getenv("MONGO_CONNECT_TIMEOUT") == "" {
		PoolCon.PoolClientOptions.SetServerSelectionTimeout( 5 * time.Second)
		PoolCon.PoolClientOptions.SetSocketTimeout(5 * time.Second)
		PoolCon.PoolClientOptions.SetConnectTimeout(5 * time.Second)
	} else {
		timeOut,_ := strconv.Atoi(os.Getenv("MONGO_CONNECT_TIMEOUT"))
		PoolCon.PoolClientOptions.SetServerSelectionTimeout(time.Duration(timeOut) * time.Second)
		PoolCon.PoolClientOptions.SetSocketTimeout(time.Duration(timeOut) * time.Second)
		PoolCon.PoolClientOptions.SetConnectTimeout(time.Duration(timeOut) * time.Second)
	}

	if os.Getenv("MONGO_MAX_POOL_SIZE") == "" {
		PoolCon.PoolClientOptions.SetMaxPoolSize(5000)
	} else {
		maxPoolSize,_ := strconv.ParseUint(os.Getenv("MONGO_MAX_POOL_SIZE"), 10, 64)
		PoolCon.PoolClientOptions.SetMaxPoolSize(maxPoolSize)
	}

	PoolCon.PoolClientOptions.SetHosts([]string{"localhost:12000"})
	PoolCon.PoolClientOptions.SetAuth(PoolCon.PoolAuth)
	PoolCon.ClientStatus = false
}


