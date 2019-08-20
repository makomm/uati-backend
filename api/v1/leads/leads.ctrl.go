package leads

import (
	"log"
	"net/http"
	"strconv"

	"gitlab.com/codenation-squad-1/backend/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
)

// Lead struct
type Lead struct {
	ID   primitive.ObjectID `bson:"_id"`
	Date string             `json:"date"`
	Sent bool               `json:"sent"`
}

// LeadDetail struct
type LeadDetail struct {
	ID         primitive.ObjectID `bson:"_id"`
	Date       string             `json:"date"`
	Sent       bool               `json:"sent"`
	Clientes   []Client           `json:"clients"`
	Recipients []Recipient        `json:"recipients"`
}

// Client struct
type Client struct {
	Nome        string  `json:"name"`
	Cargo       string  `json:"occupation"`
	Orgao       string  `json:"organization"`
	Remuneracao float64 `json:"wage"`
}

// Recipient struct
type Recipient struct {
	Username string `json:"email"`
}

const leadsCollection = "mail-leads"

func getLeads(c *gin.Context) {
	page, err := strconv.ParseInt(c.Param("page"), 10, 0)
	if err != nil {
		c.String(http.StatusBadRequest, "O número da página deve ser um número válido")
		return
	}

	limit, err := strconv.ParseInt("10", 10, 64)

	skip := 10 * page

	findOptions := options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
		Sort:  bson.D{{"date", 1}},
	}

	var collection = database.GetCollection(leadsCollection)
	var results []*Lead

	query, err := collection.Find(database.Context, bson.D{{}}, &findOptions)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao tentar consultar a base de dados")
		log.Println(err)
		return
	}

	for query.Next(database.Context) {
		var elem Lead
		err := query.Decode(&elem)
		if err != nil {
			log.Println(err)
		}

		results = append(results, &elem)
	}
	if err := query.Err(); err != nil {
		log.Println(err)
	}

	countOptions := options.CountOptions{}
	count, err := collection.CountDocuments(database.Context, bson.D{{}}, &countOptions)

	c.JSON(200, map[string]interface{}{
		"leads": results,
		"count": count,
	})
}

func getLead(c *gin.Context) {
	id := c.Param("id")

	var collection = database.GetCollection(leadsCollection)
	var result LeadDetail

	objID, _ := primitive.ObjectIDFromHex(id)
	err := collection.FindOne(database.Context, bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao tentar consultar a base de dados")
		log.Println(err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"lead": result,
	})
}
