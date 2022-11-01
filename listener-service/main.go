package main

import (
	"listener/events"
	"log"
	"os"
	"time"

	ampq "github.com/rabbitmq/amqp091-go"
)

func main(){
	//try to connect to rabbitMQ
	client,err:=connect()
	if err !=nil{
		log.Printf("error connecting to mq %s",err)
		os.Exit(1)
	}

	consumer,err:=events.NewConsumer(client)
	if err!=nil{
		panic(err)
	}
	err=consumer.Listen([]string{"log.INFO","log.WARNING,log.ERROR"})
	if err!=nil{
		panic(err)
	}
	defer client.Close()
	log.Println("Connected with rabbtMQ!")

}

func connect()(connection *ampq.Connection,err error){
	var count int64
	var backOff = 1*time.Second

	for {
		c,err:=ampq.Dial("amqp://guest:guest@rabbitmq")
		if err!=nil{
			log.Printf("RabbitMQ is not ready %s \n",err)
			count++
		}else{
			connection=c
			log.Println("Connected with rabbtMQ!")
			break
		}
		if count>5{
			log.Panicf("Can't connect to rabbit mq %s\n",err)
			return nil,err
		}

		backOff=time.Duration(count*2)
		time.Sleep(backOff)
		continue

	}

	return connection,nil

}