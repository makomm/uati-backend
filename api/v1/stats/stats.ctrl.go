package stats

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/gin-gonic/gin"
	"gitlab.com/codenation-squad-1/backend/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Funcionario struct for GOV SP funcionarios
type Funcionario struct {
	Nome        string  `bson:"nome" json:"nome"`
	Cargo       string  `bson:"cargo" json:"cargo"`
	Orgao       string  `bson:"orgao" json:"orgao"`
	Remuneracao float64 `bson:"remuneracao" json:"remuneracao"`
}

//Recipient struct for usuarios of UATI system
type Recipient struct {
	Email string `bson:"username"`
}

//Mail struct to store and handle leads
type Mail struct {
	ClientesList []Funcionario `bson:"clientes" json:"clientes"`
	Recipients   []Recipient   `bson:"recipients" json:"recipients"`
	Date         string        `bson:"date" json:"date"`
	Sent         bool          `bson:"sent" json:"sent"`
}

//CargoStats struct for data without cut in remuneracao >20k
type CargoStats struct {
	Cargo       string  `bson:"cargo" json:"cargo"`
	Mean        float64 `bson:"mean" json:"mean"`
	Std         float64 `bson:"std" json:"std"`
	Percentil75 float64 `bson:"percentil75" json:"percentil75"`
	Month       int     `bson:"month" json:"month"`
	Year        int     `bson:"year" json:"year"`
}

//OrgaoStats struct for data without cut in remuneracao >20k
type OrgaoStats struct {
	Orgao       string  `bson:"orgao" json:"orgao"`
	Mean        float64 `bson:"mean" json:"mean"`
	Std         float64 `bson:"std" json:"std"`
	Percentil75 float64 `bson:"percentil75" json:"percentil75"`
	Month       int     `bson:"month" json:"month"`
	Year        int     `bson:"year" json:"year"`
}

//TopOrgao struct for data with cut in remuneracao >20k
type TopOrgao struct {
	Orgao string `bson:"orgao" json:"orgao"`
	Total int    `bson:"total" json:"total"`
}

//TopCargo struct for data with cut in remuneracao >20k
type TopCargo struct {
	Cargo string `bson:"cargo" json:"cargo"`
	Total int    `bson:"total" json:"total"`
}

//Distribution struct for remuneracao distribution
type Distribution struct {
	Mean        float64   `bson:"mean" json:"mean"`
	Std         float64   `bson:"std" json:"std"`
	Percentil75 float64   `bson:"percentil75" json:"percentil75"`
	Month       int       `bson:"month" json:"month"`
	Year        int       `bson:"year" json:"year"`
	Bins        []float64 `bson:"bins" json:"bins"`
	Entries     []float64 `bson:"entries" json:"entries"`
}

//LeadType struct for analytics
type LeadType struct {
	Label   string  `json:"label"`
	Percent float64 `json:"percent"`
	Total   int     `json:"total"`
}

//LeadStats struct for leads analytics
type LeadStats struct {
	Alertas int64      `json:"alertas"`
	Leads   int        `json:"leads"`
	Cargos  []LeadType `json:"cargos"`
	Orgaos  []LeadType `json:"orgaos"`
	Mean    float64    `json:"mean"`
}

const leadsCollection = "mail-leads"
const statCargos = "statistic-cargos"
const statOrgao = "statistic-orgao"
const statTopCargo = "statistic-top-cargo"
const statTopOrgao = "statistic-top-orgao"
const statDistribution = "statistic-remuneracao-distribution"

func getStats(c *gin.Context) {
	var allOrgaoTop = getTop5OrgaoByMedian()
	var allCargoTop = getTop5CargoByMedian()
	var topCargoCount = getTop5CargoByCount()
	var topOrgaoCount = getTop5OrgaoByCount()
	var distributions = getDistributionRemuneracao()
	var leadStats = getLeadsInfo()
	var totalAllOrgao, _ = database.GetCollection(statOrgao).CountDocuments(context.TODO(), bson.D{{}})
	var totalAllCargos, _ = database.GetCollection(statCargos).CountDocuments(context.TODO(), bson.D{{}})
	var totalTopOrgaos, _ = database.GetCollection(statTopOrgao).CountDocuments(context.TODO(), bson.D{{}})
	var totalTopCargos, _ = database.GetCollection(statTopCargo).CountDocuments(context.TODO(), bson.D{{}})
	c.JSON(200, map[string]interface{}{
		"statistics_all": map[string]interface{}{
			"orgaos":       allOrgaoTop,
			"cargos":       allCargoTop,
			"total_orgaos": totalAllOrgao,
			"total_cargos": totalAllCargos,
		},
		"statistics_over_minimum": map[string]interface{}{
			"cargos":       topCargoCount,
			"orgaos":       topOrgaoCount,
			"total_orgaos": totalTopOrgaos,
			"total_cargos": totalTopCargos,
		},
		"distributions_remuneracao": distributions,
		"lead_stats":                leadStats,
	})
}

