/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/ForbiddenR/auto-commit/cmd"
	"github.com/joho/godotenv"
)

// The username, email and author are set in the .env file.
// The commit message is set with the -m flag.
// The commit message is required.
// The commit message is split by commas if you have multiple changes.
func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	cmd.Execute()
}