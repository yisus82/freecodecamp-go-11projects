package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/krognol/go-wolfram"
	"github.com/shomali11/slacker"
	"github.com/tidwall/gjson"

	witai "github.com/wit-ai/wit-go/v2"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Printf("%+v\n", event)
	}
}

func main() {
	godotenv.Load(".env")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	client := witai.NewClient(os.Getenv("WIT_AI_TOKEN"))
	wolframClient := wolfram.Client{AppID: os.Getenv("WOLFRAM_APP_ID")}
	go printCommandEvents(bot.CommandEvents())

	bot.Command("question - <message>", &slacker.CommandDefinition{
		Description: "Send any question to wolfram",
		Example:     "Who is the president of Spain?",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			message := request.Param("message")
			msg, err := client.Parse(&witai.MessageRequest{
				Query: message,
			})
			if err != nil {
				response.Reply("Error: " + err.Error())
				return
			}

			data, err := json.MarshalIndent(msg, "", "    ")
			if err != nil {
				response.Reply("Error: " + err.Error())
				return
			}

			rough := string(data[:])
			value := gjson.Get(rough, "entities.wit$wolfram_search_query:wolfram_search_query.0.value")
			fmt.Println(value)
			query := value.String()

			answer, err := wolframClient.GetSpokentAnswerQuery(query, wolfram.Metric, 1000)
			if err != nil {
				response.Reply("Error: " + err.Error())
				return
			}

			response.Reply(answer)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
