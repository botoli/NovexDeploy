package main

import (
	"context"
	"localVercel/db"
	"localVercel/internal/deployer"
	"localVercel/internal/queue"
	"localVercel/internal/worker"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if err := db.InitDB(); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	q := queue.NewRedisQueue(redisAddr, "", 0, "build_queue")
	d := deployer.NewDeployer()
	
	w := worker.NewWorker(q, d)

	log.Println("Starting worker...")
	w.Start(context.Background())
}
