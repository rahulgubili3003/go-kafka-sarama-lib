package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

type Comment struct {
	Text string `form:"text" json:"text"`
}

func main() {
	fmt.Println("Hello")
	app := fiber.New()
	api := app.Group("/api/v1") // /api

	api.Post("/comments", createComment)

	err := app.Listen(":3000")
	if err != nil {
		return
	}
}

func createComment(c *fiber.Ctx) error {

	// Instantiate new Message struct
	cmt := new(Comment)

	//  Parse body into comment struct
	if err := c.BodyParser(cmt); err != nil {
		log.Println(err)
		err := c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
		if err != nil {
			return err
		}
		return err
	}
	return nil
}
