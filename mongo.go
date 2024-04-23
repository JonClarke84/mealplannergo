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
    ShoppingList []ShoppingListItem `bson:"ShoppingList"`
  }

  err := collection.FindOne(context.Background(), bson.D{}).Decode(&document)
  if err != nil {
    return nil, err
  }

  return document.ShoppingList, nil
}

func getShoppingListItemFromIDHex(client *mongo.Client, IDHex string) (ShoppingListItem, error) {
  shoppingList, err := getShoppingList(client)
  if err != nil {
    fmt.Printf("Error getting shopping list: %s\n", err)
    var failedItem ShoppingListItem
    return failedItem, err
  }
  for _, item := range shoppingList {
    if item.IDHex == IDHex {
      return item, nil
    }
  }
  return ShoppingListItem{}, nil
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

func addShoppingListItem(client *mongo.Client, itemName string) (ShoppingListItem, error) {
  newIdFromTimeStamp := primitive.NewObjectIDFromTimestamp(time.Now())
  collection := client.Database("GoShopping").Collection("shopping-lists") 
  filter := bson.D{{}}
  update := bson.D{{Key: "$push", Value: bson.D{{Key: "ShoppingList", Value: bson.D{{Key: "Id", Value: newIdFromTimeStamp}, {Key: "Item", Value: itemName}, {Key: "Ticked", Value: false}}}}}}
  _, err := collection.UpdateOne(context.Background(), filter, update)
  if err != nil {
    fmt.Printf("Error adding shopping list item: %s\n", err)
    var failedItem ShoppingListItem
    return failedItem, err
  }
  newItem := ShoppingListItem {
    ID: newIdFromTimeStamp,
    IDHex: newIdFromTimeStamp.Hex(),
    Item: itemName,
    Ticked: false,
  }
  return newItem, nil
}

func updateShoppingListItem(client *mongo.Client, itemId string, newItem string) (ShoppingListItem, error) {
	collection := client.Database("GoShopping").Collection("shopping-lists")
  filter := bson.D{{}}
  update := bson.D{{Key: "$set", Value: bson.D{{Key: "ShoppingList.$[element].Item", Value: newItem}}}}
  options := options.UpdateOptions{
    ArrayFilters: &options.ArrayFilters{
      Filters: []interface{}{bson.D{{Key: "element.IDHex", Value: itemId}}},
    },
  }
  _, err := collection.UpdateOne(context.Background(), filter, update, &options)
  var shoppingListItem ShoppingListItem
  shoppingList, err := getShoppingList(client)
  if err != nil {
    fmt.Printf("Error updating shopping list item: %s\n", err)
    return shoppingListItem, err
  }
  for _, item := range shoppingList {
    if item.IDHex == itemId {
      shoppingListItem = item
    }
  }
  return shoppingListItem, nil
}

func deleteShoppingListItem(client *mongo.Client, item string) error {
	collection := client.Database("GoShopping").Collection("shopping-lists")
  filter := bson.D{{}}
  update := bson.D{{Key: "$pull", Value: bson.D{{Key: "ShoppingList", Value: bson.D{{Key: "IDHex", Value: item}}}}}}
  _, err := collection.UpdateOne(context.Background(), filter, update)
  if err != nil {
    fmt.Printf("Error deleting shopping list item: %s\n", err)
    return err
  }
  return nil
}

func sortShoppingList(client *mongo.Client, ids []string) ([]ShoppingListItem, error) {
  var newShoppingList []ShoppingListItem
  currentShoppingList, err := getShoppingList(client)
  if err != nil {
    fmt.Printf("Error getting shopping list: %s\n", err)
    return newShoppingList, err
  }
  for _, id := range ids {
    for _, item := range currentShoppingList {
      if item.IDHex == id {
        newShoppingList = append(newShoppingList, item)
      }
    }
  }
 
  collection := client.Database("GoShopping").Collection("shopping-lists")
  filter := bson.D{{}}
  update := bson.D{{Key: "$set", Value: bson.D{{Key: "ShoppingList", Value: newShoppingList}}}}
  _, err = collection.UpdateOne(context.Background(), filter, update)
  if err != nil {
    fmt.Printf("Error sorting shopping list: %s\n", err)
    return newShoppingList, err
  }
  fmt.Printf("Sorted shopping list: %v\n", newShoppingList)
  return newShoppingList, nil
}

func tickShoppingListItem(client *mongo.Client, itemId string, Ticked bool) error {
  itemIdFromHex, err := primitive.ObjectIDFromHex(itemId)
  if err != nil {
    fmt.Printf("Error converting item to object ID: %s\n", err)
    return err
  }
  collection := client.Database("GoShopping").Collection("shopping-lists")
  filter := bson.D{{}}
  update := bson.D{{Key: "$set", Value: bson.D{{Key: "ShoppingList.$[element].Ticked", Value: Ticked}}}}
  options := options.UpdateOptions{
    ArrayFilters: &options.ArrayFilters{
      Filters: []interface{}{bson.D{{Key: "element.Id", Value: itemIdFromHex}}},
    },
  }
  _, err = collection.UpdateOne(context.Background(), filter, update, &options)
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
