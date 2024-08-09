package framedata

import (
	"encoding/json"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
)

var StanfordURL string
var frameDataCache atomic.Value

type MoveFD struct {
	ID         string `json:"id"`
	Attack     string `json:"attack"`
	Input      string `json:"input"`
	AtkDisplay string `json:"atk_display"`
	Startup    string `json:"startup"`
	Active     string `json:"active"`
	Recovery   string `json:"recovery"`
	Advantage  string `json:"advantage"`
}

func init() {
	godotenv.Load()
	StanfordURL = os.Getenv("STANFORD_URL")
	frameDataCache.Store(make(map[string]map[string]map[string]MoveFD))
}

func GetFrameDataHandler(c *fiber.Ctx) error {
	frameData := frameDataCache.Load().(map[string]map[string]map[string]MoveFD)
	if len(frameData) == 0 {
		fmt.Println("Frame Data cache is empty, attempting to update cache...")
		err := getFrameDataCache()
		if err != nil {
			fmt.Println("Error getting frame data:", err)
			return c.Status(500).SendString("Error getting frame data")
		}
		frameData = frameDataCache.Load().(map[string]map[string]map[string]MoveFD)
	}
	return c.JSON(frameData)
}

func getFrameDataCache() error {
	newCache := make(map[string]map[string]map[string]MoveFD)
	url := fmt.Sprintf("%s/api/tfh-data/framedata", StanfordURL)
	_, body, err := fasthttp.Get(nil, url)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &newCache)
	if err != nil {
		return err
	}

	frameDataCache.Store(newCache)
	fmt.Println("Frame Data cached")
	return nil
}

func UpdateFrameDataHandler(c *fiber.Ctx) error {
	err := getFrameDataCache()
	if err != nil {
		return c.Status(500).SendString("Error getting frame data")
	}
	return c.SendString("Frame Data cache updated")
}

func AutoUpdateFrameDataCache() {
	err := getFrameDataCache()
	if err != nil {
		fmt.Println("Error getting frame data:", err)
	}
	ticker := time.NewTicker(8 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		err := getFrameDataCache()
		if err != nil {
			fmt.Println("Error getting frame data:", err)
		}
	}
}
