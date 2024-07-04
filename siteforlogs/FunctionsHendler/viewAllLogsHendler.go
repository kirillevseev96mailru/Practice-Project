package functionshendler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Log struct {
	Date    string `json:"date"`
	User    string `json:"user"`
	Command string `json:"command"`
	Type    string `json:"type"`
}

type Show struct {
	Pages []Log
}


func outPutAllLogs() ([]Log, error) {
	file, err := ioutil.ReadFile("New-logs.json")
	if err != nil {
		return nil, err
	}

	var logs []Log
	err = json.Unmarshal(file, &logs)
	if err != nil {
		log.Fatal(err)
	}


	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongodb:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	collection := client.Database("logsdb").Collection("logs")

	// Опции поиска (можно задать фильтры и сортировку)
	findOptions := options.Find()

	// Находим все документы в коллекции
	cur, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	defer cur.Close(ctx)

	var results []bson.M
	if err := cur.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	var newlogs []Log
	for _, result := range results {
		fmt.Println(result)
		newlog := Log{
			Date: result["date"].(string),
                	User: result["user"].(string),
                	Command: result["command"].(string),
                	Type: result["typeourlog"].(string),
        	}
		newlogs = append(newlogs, newlog)
	}
	return newlogs, nil
}


func viewAllLogsHendler(w http.ResponseWriter, r *http.Request) {
	logs, err := outPutAllLogs()
	if err != nil {
		fmt.Println(err)
	}

	web_data := Show{logs}

	tmpl, err := template.ParseFiles("templates/viewAllLogs.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	tmpl.ExecuteTemplate(w, "viewAllLogs", web_data)
}
