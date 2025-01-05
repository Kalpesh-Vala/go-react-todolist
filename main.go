package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("Hello world")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file : ", err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal("Error connecting to MongoDB : ", err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal("Error pinging to MongoDB : ", err)
	}

	fmt.Println("Connected to MongoDB")

	collection = client.Database("todo-golang-db").Collection("todos")

	app := fiber.New()

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error fetching todos"})
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

func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)

	// Parse the body into the Todo struct
	if err := c.BodyParser(todo); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate the body field
	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
	}

	// Generate a new ObjectID for the todo
	todo.ID = primitive.NewObjectID()

	// Insert the todo into the collection
	_, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error creating todo"})
	}

	// Return the created todo
	return c.Status(201).JSON(todo)
}

func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"completed": true}}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error updating todo"})
	}

	return c.Status(200).JSON(fiber.Map{"sucess": true})
}

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	filter := bson.M{"_id": objectId}
	_, err = collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error deleting todo"})
	}

	return c.Status(200).JSON(fiber.Map{"sucess": true})

}

// package main

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/joho/godotenv"
// )

// type Todo struct {
// 	Id        int    `json:"id"`
// 	Completed bool   `json:"completed"`
// 	Body      string `json:"body"`
// }

// func main() {
// 	fmt.Println("Hello world")
// 	app := fiber.New()

// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal("Error loading .env file...")
// 	}

// 	PORT := os.Getenv("PORT")

// 	todos := []Todo{}

// 	app.Get("/", func(c *fiber.Ctx) error {
// 		return c.Status(200).JSON(fiber.Map{"msg": "This is my first web app using go and react..."})
// 	})

// 	app.Get("/api/todos", func(c *fiber.Ctx) error {
// 		// return c.Status(200).JSON(fiber.Map{"msg": "This is my first web app..."})
// 		return c.Status(200).JSON(todos)
// 	})

// 	//add todo
// 	app.Post("/api/todos", func(c *fiber.Ctx) error {
// 		todo := &Todo{}

// 		// log.Println("Raw Body:", string(c.Body()))

// 		if err := c.BodyParser(todo); err != nil {
// 			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
// 		}
// 		if todo.Body == "" {
// 			return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
// 		}

// 		todo.Id = len(todos) + 1
// 		todos = append(todos, *todo)

// 		return c.Status(201).JSON(todo)
// 	})

// 	//update todo
// 	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
// 		id := c.Params("id")

// 		for i, todo := range todos {
// 			if fmt.Sprint(todo.Id) == id {
// 				todos[i].Completed = true
// 				return c.Status(200).JSON(todos[i])
// 			}
// 		}

// 		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
// 	})

// 	//delete todo
// 	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
// 		id := c.Params("id")

// 		for i, todo := range todos {
// 			if fmt.Sprint(todo.Id) == id {
// 				todos = append(todos[:i], todos[i+1:]...)
// 				return c.Status(200).JSON(fiber.Map{"sucess": true})
// 			}
// 		}

// 		return c.Status(404).JSON(fiber.Map{"error": "Todos not found"})
// 	})

// 	log.Fatal(app.Listen(":" + PORT))
// }
