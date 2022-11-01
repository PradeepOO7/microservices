package events

import (
	"log"

	ampq "github.com/rabbitmq/amqp091-go"
)

type Emitter struct{
	conn *ampq.Connection
}

func (e *Emitter) Setup() error {
	channel, err := e.conn.Channel()
	if err != nil {
		log.Printf("Error in events.Setup %s\n", err)
		return err
	}
	return declareExchange(channel)
}

func (e *Emitter) Push(event,severity string )(error) {
	channel, err := e.conn.Channel()
	if err != nil {
		log.Printf("Error in events.Setup %s\n", err)
		return err
	}
	defer channel.Close()
	err=channel.Publish(
		"logs_topic",
		severity,
		false,
		false,
		ampq.Publishing{
			ContentType: "text/plain",
			Body: []byte(event),
		},
	)

	if err != nil {
		log.Printf("Error in events.Setup %s\n", err)
		return err
	}
	return nil

}

func NewEventEmitter(conn *ampq.Connection)(Emitter,error){
	emitter:=Emitter{
		conn: conn,
	}

	err:=emitter.Setup()
	if err != nil {
		log.Printf("Error in events.Setup %s\n", err)
		return Emitter{},err
	}
	return emitter,nil
}
