package FunctionsForParser

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
	"context"
)

type Logs struct {
	Date string `json:"date"`
	User string `json:"user"`
	Command string `json:"command"`
	TypeOurLog string `json:"type"`
}

func WriterInJson(logItems ...string) {
	if len(logItems) != 4 {
		fmt.Println("Неверное количество аргументов для logItems")
		return
	}

	logs := []Logs{
		{
                        Date: logItems[0],
                        User: logItems[1],
                        Command: logItems[2],
                        TypeOurLog: logItems[3],
                },
	}

	logEntry := Logs{
                        Date: logItems[0],
                        User: logItems[1],
                        Command: logItems[2],
                        TypeOurLog: logItems[3],
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
	doc, err := bson.Marshal(logEntry)
	if err != nil {
		log.Fatal(err)
	}

	_, err = collection.InsertOne(ctx, doc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("hi")
	fileName := "New-logs.json"

	if _, err = os.Stat(fileName); os.IsNotExist(err) {

		jsonLog, err := json.MarshalIndent(logs, "", " ")
		if err != nil {
			fmt.Println(err)
			return
		}

		err = ioutil.WriteFile(fileName, jsonLog, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("File created")
		return
	}

	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("FILE ISN'T WORKING", err)
		return
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("He can't read file", err)
		return
	}

	var existingLogs []Logs
	err = json.Unmarshal(data, &existingLogs)
	if err != nil {
		fmt.Println("JsonMap", err)
		return
	}

	existingLogs = append(existingLogs, logs...)

	newData, err := json.MarshalIndent(existingLogs, "", " ")
	if err != nil {
	fmt.Println(err)
	return
	}

	err = ioutil.WriteFile(fileName, newData, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}
