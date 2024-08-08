package palettes

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
)

var stanfordURL string
var PalettesCache map[string][]map[string]string

func init() {
	godotenv.Load()
	stanfordURL = os.Getenv("STANFORD_URL")
}

func GetAbout(c *fiber.Ctx) error {
	text, err := os.ReadFile("./public/palettes/about.md")
	if err != nil {
		return c.Status(500).SendString("Error reading file")
	}
	response := parseMarkdown(string(text))
	return c.SendString(response)
}

func GetPalettesHandler(c *fiber.Ctx) error {
	return c.JSON(PalettesCache)
}

func GetPalettes() error {
	url := fmt.Sprintf("%s/api/tfh-data/palettes/all", stanfordURL)
	_, body, err := fasthttp.Get(nil, url)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &PalettesCache)
	if err != nil {
		return err
	}
	fmt.Println("Palettes cached")
	return nil
}

func PaletteAutoCache() {
	GetPalettes()
	ticker := time.NewTicker(4 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		GetPalettes()
	}
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
