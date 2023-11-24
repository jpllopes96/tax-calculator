package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"tax-calculator/internal/infra/database"
	"tax-calculator/internal/usecase"
	"tax-calculator/pkg/rabbitmq"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "root"
		dbname   = "golang_training"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	//wait all  to run and then close,
	defer db.Close()
	orderRepository := database.NewOrderRepository(db)

	uc := usecase.NewCalculateFinalPrice(orderRepository)

	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	msgRabbitmqChannel := make(chan amqp.Delivery)

	go rabbitmq.Consume(ch, msgRabbitmqChannel) // listening queue // stuck // T2

	rabbitmqWorker(msgRabbitmqChannel, uc) // T1

	// input := usecase.OrderInput{
	// 	ID:    "1213",
	// 	Price: 10.0,
	// 	Tax:   1.0,
	// }

	// output, err := uc.Execute(input)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(output)
}

func rabbitmqWorker(msgChan chan amqp.Delivery, uc *usecase.CalculateFinalPrice) {
	fmt.Println("Starting rabbitmq")
	for msg := range msgChan {
		var input usecase.OrderInput
		err := json.Unmarshal(msg.Body, &input)

		if err != nil {
			panic(err)
		}

		output, err := uc.Execute(input)
		if err != nil {
			panic(err)
		}

		msg.Ack(false)

		fmt.Println("nessage processed and saved in BDD: ", output)
	}
}
