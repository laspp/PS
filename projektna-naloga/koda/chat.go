// Minimal multi-turn chat with llama-go.
//
// Usage:
// 1. build the llama-go library
// 2. download a LLM model in GGUF format from huggingface
// 3. On Arnes cluster run: 		    module load Go
// 4. Add path to the llama-go library: export LIBRARY_PATH=../llama-go C_INCLUDE_PATH=../llama-go
// 5. Run: 		   				        srun --pty --constraint=turin go run . Qwen3-1.7B-Q8_0.gguf
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	llama "github.com/liminalpurple/llama-go"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run . /path/to/model.gguf")
	}

	// Print only errors from llama.cpp.
	os.Setenv("LLAMA_LOG", "error")
	llama.InitLogging()

	// Load the model weights, run on CPU only
	model, err := llama.LoadModel(os.Args[1], llama.WithGPULayers(0))
	if err != nil {
		log.Fatal(err)
	}
	defer model.Close()

	// Create a new context with a 4096 token context window and 4 threads.
	ctx, err := model.NewContext(llama.WithContext(4096), llama.WithThreads(4))
	if err != nil {
		log.Fatal(err)
	}
	defer ctx.Close()

	opts := llama.ChatOptions{
		MaxTokens:      llama.Int(256),    // Max Answer tokens
		EnableThinking: llama.Bool(false), // Disable thinking
	}

	// Create the conversation history and add a system prompt
	history := []llama.ChatMessage{
		{Role: "system", Content: "You are a helpful assistant. Answer briefly."},
	}

	fmt.Println("Welcome! Type your question and press Enter (Ctrl+D to quit).")
	in := bufio.NewScanner(os.Stdin)
	fmt.Print("\nyou> ")
	for in.Scan() {
		text := in.Text()
		if text != "" {
			// Send the whole history plus the new question, wait for the answer.
			question := llama.ChatMessage{Role: "user", Content: text}
			reply, err := ctx.Chat(context.Background(), append(history, question), opts)
			if err != nil {
				log.Printf("chat error: %v", err)
			} else {
				fmt.Println("bot>", reply.Content)
				// Store the question and the answer for the next turn.
				answer := llama.ChatMessage{Role: "assistant", Content: reply.Content}
				history = append(history, question, answer)
			}
		}
		fmt.Print("\nyou> ")
	}
}
