/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"darius/cmd"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment() // hoặc zap.NewProduction()
	zap.ReplaceGlobals(logger)
	cmd.Execute()
}
