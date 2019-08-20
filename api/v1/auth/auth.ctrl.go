package auth

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	"gitlab.com/codenation-squad-1/backend/database"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type CreateBody struct {
	Username string `json:"username" binding:"required"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Position string `json:"position"`
	Access   string `json:"access"`
	Token    string
	Password string
}

const userCollection = "usuarios"

func login(c *gin.Context) {
	var body RequestBody
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatus(400)
		return
	}
	username := body.Username
	password := fmt.Sprintf("%x", sha256.Sum256([]byte(body.Password)))

	user := getUserFromDatabase(username)
	if user.Password != password {
		c.AbortWithStatus(401)
		return
	}

	expTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("codenation-squad1"))
	c.JSON(200, map[string]interface{}{
		"user":  username,
		"token": tokenString,
	})
}

func create(c *gin.Context) {
	var collection = database.GetCollection(userCollection)
	var body CreateBody

	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatus(400)
		return
	}

	body.Token = fmt.Sprintf("%x", sha256.Sum256([]byte(body.Username+body.Name+body.Surname+time.Now().String())))
	_, err := collection.InsertOne(context.TODO(), body)
	if err != nil {
		log.Fatal(err)
	}

	SendCreatePassword(body.Username, body.Token, body.Name)
	c.JSON(200, map[string]interface{}{
		"message": "Usuário " + body.Username + " adicionado com sucesso!",
	})
}

func list(c *gin.Context) {
	c.Status(204)
}

func read(c *gin.Context) {
	c.Status(200)
}

func remove(c *gin.Context) {
	c.Status(204)
}

func update(c *gin.Context) {
	c.Status(200)
}

func getUserFromDatabase(username string) CreateBody {
	collection := database.GetCollection(userCollection)
	var res = CreateBody{}
	filter := bson.M{"username": username}
	err := collection.FindOne(database.Context, filter).Decode(&res)
	if err != nil {
		log.Println(err)
	}
	return res
}

func passwordCreation(c *gin.Context) {
	var body CreateBody
	collection := database.GetCollection(userCollection)

	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatus(400)
		return
	}

	user := getUserFromDatabase(body.Username)
	token := body.Token
	if token == "" || user.Token != token {
		c.AbortWithStatus(403)
		return
	}
	if user.Token == token {
		password := fmt.Sprintf("%x", sha256.Sum256([]byte(body.Password)))
		user.Password = password
		user.Token = ""
		_, err := collection.UpdateOne(database.Context, bson.M{"username": user.Username}, bson.M{"$set": user})
		fmt.Println(err)
		if err != nil {
			c.AbortWithStatus(400)
			return
		}
		c.JSON(200, user)
	}

}

func passwordReset(c *gin.Context) {
	var body CreateBody
	collection := database.GetCollection(userCollection)

	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatus(400)
		return
	}

	user := getUserFromDatabase(body.Username)
	user.Token = fmt.Sprintf("%x", sha256.Sum256([]byte(body.Username+body.Name+body.Surname+time.Now().String())))
	_, err := collection.UpdateOne(database.Context, bson.M{"username": user.Username}, bson.M{"$set": user})
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(400)
		return
	}
	SendCreatePassword(user.Username, user.Token, user.Name)
	c.Status(200)

}
