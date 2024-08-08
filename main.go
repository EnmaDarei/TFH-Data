package main

import (
	"log"
	"os"
	"tfhdata/packages/palettes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var port string
var authKey string

type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

var Logger = logger.New(logger.Config{
	Format: "[${status}] ${method} ${protocol}://${host}${path}\n",
})

func init() {
	port = os.Getenv("PORT")
	authKey = os.Getenv("AUTH_TOKEN")
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
	api.Get("/palettes/update", checkAuth, palettes.UpdateCacheHandler)
	app.Static("/", "./public")

	// 404 handler
	app.Use(NotFound)

	log.Fatal(app.Listen(":" + port))
}

func NotFound(c *fiber.Ctx) error {
	return c.Status(404).SendString("That doesn't Exist, Deer")
}

func checkAuth(c *fiber.Ctx) error {
	if c.Get("authorization") != authKey {
		// fmt.Println("Unauthorized request on", c.Path())
		response := ErrorResponse{
			Error:  "Unauthorized",
			Status: 401,
		}
		return c.Status(401).JSON(response)
	}
	return c.Next()
}
