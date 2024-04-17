package main

import (
	"context"
	"fmt"
  "time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShoppingListDocument struct {
  ID primitive.ObjectID `bson:"_id" json:"id,omitempty"`
  ShoppingList []ShoppingListItem `bson:"ShoppingList" json:"ShoppingList"`
}

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

func getShoppingList(client *mongo.Client) ([]ShoppingListItem, error) {
  collection := client.Database("GoShopping").Collection("shopping-lists")

  var document struct {
    ID           primitive.ObjectID `bson:"_id"`
    ShoppingList []struct {
      ID     primitive.ObjectID `bson:"Id"`
      Item   string             `bson:"Item"`
      Ticked bool               `bson:"Ticked"`
    } `bson:"ShoppingList"`
  }

  err := collection.FindOne(context.Background(), bson.D{}).Decode(&document)
  if err != nil {
    return nil, err
  }

  shoppingList := make([]ShoppingListItem, len(document.ShoppingList))
  for i, item := range document.ShoppingList {
    shoppingList[i] = ShoppingListItem{
      ID:     item.ID,
      IDHex:  item.ID.Hex(),
      Item:   item.Item,
      Ticked: item.Ticked,
    }
  }

  return shoppingList, nil
}

func updateMeal(client *mongo.Client, day string, meal string) error {
	collection := client.Database("GoShopping").Collection("meal-plans")
	filter := bson.D{{Key: "meals", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{Key: "day", Value: day}}}}}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "meals.$.meal", Value: meal}}}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Printf("Error updating meal: %s\n", err)
	}
	return nil
}

func addShoppingListItem(client *mongo.Client, itemName string) error {
  fmt.Println("adding item: ", itemName)
  newIdFromTimeStamp := primitive.NewObjectIDFromTimestamp(time.Now())
  fmt.Println("new id: ", newIdFromTimeStamp)
  collection := client.Database("GoShopping").Collection("shopping-lists") 
  filter := bson.D{{}}
  // push { id: newIdFromTimeStamp, Item: itemName, Ticked: false } to the ShoppingList array
  update := bson.D{{Key: "$push", Value: bson.D{{Key: "ShoppingList", Value: bson.D{{Key: "Id", Value: newIdFromTimeStamp}, {Key: "Item", Value: itemName}, {Key: "Ticked", Value: false}}}}}}
  _, err := collection.UpdateOne(context.Background(), filter, update)
  if err != nil {
    fmt.Printf("Error adding shopping list item: %s\n", err)
    return err
  }
  return nil
}

func updateShoppingListItem(client *mongo.Client, oldItem string, newItem string) error {
	collection := client.Database("GoShopping").Collection("shopping-lists")
	filter := bson.D{{Key: "ShoppingList", Value: oldItem}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "ShoppingList.$", Value: newItem}}}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Printf("Error updating shopping list item: %s\n", err)
		return err
	}
	return nil
}

func deleteShoppingListItem(client *mongo.Client, item string) error {
	objectIdFromHex, err := primitive.ObjectIDFromHex(item)
	if err != nil {
		fmt.Printf("Error converting item to object ID: %s\n", err)
		return err
	}
	collection := client.Database("GoShopping").Collection("shopping-list-items")
	filter := bson.D{{Key: "_id", Value: objectIdFromHex}}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		fmt.Printf("Error deleting shopping list item: %s\n", err)
		return err
	}
	return nil
}

func sortShoppingList(client *mongo.Client, newItemOrder []string) error {
	collection := client.Database("GoShopping").Collection("shopping-lists")
	filter := bson.D{{}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "ShoppingList", Value: newItemOrder}}}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Printf("Error sorting shopping list: %s\n", err)
		return err
	}
	return nil
}

func getMealPlan(client *mongo.Client) (MealPlan, error) {
	// find the first document in the collection
	collection := client.Database("GoShopping").Collection("meal-plans")
	// get the first document
	filter := bson.D{{}}
	var mealPlan MealPlan
	err := collection.FindOne(context.Background(), filter).Decode(&mealPlan)
	if err != nil {
		fmt.Printf("Error finding meal plan: %s\n", err)
		return mealPlan, err
	}
	return mealPlan, nil
}
