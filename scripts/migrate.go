package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Get MongoDB connection URI from environment variable
	mongoURI := os.Getenv("GO_SHOPPING_MONGO_ATLAS_URI")
	if mongoURI == "" {
		log.Fatal("Error: GO_SHOPPING_MONGO_ATLAS_URI environment variable is not set")
	}

	// Connect to MongoDB
	client, err := connectToMongoDB(mongoURI)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	fmt.Println("Starting data migration from GoShopping to GoShopping-test...")

	// Define collections to copy
	collections := []string{"shopping-lists", "meal-plans"}

	for _, collectionName := range collections {
		fmt.Printf("\nCopying collection: %s\n", collectionName)
		
		if err := copyCollection(client, collectionName); err != nil {
			log.Fatalf("Error copying collection %s: %v", collectionName, err)
		}
		
		fmt.Printf("✓ Successfully copied collection: %s\n", collectionName)
	}

	fmt.Println("\n🎉 Migration completed successfully!")
	fmt.Println("Your test database (GoShopping-test) now contains a copy of your production data.")
}

func connectToMongoDB(uri string) (*mongo.Client, error) {
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

	fmt.Println("✓ Connected to MongoDB Atlas")
	return client, nil
}

func copyCollection(client *mongo.Client, collectionName string) error {
	// Source and destination databases
	sourceDB := client.Database("GoShopping")
	destDB := client.Database("GoShopping-test")

	sourceCollection := sourceDB.Collection(collectionName)
	destCollection := destDB.Collection(collectionName)

	// Clear the destination collection first
	fmt.Printf("  Clearing existing data in GoShopping-test.%s...\n", collectionName)
	_, err := destCollection.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		return fmt.Errorf("failed to clear destination collection: %v", err)
	}

	// Get all documents from source collection
	cursor, err := sourceCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		return fmt.Errorf("failed to find documents in source collection: %v", err)
	}
	defer cursor.Close(context.TODO())

	// Collect all documents
	var documents []interface{}
	for cursor.Next(context.TODO()) {
		var doc bson.D
		if err := cursor.Decode(&doc); err != nil {
			return fmt.Errorf("failed to decode document: %v", err)
		}
		documents = append(documents, doc)
	}

	if err := cursor.Err(); err != nil {
		return fmt.Errorf("cursor error: %v", err)
	}

	// Insert documents into destination collection if any exist
	if len(documents) > 0 {
		fmt.Printf("  Copying %d documents...\n", len(documents))
		
		_, err = destCollection.InsertMany(context.TODO(), documents)
		if err != nil {
			return fmt.Errorf("failed to insert documents into destination collection: %v", err)
		}
	} else {
		fmt.Printf("  No documents found in source collection\n")
	}

	return nil
}