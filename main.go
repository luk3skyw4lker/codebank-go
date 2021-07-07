package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/luk3skyw4lker/codebank-go/infrastructure/repositories"
	"github.com/luk3skyw4lker/codebank-go/tasks"
)

func main() {
	db := setupDatabase()
	defer db.Close()
	// producer := setupKafkaProducer()
	setupTransactionTask(db)
	// serveGrpc(processTransactionUseCase)
}

func setupDatabase() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("host"),
		os.Getenv("port"),
		os.Getenv("user"),
		os.Getenv("password"),
		os.Getenv("dbname"),
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("error connection to database")
	}
	return db
}

func setupTransactionTask(db *sql.DB) tasks.TransactionTask {
	transactionRepository := repositories.NewTransactionRepository(db)
	task := tasks.NewTransactionTask(transactionRepository)
	// task.KafkaProducer = producer

	return task
}
