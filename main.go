package main

import (
	"log"
	"os"
	fd "tfhdata/packages/framedata"
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
	go fd.AutoUpdateFrameDataCache()
	app := fiber.New()
	app.Use(Logger)

	app.Use(func(c *fiber.Ctx) error {
		if c.Hostname() == "tfhdata.com" {
			return c.Redirect("https://tfh.enmadarei.com"+c.OriginalURL(), 301)
		}
		return c.Next()
	})

	//private files
	app.Get("/404.css", func(c *fiber.Ctx) error {
		return c.SendFile("./private/404/404.css")
	})

	api := app.Group("/api")
	api.Get("/palettes", palettes.GetPalettesHandler)
	api.Get("/palettes/about", palettes.GetAbout)
	api.Get("/palettes/update", checkAuth, palettes.UpdateCacheHandler)
	api.Get("/framedata", fd.GetFrameDataHandler)
	api.Get("/framedata/update", checkAuth, fd.UpdateFrameDataHandler)
	app.Static("/", "./public")

	// 404 handler
	app.Use(NotFound)

	log.Fatal(app.Listen(":" + port))
}

func NotFound(c *fiber.Ctx) error {
	return c.Status(404).SendFile("./private/404/404.html")
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
