package clients

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

const MESSAGE string = "New clients imported"

func SendImportedClients() {
	urlRabbit := os.Getenv("URL_RABBIT")
	conn, err := amqp.Dial(urlRabbit)
	logOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	logOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"lead-check", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	logOnError(err, "Failed to declare a queue")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(MESSAGE),
		})
	log.Printf(" [x] Sent %s", MESSAGE)
	logOnError(err, "Failed to publish a message")
}

func logOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
