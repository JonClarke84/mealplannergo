package main

import (
	"context"
	"fmt"

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
  // if itemName is empty, return an error
  if itemName == "" {
    return ShoppingListItem{}, fmt.Errorf("item name cannot be empty")
  }

  // Get the current shopping list to determine the new order
  shoppingList, err := getShoppingList(client)
  if err != nil {
    fmt.Printf("Error getting shopping list: %s\n", err)
    return ShoppingListItem{}, err // Directly return the zero value for ShoppingListItem
  }

  // New order is the length of the shopping list plus one
  order := len(shoppingList) + 1

  // Generate a new ObjectID
  newId := primitive.NewObjectID()

  // Prepare the shopping list item to be added
  newItem := ShoppingListItem{
    ID:     newId,
    IDHex:  newId.Hex(),
    Item:   itemName,
    Ticked: false,
    Order:  order,
  }

  // Prepare the update operation to push the new item
  filter := bson.D{} // This filter needs to be specific to the document you're updating
  update := bson.D{{Key: "$push", Value: bson.D{{Key: "ShoppingList", Value: newItem}}}}

  // Execute the update operation
  _, err = client.Database("GoShopping").Collection("shopping-lists").UpdateOne(context.Background(), filter, update)
  if err != nil {
    fmt.Printf("Error adding shopping list item: %s\n", err)
    return ShoppingListItem{}, err
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

func deleteShoppingListItem(client *mongo.Client, itemIDHex string) error {
	collection := client.Database("GoShopping").Collection("shopping-lists")

	filter := bson.M{}
	update := bson.M{"$pull": bson.M{"ShoppingList": bson.M{"IDHex": itemIDHex}}}
	if _, err := collection.UpdateOne(context.Background(), filter, update); err != nil {
		fmt.Printf("Error deleting shopping list item: %s\n", err)
		return err
	}

	var result struct {
		ShoppingList []ShoppingListItem `bson:"ShoppingList"`
	}
	if err := collection.FindOne(context.Background(), bson.M{}).Decode(&result); err != nil {
		fmt.Printf("Error retrieving updated shopping list post-deletion: %s\n", err)
		return err
	}

	for i, item := range result.ShoppingList {
		// Update the order to the current index
		update := bson.M{"$set": bson.M{"ShoppingList.$.Order": i + 1}}
		filter := bson.M{"ShoppingList.IDHex": item.IDHex}
		if _, err := collection.UpdateOne(context.Background(), filter, update); err != nil {
			fmt.Printf("Error updating order for item ID %s: %s\n", item.IDHex, err)
			return err
		}
	}

	return nil
}

func tickShoppingListItem(client *mongo.Client, itemId string, ticked bool) (ShoppingListItem, error) {
  collection := client.Database("GoShopping").Collection("shopping-lists")
  filter := bson.D{{}}
  update := bson.D{{Key: "$set", Value: bson.D{{Key: "ShoppingList.$[element].Ticked", Value: ticked}}}}
  options := options.UpdateOptions{
    ArrayFilters: &options.ArrayFilters{
      Filters: []interface{}{bson.D{{Key: "element.IDHex", Value: itemId}}},
    },
  }
  _, err := collection.UpdateOne(context.Background(), filter, update, &options)
  if err != nil {
    fmt.Printf("Error ticking shopping list item: %s\n", err)
    var failedShoppingListItem ShoppingListItem
    return failedShoppingListItem, err
  }

  shoppingListItem, err := getShoppingListItemFromIDHex(client, itemId)
  if err != nil {
    fmt.Printf("Error getting shopping list item: %s\n", err)
    return shoppingListItem, err
  }

  return shoppingListItem, nil
}

func updateShoppingListOrder(client *mongo.Client, newOrder []int) error {
  shoppingList, err := getShoppingList(client)
  if err != nil {
    fmt.Printf("Error getting shopping list: %s\n", err)
    return err
  }

  var newShoppingList []ShoppingListItem
  for _, order := range newOrder {
    for _, item := range shoppingList {
      if item.Order == order {
        newShoppingList = append(newShoppingList, item)
        break
      }
    }
  }

	collection := client.Database("GoShopping").Collection("shopping-lists")
  filter := bson.D{{}}
  update := bson.D{{Key: "$set", Value: bson.D{{Key: "ShoppingList", Value: newShoppingList}}}}
  _, err = collection.UpdateOne(context.Background(), filter, update)
  if err != nil {
    fmt.Printf("Error updating shopping list order: %s\n", err)
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
