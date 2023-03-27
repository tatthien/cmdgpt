package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/sashabaranov/go-openai"
	"github.com/tatthien/cmdgpt/internal/prompt"
	"github.com/urfave/cli/v2"
)

var version string
var green = color.New(color.FgGreen)
var boldGreen = green.Add(color.Bold)

var app = &cli.App{
	Name:     "cmdgpt",
	Usage:    "cmdgpt",
	Version:  version,
	Compiled: time.Now(),
	Authors: []*cli.Author{
		&cli.Author{
			Name:  "Thien Nguyen",
			Email: "me@thien.dev",
		},
	},
	Action: func(ctx *cli.Context) error {
		openaiApiKey := os.Getenv("OPENAI_API_KEY")
		if openaiApiKey == "" {
			return fmt.Errorf("missing OPENAI_API_KEY")
		}

		query := prompt.StringPrompt(boldGreen.Sprint("?") + " What's a command would you like to ask?")

		if query == "" {
			fmt.Println("There is nothing to ask!")
			return nil
		}

		messages := []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an AI linux commands generator. Your job is to analyze the user's input and convert it to the linux commands. Only response the command. Do not reponse the natural language or explaination about the command. Make sure the command is a valid syntax and will not contains any error. Try to figure the best commands that fit the user's input.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("Here is the input: %s", query),
			},
		}

		openaiClient := openai.NewClient(openaiApiKey)
		resp, err := openaiClient.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:       openai.GPT3Dot5Turbo,
				Messages:    messages,
				MaxTokens:   300,
				Temperature: 0.7,
				TopP:        1,
			},
		)

		if err != nil {
			return fmt.Errorf("chat completion error: %w", err)
		}

		if len(resp.Choices) == 0 {
			return fmt.Errorf("no answer")
		}

		color.Cyan(resp.Choices[0].Message.Content)
		return nil
	},
}

func Exec() error {
	if err := app.Run(os.Args); err != nil {
		return err
	}
	return nil
}
