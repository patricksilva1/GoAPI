package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const invalidIdFormat = "Invalid ID format"

type Item struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
}

// RETRIEVE ALL ITEMS
func GetItems(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var items []Item

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cur, err := collection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cur.Close(ctx)

		for cur.Next(ctx) {
			var item Item
			err := cur.Decode(&item)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			items = append(items, item)
		}

		c.JSON(http.StatusOK, items)
	}
}

// RETRIEVE ITEM BY ID
func GetItem(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var item Item
		id := c.Param("id")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": invalidIdFormat})
			return
		}

		err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&item)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": invalidIdFormat})
			return
		}

		c.JSON(http.StatusOK, item)
	}
}

// CREATE ITEM
func CreateItem(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var item Item
		if err := c.BindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := collection.InsertOne(ctx, item)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		insertedID, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting inserted ID"})
			return
		}

		item.ID = insertedID

		c.JSON(http.StatusOK, item)
	}
}

// UPDATE ITEM
func UpdateItem(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var item Item
		id := c.Param("id")

		if err := c.BindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": invalidIdFormat})
			return
		}

		_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": item})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Item updated successfully"})
	}
}

// DELETE ITEM BY ID
func DeleteItem(collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": invalidIdFormat})
			return
		}

		_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
	}
}
