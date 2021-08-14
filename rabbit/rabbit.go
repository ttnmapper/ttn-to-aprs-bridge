package rabbit

import (
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"ttn-to-aprs-bridge/utils"
)

type Credentials struct {
	Host     string
	Port     int
	User     string
	Password string
	Exchange string
	Queue    string
	AutoAck  bool

	MessagesChannel chan amqp.Delivery
}

func (credentials *Credentials) SubscribeToRabbit() {
	conn, err := amqp.Dial("amqp://" + credentials.User + ":" + credentials.Password + "@" + credentials.Host + ":" + strconv.Itoa(credentials.Port) + "/")
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Create a channel for errors
	notify := conn.NotifyClose(make(chan *amqp.Error)) //error channel

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		credentials.Exchange, // name
		"fanout",                            // type
		true,                                // durable
		false,                               // auto-deleted
		false,                               // internal
		false,                               // no-wait
		nil,                                 // arguments
	)
	utils.FailOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		credentials.Queue, // name
		false,                      // durable
		true,                     // delete when unused
		false,                     // exclusive
		false,                     // no-wait
		nil,                       // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		10,    // prefetch count
		0,     // prefetch size
		false, // global
	)
	utils.FailOnError(err, "Failed to set queue QoS")

	err = ch.QueueBind(
		q.Name,                              // queue name
		"",                                  // routing key
		credentials.Exchange, // exchange
		false,
		nil)
	utils.FailOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		credentials.AutoAck,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	utils.FailOnError(err, "Failed to register a consumer")

waitForMessages:
	for {
		select {
		case err := <-notify:
			if err != nil {
				log.Println(err.Error())
			}
			break waitForMessages
		case d := <-msgs:
			//log.Printf("[a] Packet received")
			credentials.MessagesChannel <- d
		}
	}

	log.Fatal("Subscribe channel closed")
}
