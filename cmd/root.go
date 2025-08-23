package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vinux",
	Short: "Vinux CLI is your crazy Go-powered CLI!",
	Run: func(cmd *cobra.Command, args []string) {
		printBanner()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printBanner() {
	fig := figure.NewFigure("VINUX CLI!", "starwars", true)
	lines := strings.Split(fig.String(), "\n")

	primary := color.New(color.FgHiCyan, color.Bold)

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			fmt.Println()

			continue
		}
		primary.Println(line)
	}

	primary.Println("\nUse 'vinux --help' to see available commands.")
	primary.Println("\n✨ May the Vini be with you! 🚀")
}
