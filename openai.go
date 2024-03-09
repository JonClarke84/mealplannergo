// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"os"
// )

// func getShoppingListFromAI() (map[string]string, error) {
// 	type Message struct {
// 		Role    string `json:"role"`
// 		Content string `json:"content"`
// 	}

// 	type RequestBody struct {
// 		Model          string          `json:"model"`
// 		Messages       []Message       `json:"messages"`
// 		ResponseFormat json.RawMessage `json:"response_format"`
// 	}

// 	type ResponseBody struct {
// 		Choices []struct {
// 			Message struct {
// 				Content string `json:"content"`
// 			} `json:"message"`
// 		} `json:"choices"`
// 	}

// 	apiKey := os.Getenv("OPENAI_API_KEY") // Set your API key in your environment variables
// 	if apiKey == "" {
// 		return nil, fmt.Errorf("OPENAI_API_KEY is not set in environment variables")
// 	}

// 	data := RequestBody{
// 		Model: "gpt-4",
// 		Messages: []Message{
// 			{
// 				Role:    "system",
// 				Content: "You are a helpful assistant designed to output JSON.",
// 			},
// 			{
// 				Role:    "user",
// 				Content: "I need a meal plan for the week with each day's meal suggestions.",
// 			},
// 		},
// 		ResponseFormat: json.RawMessage(`{"type": "json_object"}`),
// 	}

// 	requestBytes, err := json.Marshal(data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBytes))
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	respBytes, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var responseBody ResponseBody
// 	err = json.Unmarshal(respBytes, &responseBody)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Assuming the AI response is a JSON object with days as keys and meal suggestions as values
// 	var shoppingList map[string]string
// 	err = json.Unmarshal([]byte(responseBody.Choices[0].Message.Content), &shoppingList)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println("Shopping list from AI:", shoppingList)
// 	return shoppingList, nil
// }
