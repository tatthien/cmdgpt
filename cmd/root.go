package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	"github.com/tatthien/cmdgpt/internal/prompt"
	"github.com/urfave/cli/v2"
)

var version = "dev"
var green = color.New(color.FgGreen)
var boldGreen = green.Add(color.Bold)
var cyan = color.New(color.FgCyan)
var boldCyan = cyan.Add(color.Bold)

func init() {
	cli.VersionFlag = &cli.BoolFlag{Name: "version", Aliases: []string{"v"}}
	cli.VersionPrinter = func(cCtx *cli.Context) {
		fmt.Fprintf(cCtx.App.Writer, "version=%s\n", cCtx.App.Version)
	}
}

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

		query := prompt.StringPrompt("❓ What's a command would you like to ask?")

		if query == "" {
			fmt.Println("There is nothing to ask!")
			return nil
		}

		messages := []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a system that parses natural language to linux commands. You may not use natural language in your responses. You can respond with this format: <command>%sep%<explanation>. Only send one command per message.",
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

		response := resp.Choices[0].Message.Content
		substrings := strings.Split(response, "%sep%")
		cmd := ""
		explanation := "There are no explanations at this moment"

		if len(substrings) == 1 {
			cmd = substrings[0]
		} else {
			cmd = strings.TrimSpace(substrings[0])
			explanation = strings.TrimSpace(substrings[1])
		}

		fmt.Printf("💡 %s\n", boldCyan.Sprint("Command"))
		fmt.Println()
		fmt.Printf("   %s\n\n", cmd)
		fmt.Printf("💬 %s\n\n", boldCyan.Sprint("Explanation"))
		fmt.Printf("   %s\n\n", explanation)

		// Select options
		answer := ""
		prompt := &survey.Select{
			Message: "Select an action",
			Options: []string{"run", "copy", "cancel"},
			VimMode: true,
		}
		survey.AskOne(prompt, &answer)

		if answer == "run" {
			strUuid := uuid.New().String()
			strPath := fmt.Sprintf("/tmp/%s.sh", strUuid)
			ioutil.WriteFile(strPath, []byte(cmd), 0744)
			out, err := exec.Command("/bin/bash", strPath).CombinedOutput()
			if err != nil {
				return nil
			}
			os.Remove(strPath)
			fmt.Println(string(out))
		}

		if answer == "copy" {
			var copyCmd *exec.Cmd
			if runtime.GOOS == "darwin" {
				copyCmd = exec.Command("pbcopy")
			} else if runtime.GOOS == "linux" {
				copyCmd = exec.Command("xclip")
			} else {
				return fmt.Errorf("copy is not supported in %s", runtime.GOOS)
			}

			in, err := copyCmd.StdinPipe()
			if err != nil {
				return err
			}
			if _, err := in.Write([]byte(cmd)); err != nil {
				return err
			}
			if err := copyCmd.Start(); err != nil {
				return err
			}
			fmt.Println("Copied to clipboard")
		}

		return nil
	},
}

func Exec() error {
	if err := app.Run(os.Args); err != nil {
		return err
	}
	return nil
}
