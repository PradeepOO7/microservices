package main

import (
	"context"
	"fmt"
	"log"
	"logger/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "8080"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "5001"
)

type Config struct {
	Models data.Models
}

var client *mongo.Client

func main() {

	mongoClient, err := connectToMongoDb()
	if err != nil {
		log.Panicf("Error connection to mongodb %s", err)
	}
	client = mongoClient

	ctx,cancel:=context.WithTimeout(context.Background(),15*time.Second)

	defer cancel()

	defer func(){
		if err=client.Disconnect(ctx);err!=nil{
			log.Panicf("Database Disconnected %s\n",err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}
	log.Printf("Starting Logger service on port %s\n", webPort)

	app.serve()


}

func (app *Config) serve(){
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Panicf("Something went wrong %s", err)
	}

}

func connectToMongoDb() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)

	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Panicf("Something went wrong %s", err)
		return nil, err
	}else{
		log.Println("Conected to mongoDB")
	}
	return c, nil
}
