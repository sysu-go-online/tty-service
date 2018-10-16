package controller

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

// MysqlEngine is mysql engine
var MysqlEngine *xorm.Engine

// AuthRedisClient is redis client storing invalid auth information
var AuthRedisClient *redis.Client

// DomainNameRedisClient is redis client storing used domain name
var DomainNameRedisClient *redis.Client

func init() {
	// init mysql connection
	DBADDRESS := os.Getenv("DATABASE_ADDRESS")
	if len(DBADDRESS) == 0 {
		DBADDRESS = "localhost"
	}
	DBPORT := os.Getenv("DATABASE_PORT")
	if len(DBPORT) != 0 && DBPORT[0] != ':' {
		DBPORT = ":" + DBPORT
	}
	url := fmt.Sprintf("root:root@tcp(%s%s)/mydb?charset=utf8", DBADDRESS, DBPORT)
	var err error
	engine, err := xorm.NewEngine("mysql", url)
	if err != nil {
		fmt.Println(err)
		return
	}
	MysqlEngine = engine
	if os.Getenv("DEVELOP") == "TRUE" {
		MysqlEngine.Ping()
		MysqlEngine.ShowSQL(true)
		MysqlEngine.Logger().SetLevel(core.LOG_DEBUG)
	}

	// init redis connection
	REDISADDRESS := os.Getenv("REDIS_ADDRESS")
	if len(REDISADDRESS) == 0 {
		REDISADDRESS = "localhost"
	}
	REDISPORT := os.Getenv("REDIS_PORT")
	if len(REDISPORT) == 0 {
		REDISPORT = "6379"
	}
	if len(REDISPORT) != 0 && REDISPORT[0] != ':' {
		REDISPORT = ":" + REDISPORT
	}
	client := redis.NewClient(&redis.Options{
		Addr:     REDISADDRESS + REDISPORT,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	AuthRedisClient = client

	client = redis.NewClient(&redis.Options{
		Addr:     REDISADDRESS + REDISPORT,
		Password: "", // no password set
		DB:       1,  // use domain name DB
	})
	DomainNameRedisClient = client

	if os.Getenv("DEBUG") == "TRUE" {
		pong, err := client.Ping().Result()
		fmt.Println(pong, err)
	}
}
