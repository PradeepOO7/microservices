package events

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	ampq "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *ampq.Connection
	queueName string
}

var wg sync.WaitGroup

func NewConsumer(conn *ampq.Connection) (Consumer, error) {
	consmer := Consumer{
		conn: conn,
	}
	err := consmer.Setup()
	if err != nil {
		log.Printf("Error in events.NewConsumer %s\n", err)
		return Consumer{}, err
	}
	return consmer, nil
}

func (cons *Consumer) Setup() error {
	channel, err := cons.conn.Channel()
	if err != nil {
		log.Printf("Error in events.Setup %s\n", err)
		return err
	}
	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

func (cons *Consumer) Listen(topics []string) error {
	channel, err := cons.conn.Channel()
	if err != nil {
		log.Printf("Error in events.Setup %s\n", err)
		return err
	}
	defer channel.Close()
	queue, err := declareRandomQueue(channel)

	for _, s := range topics {
		channel.QueueBind(
			queue.Name,
			s,
			"logs_topic",
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	messages, err := channel.Consume(queue.Name, "", true, false, false, false, nil)

	//forever := make(chan bool)

	wg.Add(2)
	go func() {
		defer wg.Done()
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)
			go handlePayload(payload)
		}
	}()

	log.Printf("Waiting for messages in queue %s\n", queue.Name)
	wg.Wait()
	// var s=<-forever
	// log.Printf("forever channel %v \n",s)
	return nil

}

func handlePayload(payload Payload) error {
	defer wg.Done()
	switch payload.Name {
	case "log", "event":
		err := logEvent(payload)
		if err != nil {
			return err
		}
	}
	return nil
}

func logEvent(logPayload Payload) error {
	jsonData, _ := json.MarshalIndent(logPayload, "", "\t")
	log.Panicln("We are in logEvent Function")
	request, err := http.NewRequest("POST", "http://logger-service:8080/log", bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return err
	}
	
	return nil

}
