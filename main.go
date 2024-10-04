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
	openaiModel = openai.GPT4oMini20240718
	// openaiModel = openai.GPT4oMini
	// openaiModel = openai.GPT4o
	// openaiModel = openai.O1Preview
	geminiModel = "gemini-1.5-pro"
	outputFile  = "/tmp/chat-ai.log"
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

	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return
	}
	defer file.Close()

	switch aiProvider {
	case "chatgpt":
		runChatGPT(apiKey, reader, aiProvider, file)
	case "gemini":
		runGemini(apiKey, reader, aiProvider)
	default:
		log.Fatalf("Invalid AI provider, select 'chatgpt' or 'gemini'.")
	}
}

// runChatGPT handles the interaction loop with ChatGPT, processing user commands and displaying responses.
func runChatGPT(apiKey string, reader *bufio.Reader, aiProvider string, file *os.File) {
	openaiClient := openai.NewClient(apiKey)

	var conversation []openai.ChatCompletionMessage

	for {
		command, shouldContinue := prompt(reader, aiProvider)

		if !shouldContinue {
			break
		}

		if command == "h" && shouldContinue {
			displayHelp()
			continue
		} else {
			writeChatMessage(file, "> "+command+"\n")
		}

		// Add the user's message to the conversation history.
		conversation = append(conversation, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: command,
		})

		resp, err := chatGPTChat(openaiClient, conversation)
		// resp, err := chatGPTChat(openaiClient, command)
		if err != nil {
			log.Fatalf("ChatGPT failed: %v", err)
		}
		fmt.Println(resp.Choices[0].Message.Content)

		writeChatMessage(file, resp.Choices[0].Message.Content+"\n"+"\n")

		conversation = append(conversation, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: resp.Choices[0].Message.Content,
		})
	}
}

// chatGPTChat sends the user's command to the ChatGPT AI and returns the response.
func chatGPTChat(openaiClient *openai.Client, conversation []openai.ChatCompletionMessage) (resp openai.ChatCompletionResponse, err error) {
	fmt.Println(color.Format(color.GREEN, "> Waiting for ChatGPT.."))

	resp, err = openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openaiModel,
			Messages: conversation,
		},
	)

	if err != nil {
		return
	}

	return resp, nil
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
		log.Fatalf("Gemini genai NewClient failed: %v", err)
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
			log.Fatalf("Gemini failed: %v", err)
		}

		fmt.Println(respGemini.Candidates[0].Content.Parts[0])
	}
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

// writeChatMessage writes ChatGPT answer message to a specified file, it's like keeping a history.
func writeChatMessage(file *os.File, message string) {
	_, err := file.WriteString(message)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}
