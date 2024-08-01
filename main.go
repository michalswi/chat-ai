package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/vertexai/genai"
	"github.com/michalswi/color"
	"google.golang.org/api/option"
)

func main() {

	apiKey := os.Getenv("API_KEY")
	projectID := os.Getenv("VAI_PROJECT_ID")
	region := os.Getenv("VAI_REGION")

	reader := bufio.NewReader(os.Stdin)

	for {
		ctx := context.Background()
		appName := "chat-ai"

		fmt.Printf(time.Now().UTC().Format(time.RFC1123)+" [%s]: ", appName)
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		if command == "q" {
			break
		}

		client, err := genai.NewClient(ctx, projectID, region, option.WithCredentialsFile(apiKey))
		if err != nil {
			log.Fatalf("genai new client failed: %v\n", err)
		}

		// https://pkg.go.dev/cloud.google.com/go/vertexai/genai#Client.GenerativeModel
		model := client.GenerativeModel("gemini-1.5-pro")

		const ChatTemperature float32 = 0.1
		temperature := ChatTemperature
		model.Temperature = &temperature

		chatSession := model.StartChat()

		var builder strings.Builder

		fmt.Fprintln(&builder, command)
		introductionString := builder.String()
		var response *genai.GenerateContentResponse

		fmt.Println(color.Format(color.GREEN, "Gemini review started.."))

		response, err = chatSession.SendMessage(ctx, genai.Text(introductionString))
		if err != nil {
			log.Fatalf("Gemini review failed: %v\n", err)
		}

		fmt.Println(response.Candidates[0].Content.Parts[0])
	}
}
