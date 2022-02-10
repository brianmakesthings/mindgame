// go run amqp_server.go
// based on https://www.rabbitmq.com/tutorials/tutorial-six-go.html

package main

import (
	"encoding/json"
	"log"
	"mindgame_go/mindgame"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	var userIds = make(map[string]bool)
	var rooms = make(map[string]*(mindgame.RoomInfo))
	// var sids = make(map[string]string)
	conn, err := amqp.Dial("amqp://localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		true,        // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			println("dealing with request")
			var request map[string]string
			json.Unmarshal(d.Body, &request)
			failOnError(err, "Failed to convert body")

			log.Printf("%s", request)

			var response map[string]string
			switch request["request"] {
			case "generateUserId":
				response = mindgame.GenerateUserId(request, &userIds)
			case "createRoom":
				response = mindgame.CreateRoom(request, &rooms)
			case "joinRoom":
				response = mindgame.JoinRoom(request, &rooms)
			case "getParticipants":
				response = mindgame.GetParticipants(request, &rooms)
			case "startGame":
				response = mindgame.StartGame(request, &rooms)
			case "getGameState":
				response = mindgame.GetGameState(request, &rooms)
			case "playCard":
				response = mindgame.PlayCard(request, &rooms)
			case "removeUser":
				response = mindgame.RemoveUser(request, &rooms)
			case "deleteRoom":
				response = mindgame.DeleteRoom(request, &rooms)
			// case "addSid" :
			// 	resposne = AddSid(request, &sids)
			// case "popSid" :
			// 	resposne = PopSid(request, &sids)
			default:
				response = map[string]string{"error": "could not determine request"}
			}
			log.Printf("%s", response)
			body, _ := json.Marshal(response)

			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          body,
				})
			failOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf("Awaiting RPC requests")
	<-forever
}
