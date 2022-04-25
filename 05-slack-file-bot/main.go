package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	channels := []string{os.Getenv("SLACK_CHANNEL_ID")}
	files := []string{"example.txt"}

	for i := 0; i < len(files); i++ {
		params := slack.FileUploadParameters{
			File:     files[i],
			Channels: channels,
		}
		file, err := api.UploadFile(params)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		fmt.Printf("Name: %s, URL: %s\n", file.Name, file.URL)
	}
}
