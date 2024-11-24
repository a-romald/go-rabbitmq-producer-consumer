package broker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/a-romald/go-rabbitmq-producer-consumer/models"
	"github.com/a-romald/go-rabbitmq-producer-consumer/utils"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"words_ex",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload models.Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [words_ex, %s]\n", q.Name)
	<-forever

	return nil
}

func handlePayload(payload models.Payload) {
	switch payload.Action {
	case "reverse":
		reverse := reverseWord(payload.Word)

		wrapper := make(map[string]interface{})
		wrapper["word"] = payload.Word
		wrapper["reverse"] = reverse

		jsonData, err := json.Marshal(wrapper)
		if err != nil {
			fmt.Println(err)
		}

		workerURL := "http://worker:" + os.Getenv("WORKER_PORT") + "/messages"

		request, err := http.NewRequest("POST", workerURL, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println(err)
		}

		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{}

		response, err := client.Do(request)
		if err != nil {
			log.Println(err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusAccepted {
			log.Println(err)
		}

	case "translate":
		// translate word

	default:
		//
	}
}

func reverseWord(word string) string {
	reverseWord := utils.Reverse(word)
	return reverseWord
}
