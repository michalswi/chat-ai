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

const (
	appName = "chat-ai"
	// https://pkg.go.dev/github.com/sashabaranov/go-openai#pkg-constants
	openaiModel = openai.GPT4oMini
	// openaiModel = openai.GPT4o
	// openaiModel = openai.O1Preview
	geminiModel = "gemini-1.5-pro"
)

func main() {

	var aiProvider string
	flag.StringVar(&aiProvider, "p", "", "AI provider [chatgpt, gemini]")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("Please set the API_KEY env variable")
	}

	switch aiProvider {
	case "chatgpt":
		runChatGPT(apiKey, reader, aiProvider)
	case "gemini":
		runGemini(apiKey, reader, aiProvider)
	default:
		log.Fatalf("Invalid AI provider, select 'chatgpt' or 'gemini'.")
	}
}

// runChatGPT handles the interaction loop with ChatGPT, processing user commands and displaying responses.
func runChatGPT(apiKey string, reader *bufio.Reader, aiProvider string) {
	openaiClient := openai.NewClient(apiKey)

	for {
		command, shouldContinue := prompt(reader, aiProvider)

		if !shouldContinue {
			break
		}

		if command == "h" && shouldContinue {
			displayHelp()
			continue
		}

		resp, err := chatGPTChat(openaiClient, command)
		if err != nil {
			log.Fatalf("ChatGPT failed: %v\n", err)
		}
		fmt.Println(resp.Choices[0].Message.Content)
	}
}

// runGemini handles the interaction loop with Gemini, processing user commands and displaying responses.
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

		if command == "h" && shouldContinue {
			displayHelp()
			continue
		}

		respGemini, err := geminiChat(ctx, client, command)
		if err != nil {
			log.Fatalf("Gemini failed: %v\n", err)
		}

		fmt.Println(respGemini.Candidates[0].Content.Parts[0])
	}
}

// prompt displays a prompt to the user, reads their input, and determines whether to continue the loop or handle special commands like "h" for help.
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
		fmt.Println(color.Format(color.GREEN, "Exiting chat. bye!"))
		return "", false
	case "h":
		return "h", true
	}

	return command, true
}

// displayHelp prints out the available commands and their descriptions to the user.
func displayHelp() {
	fmt.Println(color.Format(color.GREEN, "Commands:"))
	fmt.Println(color.Format(color.GREEN, "  q - quit: Exit the chat."))
	fmt.Println(color.Format(color.GREEN, "  h - help: Display this help message."))
}

// geminiChat sends the user's command to the Gemini AI and returns the response.
func geminiChat(ctx context.Context, client *genai.Client, command string) (resp *genai.GenerateContentResponse, err error) {
	fmt.Println(color.Format(color.GREEN, "> Waiting for Gemini.."))

	// https://pkg.go.dev/cloud.google.com/go/vertexai/genai#Client.GenerativeModel
	model := client.GenerativeModel(geminiModel)

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

// chatGPTChat sends the user's command to the ChatGPT AI and returns the response.
func chatGPTChat(openaiClient *openai.Client, command string) (resp openai.ChatCompletionResponse, err error) {
	fmt.Println(color.Format(color.GREEN, "> Waiting for ChatGPT.."))

	resp, err = openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openaiModel,
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
