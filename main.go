package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// represents the data about a record album
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"Price"`
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	var results []album
	client := ConnectDB()
	collection := GetCollection(client)
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	if err := cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, results)
}

func postAlbums(c *gin.Context) {
	client := ConnectDB()
	collection := GetCollection(client)
	var newAlbum album

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}
	_, err := collection.InsertOne(ctx, newAlbum)

	if err != nil {
		c.IndentedJSON(http.StatusConflict, newAlbum)
	} else {
		c.IndentedJSON(http.StatusCreated, newAlbum)
	}

}
func getusingID(c *gin.Context) {
	id := c.Param("id")
	var results []album
	client := ConnectDB()
	collection := GetCollection(client)
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	if err := cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	for _, a := range results {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}

	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
func ConnectDB() *mongo.Client {
	Mongo_URL := "mongodb://127.0.0.1:27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(Mongo_URL))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	defer cancel()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to mongoDB")
	return client
}
func GetCollection(client *mongo.Client) *mongo.Collection {
	collection := client.Database("local").Collection("new")
	return collection
}
func main() {

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getusingID)
	router.POST("/album", postAlbums)
	router.Run("localhost:3000")
}
