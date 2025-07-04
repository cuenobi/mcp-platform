/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"

	"github.com/cuenobi/mcp-platform/mcp-server-jira/cmd"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Could not load .env file; using system environment")
	}

	cmd.Execute()
}
