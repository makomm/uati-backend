package auth

import (
	"encoding/json"
	"html"
	"log"
	"os"
	"strings"

	"gitlab.com/codenation-squad-1/backend/database"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type Email struct {
	Recipient string
	Template  string
	Login     string
	Senha     string
	Subject   string
}
type Template struct {
	Name string `bson:"name"`
	Html string `bson:"html"`
}

func SendCreatePassword(username string, token string, name string) {
	url := os.Getenv("URL_APP")
	email := os.Getenv("EMAIL_LOGIN")
	senha := os.Getenv("EMAIL_SENHA")
	urlRabbit := os.Getenv("URL_RABBIT")
	conn, err := amqp.Dial(urlRabbit)
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

	template := templateToken(name, token, url, username)
	body, _ := json.Marshal(Email{username, template, email, senha, "Registro de senha"})
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}

func templateToken(name string, token string, url string, email string) string {
	var template Template
	colTemplate := database.GetCollection("templates")
	colTemplate.FindOne(database.Context, bson.M{"name": "create-password"}).Decode(&template)
	var body = strings.ReplaceAll(strings.ReplaceAll(html.UnescapeString(template.Html), "[#NAME#]", name), "[#LINK#]", url+"/create-password?token="+token+"&user="+email)
	return body
}
