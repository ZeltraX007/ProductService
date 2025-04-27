package app

import (
	"ProductService/config"
	"ProductService/db/connector"
)

func Start() {
	config.InitializeEnv() //reading env
	config.InitDB()        //establishing db connection
	config.InitRedis()     //establishing redis connection
	connector.Connector()
	runserver()
}
