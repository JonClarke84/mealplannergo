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
	ShoppingList  []string
}

func getNewestList(client *mongo.Client) (MealList, error) {
	// find the first document in the collection
	collection := client.Database("GoShopping").Collection("shopping-lists")
	filter := bson.D{{}}
	var mealList MealList
	err := collection.FindOne(context.Background(), filter).Decode(&mealList)
	if err != nil {
		fmt.Printf("Error finding meal list: %s\n", err)
		return mealList, err
	}
	return mealList, nil
}

func getShoppingList(client *mongo.Client) ([]string, error) {
	// find the first document in the collection
	collection := client.Database("GoShopping").Collection("shopping-lists")
	filter := bson.D{{}}
	var mealList MealList
	err := collection.FindOne(context.Background(), filter).Decode(&mealList)
	if err != nil {
		fmt.Printf("Error finding shopping list: %s\n", err)
		return mealList.ShoppingList, err
	}
	return mealList.ShoppingList, nil
}

func addShoppingListItem(client *mongo.Client, item string) error {
	// add an item to the shopping list
	collection := client.Database("GoShopping").Collection("shopping-lists")
	filter := bson.D{{}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "ShoppingList", Value: item}}}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Printf("Error adding shopping list item: %s\n", err)
		return err
	}
	return nil
}

func updateMeal(client *mongo.Client, day string, meal string) error {
	// update the document
	collection := client.Database("GoShopping").Collection("shopping-lists")
	filter := bson.D{{}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: day, Value: meal}}}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Printf("Error updating meal: %s\n", err)
		return err
	}
	return nil
}

func deleteShoppingListItem(client *mongo.Client, item string) error {
	// delete an item from the shopping list
	collection := client.Database("GoShopping").Collection("shopping-lists")
	filter := bson.D{{}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "ShoppingList", Value: item}}}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Printf("Error deleting shopping list item: %s\n", err)
		return err
	}
	return nil
}
