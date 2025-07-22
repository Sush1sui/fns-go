package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sush1sui/fns-go/internal/bot"
	"github.com/Sush1sui/fns-go/internal/bot/helpers"
	"github.com/Sush1sui/fns-go/internal/common"
	"github.com/Sush1sui/fns-go/internal/config"
	"github.com/Sush1sui/fns-go/internal/repository"
	"github.com/Sush1sui/fns-go/internal/repository/mongodb"
	"github.com/Sush1sui/fns-go/internal/server/routes"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	// Initialize global variables
	common.InitializeGlobalVars()


	mongoClient := config.MongoConnection() // Initialize MongoDB connection
	defer mongoClient.Disconnect(context.Background())
	if err := mongoClient.Ping(context.Background(), nil); err != nil {
		panic(fmt.Sprintf("Failed to connect to MongoDB: %v", err))
	}

	// get collection
	stickyCollection := mongoClient.Database(os.Getenv("MONGODB_NAME")).Collection(os.Getenv("MONGODB_STICKY_CHANNELS_COLLECTION"))
	nicknameRequestCollection := mongoClient.Database(os.Getenv("MONGODB_NAME")).Collection(os.Getenv("MONGODB_NICKNAME_REQUESTS_COLLECTION"))
	exemptedCollection := mongoClient.Database(os.Getenv("MONGODB_NAME")).Collection(os.Getenv("MONGODB_EXEMPTED_USERS_COLLECTION"))

	repository.StickyService = repository.StickyServiceType{
		DBClient: &mongodb.MongoClient{
			Client: stickyCollection,
		},
	}
	repository.NicknameRequestService = repository.NicknameRequestServiceType{
		DBClient: &mongodb.MongoClient{
			Client: nicknameRequestCollection,
		},
	}
	repository.ExemptedService = repository.ExemptedServiceType{
		DBClient: &mongodb.MongoClient{
			Client: exemptedCollection,
		},
	}

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	router := routes.NewRouter()
	fmt.Printf("Server listening on Port:%s\n", cfg.ServerPort)

	// Run HTTP server in a goroutine
	go func() {
		if err := http.ListenAndServe(addr, router); err != nil {
			// Log error instead of panicking to avoid crashing the service
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	// Run Discord bot in a goroutine
	go bot.StartBot()

	// Run PingServerLoop in a goroutine
	go helpers.PingServerLoop(cfg.ServerURL)

	// Block main goroutine until interrupt signal (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	fmt.Println("Shutting down...")
}