package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/vertexai/genai"
	"github.com/michalswi/color"
	openai "github.com/sashabaranov/go-openai"
	"google.golang.org/api/option"
)

const appName = "chat-ai"

func main() {

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("Please set the API_KEY env variable")
	}

	reader := bufio.NewReader(os.Stdin)

	var aiProvider string
	flag.StringVar(&aiProvider, "p", "", "AI provider [chatgpt, gemini]")
	flag.Parse()

	switch aiProvider {
	case "chatgpt":
		runChatGPT(apiKey, reader, aiProvider)
	case "gemini":
		runGemini(apiKey, reader, aiProvider)
	default:
		log.Fatalf("Invalid AI provider: %s", aiProvider)
	}
}

func runChatGPT(apiKey string, reader *bufio.Reader, aiProvider string) {
	openaiClient := openai.NewClient(apiKey)

	for {
		command, shouldContinue := prompt(reader, aiProvider)
		if !shouldContinue {
			break
		}

		resp, err := chatGPTChat(openaiClient, command)
		if err != nil {
			log.Fatalf("ChatGPT failed: %v\n", err)
		}
		fmt.Println(resp.Choices[0].Message.Content)
	}
}

func runGemini(apiKey string, reader *bufio.Reader, aiProvider string) {
	ctx := context.Background()

	projectID := os.Getenv("VAI_PROJECT_ID")
	if projectID == "" {
		log.Fatal("Please set the VAI_PROJECT_ID env variable")
	}

	region := os.Getenv("VAI_REGION")
	if region == "" {
		log.Fatal("Please set the VAI_REGION env variable")
	}

	client, err := genai.NewClient(ctx, projectID, region, option.WithCredentialsFile(apiKey))
	if err != nil {
		log.Fatalf("Gemini genai NewClient failed: %v\n", err)
	}

	for {
		command, shouldContinue := prompt(reader, aiProvider)
		if !shouldContinue {
			break
		}

		respGemini, err := geminiChat(ctx, client, command)
		if err != nil {
			log.Fatalf("Gemini failed: %v\n", err)
		}

		fmt.Println(respGemini.Candidates[0].Content.Parts[0])
	}
}

func prompt(reader *bufio.Reader, aiProvider string) (string, bool) {
	prompt := fmt.Sprintf("%s [%s:%s]: ", time.Now().UTC().Format(time.RFC1123), appName, aiProvider)
	fmt.Printf(color.Format(color.YELLOW, prompt))

	command, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading command:", err)
		return "", true
	}
	command = strings.TrimSpace(command)

	switch command {
	case "q":
		fmt.Println(color.Format(color.YELLOW, "Exiting chat. bye!"))
		return "", false
	}

	return command, true
}

func geminiChat(ctx context.Context, client *genai.Client, command string) (resp *genai.GenerateContentResponse, err error) {
	fmt.Println(color.Format(color.GREEN, "> Waiting for Gemini.."))

	// https://pkg.go.dev/cloud.google.com/go/vertexai/genai#Client.GenerativeModel
	model := client.GenerativeModel("gemini-1.5-pro")

	const ChatTemperature float32 = 0.1
	temperature := ChatTemperature
	model.Temperature = &temperature

	chatSession := model.StartChat()

	var builder strings.Builder

	fmt.Fprintln(&builder, command)
	introductionString := builder.String()

	resp, err = chatSession.SendMessage(ctx, genai.Text(introductionString))
	if err != nil {
		return
	}

	return resp, nil
}

func chatGPTChat(openaiClient *openai.Client, command string) (resp openai.ChatCompletionResponse, err error) {
	fmt.Println(color.Format(color.GREEN, "> Waiting for ChatGPT.."))

	resp, err = openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: command,
				},
			},
			// max number of tokens to generate in response
			// MaxTokens: 2000,
		},
	)

	if err != nil {
		return
	}

	return resp, nil
}
