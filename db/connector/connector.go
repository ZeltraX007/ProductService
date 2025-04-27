package connector

import (
	"ProductService/config"
	"ProductService/db"
)

var (
	PGDBConnector  db.DBOperations
	RedisConnector db.CacheInterface
)

func Connector() {
	PGDBConnector = db.NewPGConnector(config.PostgresConn)
	RedisConnector = db.NewRedisConnector(config.RedisClient)
}
