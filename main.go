package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"goapi/controllers"
)

const itemIDPath = "/items/:id"

func main() {
	gin.SetMode(gin.ReleaseMode)

	// Connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Error creating MongoDB client:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	collection := client.Database("GOtestdb").Collection("items")

	r := gin.Default()

	r.GET("/items", controllers.GetItems(collection))
	r.POST("/items", controllers.CreateItem(collection))
	r.GET(itemIDPath, controllers.GetItem(collection))
	r.PUT(itemIDPath, controllers.UpdateItem(collection))
	r.DELETE(itemIDPath, controllers.DeleteItem(collection))

	r.Run(":8080")
}
