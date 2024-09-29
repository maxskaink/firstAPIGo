package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("Hello, World!")

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file: ", err)
	}

	MONGO_URI := os.Getenv("MONGO_URI")

	clientOptions := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection = client.Database("golan_db").Collection("todos")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:8081",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodos)
	app.Patch("/api/todos/:id", updateTodos)
	app.Delete("/api/todos/:id", deleteTodos)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	log.Fatal(app.Listen(":" + PORT))

}

func getTodos(c *fiber.Ctx) error {
	fmt.Println("Get Todos")
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}
	return c.JSON(todos)
}
func createTodos(c *fiber.Ctx) error {
	fmt.Println("Create Todos")
	todo := new(Todo)

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Body == "" {
		return c.Status(400).SendString("Body is required")
	}
	insertResult, err := collection.InsertOne(context.Background(), todo)

	if err != nil {
		return err
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(todo)
}
func updateTodos(c *fiber.Ctx) error {
	fmt.Println("Update Todos")
	id := c.Params("id")
	ObjectID, err := primitive.ObjectIDFromHex(id)

	cursor := collection.FindOne(context.Background(), bson.M{"_id": ObjectID})

	if err := cursor.Err(); err != nil {
		return c.Status(400).SendString("Todo not found")
	}
	todo := new(Todo)
	if err := cursor.Decode(&todo); err != nil {
		return err
	}

	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	filter := bson.M{"_id": ObjectID}
	update := bson.M{"$set": bson.M{"completed": !todo.Completed}}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"message": "Todo updated"})
}
func deleteTodos(c *fiber.Ctx) error {
	fmt.Println("Delete Todos")
	id := c.Params("id")
	ObjectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}
	filter := bson.M{"_id": ObjectID}
	_, err = collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"message": "Todo deleted"})
}
