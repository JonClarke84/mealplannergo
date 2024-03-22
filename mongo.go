package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func db(uri string) (*mongo.Client, error) {
	// connect to MongoDB
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return client, nil
}

type MealList struct {
	SaturdayMeal  string
	SundayMeal    string
	MondayMeal    string
	TuesdayMeal   string
	WednesdayMeal string
	ThursdayMeal  string
	FridayMeal    string
}

func getThisWeeksMeals(client *mongo.Client) (MealList, error) {
	// find the first document in the collection
	collection := client.Database("GoShopping").Collection("shopping-lists")
	filter := bson.D{{}}
	var mealList MealList
	err := collection.FindOne(context.Background(), filter).Decode(&mealList)
	if err != nil {
		fmt.Printf("Error finding document: %s\n", err)
		return mealList, err
	}
	return mealList, nil
}
