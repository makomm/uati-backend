package main

import (
	"log"
	"fmt"
	"github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/bson"
	"context"
	"leads/database"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/streadway/amqp"
	"time"
	"os"
	"encoding/json"
	"strconv"
	
)
//Recipient struct for usuarios of UATI system
type Recipient struct {
	Email string `bson:"username"`
}

//Cliente struct for UATI client
type Cliente struct {
	Nome string `bson:"nome"`
	Lead bool `bson:"lead"`
}
//Funcionario struct for GOV SP funcionarios
type Funcionario struct {
	Nome string `bson:"nome" json:"nome"`
	Cargo string `bson:"cargo" json:"cargo"`
	Orgado string `bson:"orgao" json:"orgao"`
	Remuneracao float64 `bson:"remuneracao" json:"remuneracao"`
}
//Mail struct to store and handle leads
type Mail struct {
	ClientesList []Funcionario `bson:"clientes" json:"clientes"`
	Recipients []Recipient `bson:"recipients" json:"clientes"`
	Date string `bson:"date" json:"date"`
	Sent bool `bson:"sent" json:"sent"`
}
//Email struct
type Email struct {
	Recipient string
	Template  string
	Login     string
	Senha     string
	Subject   string
}

var limit = 2000

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	database.Initialize()
	queueListener()
}

func searchLeads(monthYear string) []Funcionario{
	var clientes []string
	var result []Funcionario
	var colFuncionarios = database.GetCollection("funcionarios-" + monthYear)
	var colClientes = database.GetCollection("clientes-uati")
	findOptions := options.Find()
	cur,err:=colClientes.Find(context.TODO(), bson.M{"lead":bson.M{"$ne":true}},findOptions)
	if err!= nil {
		fmt.Println(err)
	}
	for cur.Next(context.TODO()){
		var cliente Cliente
		cur.Decode(&cliente)
		clientes = append(clientes,cliente.Nome)
	}
	// Teste clients
	// clientes = append(clientes,"AAIRON TELES DE CAMARGO")
	// clientes = append(clientes,"ABARE VAZ DE LIMA")
	// clientes = append(clientes,"ADAIR RIBEIRO JUNIOR")

	cur2,err := colFuncionarios.Find(context.TODO(),bson.M{"nome":bson.M{"$in":clientes},"remuneracao":bson.M{"$gt":limit}})

	for cur2.Next(context.TODO()){
		var funcionario Funcionario
		cur2.Decode(&funcionario)
		cliente := Cliente{funcionario.Nome,true}
		_, err := colClientes.UpdateOne(database.Context, bson.M{"nome": funcionario.Nome}, bson.M{"$set": cliente})
		if err!=nil{
			fmt.Println(err)
		}
		result = append(result,funcionario)
	}

	return result
}

//Generates a Mail List
func storeMailList(list []Funcionario) {
	var usersCollection = database.GetCollection("usuarios")
	var mailCollection = database.GetCollection("mail-leads")
	var recipients []Recipient
	var template = ""
	usrList, err := usersCollection.Find(context.TODO(), bson.D{{}})
	if err!= nil {
		fmt.Println(err)
	}
	for _, funcionario := range list{
		template = template + funcionario.Nome + ": " + strconv.FormatFloat(funcionario.Remuneracao, 'E', -1, 64) +"\n"
	}
	for usrList.Next(context.TODO()){
		var recipient Recipient
		usrList.Decode(&recipient)
		msgMail(recipient.Email, template)
		recipients = append(recipients,recipient)
	}
	var mail = Mail{list,recipients,time.Now().String(),true}
	_,errMail := mailCollection.InsertOne(context.TODO(), mail)
	if errMail != nil {
		fmt.Println(errMail)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
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
		"lead-check", // name
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
			var leads = searchLeads(string(d.Body))
			if len(leads) > 0 {
				storeMailList(leads)
			}
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func msgMail(username string, template string) {
	email := os.Getenv("EMAIL_LOGIN")
	senha := os.Getenv("EMAIL_SENHA")
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

	body, _ := json.Marshal(Email{username, template, email, senha, "Leads Encontrados!"})
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

