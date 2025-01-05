package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	Id        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	fmt.Println("Hello world")
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file...")
	}

	PORT := os.Getenv("PORT")

	todos := []Todo{}

	app.Get("/api/todos", func(c *fiber.Ctx) error {
		// return c.Status(200).JSON(fiber.Map{"msg": "This is my first web app..."})
		return c.Status(200).JSON(todos)
	})

	//add todo
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}

		// log.Println("Raw Body:", string(c.Body()))

		if err := c.BodyParser(todo); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}
		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
		}

		todo.Id = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(201).JSON(todo)
	})

	//update todo
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.Id) == id {
				todos[i].Completed = true
				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	//delete todo
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.Id) == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(fiber.Map{"sucess": true})
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Todos not found"})
	})

	log.Fatal(app.Listen(":" + PORT))
}
