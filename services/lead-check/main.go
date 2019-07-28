package main

import (
	"log"
	"fmt"
	"github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/bson"
	"context"
	"leads/database"

	// "encoding/hex"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"html"
	"strings"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/streadway/amqp"
	"time"
	"os"
	"encoding/json"
	
)
//Recipient struct for usuarios of UATI system
type Recipient struct {
	Email string `bson:"username"`
	Name string `bson:"name"`
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
//Template struct
type Template struct {
	Name string `bson:"name"`
	Html string `bson:"html"`
}


var limit = 20000

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	database.Initialize()
	queueListener()
}

func searchLeads() []Funcionario{
	var clientes []string
	var result []Funcionario
	var colFuncionarios = database.GetCollection("funcionarios")
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
	usrList, err := usersCollection.Find(context.TODO(), bson.D{{}})
	if err!= nil {
		fmt.Println(err)
	}
	
	for usrList.Next(context.TODO()){
		var recipient Recipient
		usrList.Decode(&recipient)
		recipients = append(recipients,recipient)
	}
	var mail = Mail{list,recipients,time.Now().String(),true}
	obj,errMail := mailCollection.InsertOne(context.TODO(), mail)
	if errMail != nil {
		fmt.Println(errMail)
	}
	oid, _ := obj.InsertedID.(primitive.ObjectID)
	for _,rec := range recipients{
		template:=templateLeads(oid.Hex(), rec.Name)
		msgMail(rec.Email, template)
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
			var leads = searchLeads()
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

func templateLeads(ID string, name string) string {
	url := os.Getenv("URL_APP")
	var template Template
	colTemplate := database.GetCollection("templates")
	colTemplate.FindOne(database.Context, bson.M{"name": "new-leads"}).Decode(&template)
	var body = strings.ReplaceAll(strings.ReplaceAll(html.UnescapeString(template.Html), "[#NAME#]", name), "[#LINK#]", url+"/lead-detail/"+ID)
	return body
}