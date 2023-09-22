package main

import (
	"log"

	"github.com/Tus1688/kim-hackathon-2023-api/authutil"
	"github.com/Tus1688/kim-hackathon-2023-api/database"
)

func main() {
	err := database.InitMysql()
	if err != nil {
		log.Fatal("unable to connect to mysql", err)
	}
	err = database.InitRedis()
	if err != nil {
		log.Fatal("unable to connect to redis", err)
	}
	err = authutil.InitializeJWTKey()
	if err != nil {
		log.Fatal("unable to initialize jwt key", err)
	}

}
