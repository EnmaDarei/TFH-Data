package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
)

var Logger = logger.New(logger.Config{
	Format: "[${status}] ${method} ${protocol}://${host}${path}\n",
})

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	stanfordURL := os.Getenv("STANFORD_URL")

	app := fiber.New()
	app.Use(Logger)
	app.Use(func(c *fiber.Ctx) error {
		if c.Hostname() == "tfhdata.com" {
			return c.Redirect("https://tfh.enmadarei.com"+c.OriginalURL(), 301)
		}
		return c.Next()
	})

	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.Redirect("https://images.candyfloof.com/tfh-data/arieyes.ico", 301)
	})

	api := app.Group("/api")

	palettes := api.Group("/palettes")

	palettes.Get("/about-text", func(c *fiber.Ctx) error {
		text, err := os.ReadFile("./public/palettes/about.txt")
		if err != nil {
			return c.Status(500).SendString("Error reading file")
		}
		response := parseMarkdown(string(text))
		return c.SendString(response)
	})
	palettes.Get("/:character", func(c *fiber.Ctx) error {
		character := c.Params("character")
		paletteInfo, err := GetPalettes(character, stanfordURL)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(paletteInfo)
	})

	app.Static("/", "./public")

	// 404 handler
	app.Use(NotFound)

	log.Fatal(app.Listen(":" + port))
}

func GetPalettes(character, stanfordURL string) (map[int]map[string]string, error) {
	url := fmt.Sprintf("%s/api/tfh-data/palettes/%s", stanfordURL, character)
	_, body, err := fasthttp.Get(nil, url)
	if err != nil {
		return nil, err
	}

	var paletteInfo []string
	err = json.Unmarshal(body, &paletteInfo)
	if err != nil {
		return nil, err
	}

	palettes := make(map[int]map[string]string)
	for i, palette := range paletteInfo {
		if palette == "" {
			continue
		}
		imageName := strings.ReplaceAll(strings.ReplaceAll(palette, "/", "-"), ":", "-")
		image := fmt.Sprintf("https://images.candyfloof.com/tfh-data/palettes/%s/%s.png", character, imageName)
		palettes[i] = map[string]string{"name": palette, "image": image}
	}

	return palettes, nil
}

func NotFound(c *fiber.Ctx) error {
	return c.Status(404).SendString("That doesn't Exist, Deer")
}

func parseMarkdown(markdown string) string {
	// Handle headers
	markdown = regexp.MustCompile(`(?m)^# (.+)$`).ReplaceAllString(markdown, "<h1>$1</h1>")
	markdown = regexp.MustCompile(`(?m)^## (.+)$`).ReplaceAllString(markdown, "<h2>$1</h2>")
	markdown = regexp.MustCompile(`(?m)^### (.+)$`).ReplaceAllString(markdown, "<h3>$1</h3>")
	markdown = regexp.MustCompile(`(?m)^#### (.+)$`).ReplaceAllString(markdown, "<h4>$1</h4>")

	// Handle bold
	markdown = regexp.MustCompile(`\*\*(.*?)\*\*`).ReplaceAllString(markdown, "<strong>$1</strong>")

	// Handle italic
	markdown = regexp.MustCompile(`\*(.*?)\*`).ReplaceAllString(markdown, "<em>$1</em>")

	// Handle links
	markdown = regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`).ReplaceAllString(markdown, `<a href="$2" target="_blank">$1</a>`)

	// Handle underlined
	markdown = regexp.MustCompile(`__(.*?)__`).ReplaceAllString(markdown, "<u>$1</u>")

	// Handle unordered lists
	markdown = regexp.MustCompile(`(?m)^\s*\*\s(.+)$`).ReplaceAllString(markdown, "<li>$1</li>")
	markdown = strings.ReplaceAll(markdown, "</li>\n<li>", "</li><li>")
	markdown = regexp.MustCompile(`(<li>.*</li>)`).ReplaceAllString(markdown, "<ul>$1</ul>")

	// Handle paragraphs
	markdown = regexp.MustCompile(`(?m)^([^<].+)$`).ReplaceAllString(markdown, "<p>$1</p>")

	return markdown
}
