package main

import (
	"fmt"
	"log"
	"net/http"
	//"os"
	"time"

	ampq "github.com/rabbitmq/amqp091-go"
)

const webPort = "8080"

type Config struct{
	Rabbit *ampq.Connection
}

func main() {
		//try to connect to rabbitMQ
		client,err:=connect()
		if err !=nil{
			log.Printf("error connecting to mq %s\n",err)
			//os.Exit(1)
		}
	app := Config{
		Rabbit: client,
	}
	log.Printf("Starting broker service on port %s\n", webPort)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Panicf("Something went wrong %s", err)
	}
	
}

func connect()(connection *ampq.Connection,err error){
	var count int64
	var backOff = 1*time.Second

	for {
		c,err:=ampq.Dial("amqp://guest:guest@rabbitmq")
		if err!=nil{
			log.Printf("RabbitMQ is not ready s %s \n",err)
			count++
		}else{
			connection=c
			log.Println("Connected with rabbitMQ!")
			break
		}
		if count>5{
			log.Panicf("Can't connect to rabbit mq %s\n",err)
			return nil,err
		}

		backOff=time.Duration(count*2)
		log.Printf("backoff for %v seconds\n",backOff)
		time.Sleep(backOff)
		continue

	}

	return connection,nil

}
