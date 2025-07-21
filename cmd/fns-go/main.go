package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Sush1sui/fns-go/internal/bot"
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

	mongoClient := config.MongoConnection() // Initialize MongoDB connection
	defer mongoClient.Disconnect(context.Background())
	if err := mongoClient.Ping(context.Background(), nil); err != nil {
		panic(fmt.Sprintf("Failed to connect to MongoDB: %v", err))
	}

	// get collection
	stickyCollection := mongoClient.Database(os.Getenv("MONGODB_NAME")).Collection(os.Getenv("MONGODB_STICKY_CHANNELS_COLLECTION"))

	repository.StickyService = repository.StickyServiceType{
		DBClient: mongodb.MongoClient{
			Client: stickyCollection,
		},
	}

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	router := routes.NewRouter()
	fmt.Printf("Server listening on Port:%s\n", cfg.ServerPort)
	bot.StartBot()
	err = http.ListenAndServe(addr, router)
	if err != nil {
		panic(err)
	}
}