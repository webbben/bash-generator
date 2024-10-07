package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	llama "github.com/webbben/ollama-wrapper"
)

const sysPrompt = `You are an assistant that generates bash commands. 
Given a prompt, you should generate a bash command that fulfills the prompt and works for %s shell running on %s OS.
Give a short description of what the command does. If there are arguments, flags, chained commands, etc, list and briefly explain them.

Do not say anything else besides this, unless you think it is very important context for the user to know.

Examples:

1) "Create a new directory called 'test'."

mkdir test

Description: Create the directory 'test'.

2) "List all files in the current directory, including the size of each file."

ls -lh

Description: List information about the files (the current directory by default).

l: use a long listing format
h: with -l, print sizes in human readable format (e.g., 1K 234M 2G)
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a sentence as an argument.")
		fmt.Println("Usage: bash-gen \"<your-sentence>\"")
		return
	}

	userPrompt := os.Args[1]

	client, err := llama.GetClient()
	if err != nil {
		log.Fatal(err)
	}

	// Start spinner
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				for _, r := range `▁▂▃▄▅▆▇█▇▆▅▄▃▁` {
					fmt.Printf("\r%c ", r)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}()

	// get OS and shell

	shell := "bash"
	if env, ok := os.LookupEnv("SHELL"); !ok {
		fmt.Println("SHELL environment variable not set. Defaulting to regular bash.")
	} else {
		parts := strings.Split(env, "/")
		shell = parts[len(parts)-1]
	}
	fmt.Println("Shell:", shell)
	fmt.Println("OS:", runtime.GOOS)
	sys := fmt.Sprintf(sysPrompt, runtime.GOOS, shell)

	res, err := llama.GenerateCompletion(client, sys, userPrompt)
	done <- true
	time.Sleep(100 * time.Millisecond)
	fmt.Print("\r") // Clear spinner

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\033[2J\033[1;1H") // Clear screen and move cursor to top-left corner
	fmt.Println(res)
}
