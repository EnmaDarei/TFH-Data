package palettes

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	database "tfhdata/packages/db"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Palette struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Palettes struct {
	Velvet     []Palette `json:"velvet"`
	Tianhuo    []Palette `json:"tianhuo"`
	Arizona    []Palette `json:"arizona"`
	Oleander   []Palette `json:"oleander"`
	Paprika    []Palette `json:"paprika"`
	Pom        []Palette `json:"pom"`
	Shanty     []Palette `json:"shanty"`
	Stronghoof []Palette `json:"stronghoof"`
	Texas      []Palette `json:"texas"`
}

var palettesCache Palettes
var cache_mux sync.Mutex

func GetAbout(c *fiber.Ctx) error {
	text, err := os.ReadFile("./public/palettes/about.md")
	if err != nil {
		return c.Status(500).SendString("Error reading file")
	}
	response := parseMarkdown(string(text))
	return c.SendString(response)
}

func GetPalettesHandler(c *fiber.Ctx) error {
	cache_mux.Lock()
	palettes := palettesCache
	cache_mux.Unlock()
	return c.JSON(palettes)
}

func UpdateCacheHandler(c *fiber.Ctx) error {
	err := GetPalettes()
	if err != nil {
		return c.Status(500).SendString("Error getting palettes")
	}
	return c.SendString("Palettes cache updated")
}

func GetPalettes() error {
	var new_cache Palettes
	characters := []string{"velvet", "arizona", "paprika", "tianhuo", "oleander", "pom", "shanty", "stronghoof", "texas"}
	queryString := fmt.Sprintf("SELECT %v FROM palettes ORDER BY slot asc", strings.Join(characters, ","))
	rows, err := database.DB.Query(queryString)
	if err != nil {
		log.Println("Error querying the database:", err)
	}
	defer rows.Close()
	for rows.Next() {
		paletteNames := make([]string, len(characters))
		scanArgs := make([]interface{}, len(characters))
		for i := range characters {
			scanArgs[i] = &paletteNames[i]
		}

		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("Error scanning row:", err)
			return err
		}

		for i, char := range characters {
			if paletteNames[i] != "-" {
				palette := Palette{
					Name:  paletteNames[i],
					Image: fmt.Sprintf("https://images.candyfloof.com/tfh-data/palettes/%s/%s.png", char, paletteNames[i]),
				}
				switch char {
				case "velvet":
					new_cache.Velvet = append(new_cache.Velvet, palette)
				case "arizona":
					new_cache.Arizona = append(new_cache.Arizona, palette)
				case "paprika":
					new_cache.Paprika = append(new_cache.Paprika, palette)
				case "tianhuo":
					new_cache.Tianhuo = append(new_cache.Tianhuo, palette)
				case "oleander":
					new_cache.Oleander = append(new_cache.Oleander, palette)
				case "pom":
					new_cache.Pom = append(new_cache.Pom, palette)
				case "shanty":
					new_cache.Shanty = append(new_cache.Shanty, palette)
				case "stronghoof":
					new_cache.Stronghoof = append(new_cache.Stronghoof, palette)
				case "texas":
					new_cache.Texas = append(new_cache.Texas, palette)
				}
			}
		}
	}
	cache_mux.Lock()
	defer cache_mux.Unlock()
	palettesCache = new_cache
	fmt.Println("Palettes cached")
	return nil
}

func PaletteAutoCache() {
	err := GetPalettes()
	if err != nil {
		fmt.Println("Error getting palettes:", err)
	}
	ticker := time.NewTicker(4 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		err := GetPalettes()
		if err != nil {
			fmt.Println("Error getting palettes:", err)
		}
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
