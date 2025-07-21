package config

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// Config holds the application's configuration
type Config struct {
	DiscordToken string
	ServerPort   string
	AppID        string
}

// New loads configuration from environment variables
func New() (*Config, error) {
	// In development, load .env file. In production, env vars are set directly.
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	token := os.Getenv("bot_token")
	if token == "" {
		return nil, fmt.Errorf("bot_token environment variable not set")
	}

	appID := os.Getenv("app_id")
	if appID == "" {
		return nil, fmt.Errorf("app_id environment variable not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5469" // Default port
	}

	return &Config{
		DiscordToken: token,
		AppID:        appID,
		ServerPort:   port,
	}, nil
}

func MongoConnection() *mongo.Client {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
  client, err := mongo.Connect(opts)
  if err != nil {
    panic(err)
  }

  // Send a ping to confirm a successful connection
  if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
    panic(err)
  }
  fmt.Println("DB Connected!")

	return client
}