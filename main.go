package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Todo struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

var todos = []*Todo{
	{Id: 1, Name: "Cleber", Completed: false},
	{Id: 2, Name: "Claudio", Completed: false},
}

func main() {
	app := fiber.New()

	// Default middleware config
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	SetupTodosRoutes(app)

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}

func SetupApiV1(app *fiber.App) {
	v1 := app.Group("/v1")

	SetupTodosRoutes(v1)
}

func SetupTodosRoutes(grp fiber.Router) {
	todosRoutes := grp.Group("/todos")
	todosRoutes.Get("/", GetTodos)
	todosRoutes.Post("/", CreateTodo)
	todosRoutes.Get("/:id", GetTodo)
	todosRoutes.Delete("/:id", DeleteTodo)
	todosRoutes.Patch("/:id", UpdateTodo)
}

func GetTodos(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(todos)
}

func CreateTodo(c *fiber.Ctx) error {
	type request struct {
		Name string `json:"name"`
	}

	var body request

	err := c.BodyParser(&body)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse json",
		})
		return c.Status(fiber.StatusOK).JSON("OK")
	}
	todo := &Todo{
		Id:        len(todos) + 1,
		Name:      body.Name,
		Completed: false,
	}

	todos = append(todos, todo)

	return c.Status(fiber.StatusCreated).JSON(todos)
}

func GetTodo(c *fiber.Ctx) error {
	paramsId := c.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
	}

	for _, todo := range todos {
		if todo.Id == id {
			return c.Status(fiber.StatusOK).JSON(todo)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON("Not found.")
}

func DeleteTodo(c *fiber.Ctx) error {
	paramsId := c.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
	}

	for i, todo := range todos {
		if todo.Id == id {
			todos = append(todos[0:i], todos[i+1:]...)
			return c.Status(fiber.StatusOK).JSON("Todo Deleted.")
		}
	}

	return c.Status(fiber.StatusNotFound).JSON("Not found.")
}

func UpdateTodo(c *fiber.Ctx) error {
	type request struct {
		Name      *string `json:"name"` // o * serve para deixar o campo automaticamente obrigatório
		Completed *bool   `json:"completed"`
	}

	paramsId := c.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
	}

	var body request
	err = c.BodyParser(&body)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse body",
		})
	}

	var todo *Todo

	for _, t := range todos {
		if t.Id == id {
			todo = t
			break
		}
	}

	if todo == nil {
		return c.Status(fiber.StatusNotFound).JSON("Todo not found...")
	}

	if body.Name != nil {
		todo.Name = *body.Name // o * é utilizado para receber o valor do ponteiro
	}

	if body.Completed != nil {
		todo.Completed = *body.Completed
	}

	// agora é certeza que o Todo foi atualizado
	return c.Status(fiber.StatusOK).JSON(todo)
}
