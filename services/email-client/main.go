package main

import (
	"log"
	"strings"

	"encoding/json"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	queueListener()
}

//Email struct to send emails
type Email struct {
	Recipient string
	Template  string
	Login     string
	Senha     string
	Subject   string
}

func sendMail(recipient string, template string, login string, senha string, subject string) {
	// Set up authentication information.
	auth := sasl.NewPlainClient("", login+"@gmail.com", senha)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{recipient}
	// We need to get the subject dynamically
	msg := strings.NewReader("To: " + recipient + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + template +
		"\r\n")
	err := smtp.SendMail("smtp.gmail.com:587", auth, login+"@gmail.com", to, msg)
	if err != nil {
		log.Fatal(err)
	}
}

func queueListener() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbit:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"mail", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var rmail Email
			json.Unmarshal(d.Body, &rmail)
			sendMail(rmail.Recipient, rmail.Template, rmail.Login, rmail.Senha, rmail.Subject)
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
