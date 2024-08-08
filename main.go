package main

import (
	"log"
	"os"
	"tfhdata/packages/palettes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var Logger = logger.New(logger.Config{
	Format: "[${status}] ${method} ${protocol}://${host}${path}\n",
})

var port string

func init() {
	port = os.Getenv("PORT")
}

func main() {
	go palettes.PaletteAutoCache()
	app := fiber.New()
	app.Use(Logger)

	app.Use(func(c *fiber.Ctx) error {
		if c.Hostname() == "tfhdata.com" {
			return c.Redirect("https://tfh.enmadarei.com"+c.OriginalURL(), 301)
		}
		return c.Next()
	})

	api := app.Group("/api")
	api.Get("/palettes", palettes.GetPalettesHandler)
	api.Get("/palettes/about", palettes.GetAbout)
	app.Static("/", "./public")

	// 404 handler
	app.Use(NotFound)

	log.Fatal(app.Listen(":" + port))
}

func NotFound(c *fiber.Ctx) error {
	return c.Status(404).SendString("That doesn't Exist, Deer")
}
