package main

import (
	"embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/a-romald/go-rabbitmq-producer-consumer/broker"
	"github.com/a-romald/go-rabbitmq-producer-consumer/models"
	"github.com/a-romald/go-rabbitmq-producer-consumer/utils"
)

//go:embed templates
var templatesFS embed.FS // import embed

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(templatesFS, "templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	data := make(map[string]string)
	data["Port"] = os.Getenv("WORKER_PORT")
	tmpl.Execute(w, data)
}

func (app *Application) HandleForm(w http.ResponseWriter, r *http.Request) {
	var payload models.Payload
	_ = json.NewDecoder(r.Body).Decode(&payload)
	defer r.Body.Close()

	err := app.pushToQueue(payload.Word, payload.Action)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	utils.PrintJSON(w, http.StatusOK, payload, "payload")
}

// pushToQueue pushes a message into RabbitMQ
func (app *Application) pushToQueue(word, action string) error {
	producer, err := broker.NewDataEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := models.Payload{
		Word:   word,
		Action: action,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = producer.Push(string(j), "word.reverse") // data & routing key
	if err != nil {
		return err
	}
	return nil
}
