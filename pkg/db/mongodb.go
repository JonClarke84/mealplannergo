package db

import (
	"context"
	"fmt"

	"github.com/JonClarke84/mealplannergo/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB represents a MongoDB client connection
type MongoDB struct {
	Client *mongo.Client
}

// Ensure MongoDB implements DBInterface
var _ DBInterface = (*MongoDB)(nil)

// NewMongoDB creates a new MongoDB connection
func NewMongoDB(uri string) (*MongoDB, error) {
	// connect to MongoDB
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return nil, err
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return &MongoDB{Client: client}, nil
}

// GetShoppingList retrieves the shopping list from the database
func (m *MongoDB) GetShoppingList() ([]models.ShoppingListItem, error) {
	collection := m.Client.Database("GoShopping").Collection("shopping-lists")

	var document struct {
		ID           primitive.ObjectID        `bson:"_id"`
		ShoppingList []models.ShoppingListItem `bson:"ShoppingList"`
		SortOrder    []primitive.ObjectID      `bson:"SortOrder"`
	}

	err := collection.FindOne(context.Background(), bson.D{}).Decode(&document)
	if err != nil {
		return nil, err
	}

	// Create a map of item ID to ShoppingListItem for easy lookup
	itemMap := make(map[primitive.ObjectID]models.ShoppingListItem)
	for _, item := range document.ShoppingList {
		itemMap[item.ID] = item
	}

	// Sort the shopping list based on the SortOrder
	var sortedShoppingList []models.ShoppingListItem
	for _, id := range document.SortOrder {
		if item, exists := itemMap[id]; exists {
			sortedShoppingList = append(sortedShoppingList, item)
		}
	}

	return sortedShoppingList, nil
}

// GetShoppingListItemFromIDHex retrieves a shopping list item by its hex ID
func (m *MongoDB) GetShoppingListItemFromIDHex(IDHex string) (models.ShoppingListItem, error) {
	shoppingList, err := m.GetShoppingList()
	if err != nil {
		var failedItem models.ShoppingListItem
		return failedItem, err
	}
	for _, item := range shoppingList {
		if item.IDHex == IDHex {
			return item, nil
		}
	}
	return models.ShoppingListItem{}, nil
}

// UpdateMeal updates a meal for a specific day
func (m *MongoDB) UpdateMeal(day string, meal string) error {
	collection := m.Client.Database("GoShopping").Collection("meal-plans")
	filter := bson.D{{Key: "meals", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{Key: "day", Value: day}}}}}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "meals.$.meal", Value: meal}}}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Printf("Error updating meal: %s\n", err)
	}
	return nil
}

// AddShoppingListItem adds a new item to the shopping list
func (m *MongoDB) AddShoppingListItem(itemName string) (models.ShoppingListItem, error) {
	// if itemName is empty, return an error
	if itemName == "" {
		return models.ShoppingListItem{}, fmt.Errorf("item name cannot be empty")
	}

	// Generate a new ObjectID
	newId := primitive.NewObjectID()

	// Prepare the shopping list item to be added
	newItem := models.ShoppingListItem{
		ID:     newId,
		IDHex:  newId.Hex(),
		Item:   itemName,
		Ticked: false,
	}

	// Prepare the update operation to push the new item
	filter := bson.D{} // This filter needs to be specific to the document you're updating
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "ShoppingList", Value: newItem}}}}

	// Execute the update operation
	_, err := m.Client.Database("GoShopping").Collection("shopping-lists").UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Printf("Error adding shopping list item: %s\n", err)
		return models.ShoppingListItem{}, err
	}

	// Add the new item to the Order
	err = m.AddShoppingListIdToShoppingListOrder(newItem.IDHex)

	return newItem, nil
}

// AddShoppingListIdToShoppingListOrder adds a new item ID to the sort order
func (m *MongoDB) AddShoppingListIdToShoppingListOrder(itemId string) error {
	// Prepare the update operation to push the new item
	filter := bson.D{} // This filter needs to be specific to the document you're updating
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "SortOrder", Value: itemId}}}}

	// Execute the update operation
	_, err := m.Client.Database("GoShopping").Collection("shopping-lists").UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Printf("Error adding shopping list item to order: %s\n", err)
		return err
	}

	return nil
}

// UpdateShoppingListItem updates an existing shopping list item
func (m *MongoDB) UpdateShoppingListItem(itemId string, newItem string) (models.ShoppingListItem, error) {
	collection := m.Client.Database("GoShopping").Collection("shopping-lists")
	filter := bson.D{{}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "ShoppingList.$[element].Item", Value: newItem}}}}
	options := options.UpdateOptions{
		ArrayFilters: &options.ArrayFilters{
			Filters: []interface{}{bson.D{{Key: "element.IDHex", Value: itemId}}},
		},
	}
	_, err := collection.UpdateOne(context.Background(), filter, update, &options)
	var shoppingListItem models.ShoppingListItem
	shoppingList, err := m.GetShoppingList()
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

// DeleteShoppingListItem removes an item from the shopping list
func (m *MongoDB) DeleteShoppingListItem(itemIDHex string) error {
	collection := m.Client.Database("GoShopping").Collection("shopping-lists")

	filter := bson.M{}
	update := bson.M{"$pull": bson.M{"ShoppingList": bson.M{"IDHex": itemIDHex}}}
	if _, err := collection.UpdateOne(context.Background(), filter, update); err != nil {
		fmt.Printf("Error deleting shopping list item: %s\n", err)
		return err
	}

	var result struct {
		ShoppingList []models.ShoppingListItem `bson:"ShoppingList"`
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

// TickShoppingListItem toggles the ticked status of a shopping list item
func (m *MongoDB) TickShoppingListItem(itemId string, ticked bool) (models.ShoppingListItem, error) {
	collection := m.Client.Database("GoShopping").Collection("shopping-lists")
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
		var failedShoppingListItem models.ShoppingListItem
		return failedShoppingListItem, err
	}

	shoppingListItem, err := m.GetShoppingListItemFromIDHex(itemId)
	if err != nil {
		fmt.Printf("Error getting shopping list item: %s\n", err)
		return shoppingListItem, err
	}

	return shoppingListItem, nil
}

// GetMealPlan retrieves the meal plan from the database
func (m *MongoDB) GetMealPlan() (models.MealPlan, error) {
	// find the first document in the collection
	collection := m.Client.Database("GoShopping").Collection("meal-plans")
	// get the first document
	filter := bson.D{{}}
	var mealPlan models.MealPlan
	err := collection.FindOne(context.Background(), filter).Decode(&mealPlan)
	if err != nil {
		fmt.Printf("Error finding meal plan: %s\n", err)
		return mealPlan, err
	}
	return mealPlan, nil
}

// SortShoppingList updates the order of items in the shopping list
func (m *MongoDB) SortShoppingList(newOrder []models.Order) error {
	collection := m.Client.Database("GoShopping").Collection("shopping-lists")

	// Create a new array of ObjectIDs in the new order
	var newSortOrder []primitive.ObjectID
	for _, order := range newOrder {
		id, err := primitive.ObjectIDFromHex(order.ID)
		if err != nil {
			return fmt.Errorf("invalid object ID: %s", order.ID)
		}
		newSortOrder = append(newSortOrder, id)
	}

	// Update the SortOrder in the shopping list document
	filter := bson.D{} // This filter needs to be specific to the document you're updating
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "SortOrder", Value: newSortOrder}}}}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Printf("Error updating shopping list order: %s\n", err)
		return err
	}

	return nil
}

// Close closes the MongoDB connection
func (m *MongoDB) Close() {
	m.Client.Disconnect(context.TODO())
}