func getTop5OrgaoByMedian() []OrgaoStats {
	var result []OrgaoStats
	var orgaoCollection = database.GetCollection(statOrgao)
	var number int64
	number = 5
	findOptions := options.FindOptions{
		Limit: &number,
		Sort:  bson.D{{"mean", -1}},
	}
	cur, err := orgaoCollection.Find(context.TODO(), bson.D{{}}, &findOptions)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for cur.Next(context.TODO()) {
		var orgao OrgaoStats
		cur.Decode(&orgao)
		if math.IsNaN(orgao.Std) {
			orgao.Std = 0
		}
		result = append(result, orgao)
	}
	return result
}

func getTop5CargoByMedian() []CargoStats {
	var result []CargoStats
	var cargoCollection = database.GetCollection(statCargos)
	var number int64
	number = 5
	findOptions := options.FindOptions{
		Limit: &number,
		Sort:  bson.D{{"mean", -1}},
	}
	cur, err := cargoCollection.Find(context.TODO(), bson.D{{}}, &findOptions)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for cur.Next(context.TODO()) {
		var cargo CargoStats
		cur.Decode(&cargo)
		if math.IsNaN(cargo.Std) {
			cargo.Std = 0
		}
		result = append(result, cargo)
	}
	return result
}

func getTop5CargoByCount() []TopCargo {
	var result []TopCargo
	var cargoCollection = database.GetCollection(statTopCargo)
	var number int64
	number = 5
	findOptions := options.FindOptions{
		Limit: &number,
		Sort:  bson.D{{"total", -1}},
	}
	cur, err := cargoCollection.Find(context.TODO(), bson.D{{}}, &findOptions)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for cur.Next(context.TODO()) {
		var cargo TopCargo
		cur.Decode(&cargo)
		result = append(result, cargo)
	}
	return result
}

func getTop5OrgaoByCount() []TopOrgao {
	var result []TopOrgao
	var orgaoCollection = database.GetCollection(statTopOrgao)
	var number int64
	number = 5
	findOptions := options.FindOptions{
		Limit: &number,
		Sort:  bson.D{{"total", -1}},
	}
	cur, err := orgaoCollection.Find(context.TODO(), bson.D{{}}, &findOptions)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for cur.Next(context.TODO()) {
		var orgao TopOrgao
		cur.Decode(&orgao)
		result = append(result, orgao)
	}
	return result
}

func getDistributionRemuneracao() []Distribution {
	var distributions []Distribution
	var distributionCollection = database.GetCollection(statDistribution)
	cur, err := distributionCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for cur.Next(context.TODO()) {
		var distribution Distribution
		cur.Decode(&distribution)
		distributions = append(distributions, distribution)
	}
	return distributions
}

func getLeadsInfo() LeadStats {
	var leadStats LeadStats
	var mails []Mail
	var funcionarios []Funcionario
	var cargos = make(map[string]int)
	var orgaos = make(map[string]int)
	var meanRemuneracao float64
	var distributionCollection = database.GetCollection(leadsCollection)

	cur, err := distributionCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		fmt.Println(err)
		return LeadStats{}
	}
	leadStats.Alertas, _ = distributionCollection.CountDocuments(context.TODO(), bson.D{{}})

	for cur.Next(context.TODO()) {
		var mail Mail
		cur.Decode(&mail)
		mails = append(mails, mail)

		for _, funcionario := range mail.ClientesList {
			funcionarios = append(funcionarios, funcionario)
			cargos[funcionario.Cargo] = cargos[funcionario.Cargo] + 1
			orgaos[funcionario.Orgao] = orgaos[funcionario.Orgao] + 1
			meanRemuneracao = meanRemuneracao + funcionario.Remuneracao
		}

	}
	leadStats.Cargos = getTopFive(cargos)
	leadStats.Orgaos = getTopFive(orgaos)
	leadStats.Leads = len(funcionarios)
	leadStats.Mean = meanRemuneracao / float64(len(funcionarios))
	return leadStats
}

func getTopFive(entries map[string]int) []LeadType {
	type kv struct {
		Key   string
		Value int
	}
	var results []LeadType
	var ss []kv
	var total = 0
	for k, v := range entries {
		ss = append(ss, kv{k, v})
		total = total + v
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	if len(ss) > 5 {
		ss = ss[:5]
	}
	for _, val := range ss {
		results = append(results, LeadType{
			Label:   val.Key,
			Total:   val.Value,
			Percent: float64(val.Value) / float64(total),
		})
	}
	return results
}
