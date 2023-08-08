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

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	commit.User = os.Getenv("USERNAME")
	commit.Email = os.Getenv("EMAIL")
	cmd.Execute()
}
