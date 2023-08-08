/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"

	"github.com/ForbiddenR/autocommit/cmd"
	"github.com/ForbiddenR/autocommit/cmd/commit"
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
	commit.User = os.Getenv("USERNAME")
	commit.Email = os.Getenv("EMAIL")
	commit.Author = os.Getenv("AUTHOR")
	cmd.Execute()
}