package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/viniciuscg/vinux/internal/input"
	"github.com/viniciuscg/vinux/internal/notify"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Sets up the linux based on VINI standards",
	Run: func(cmd *cobra.Command, args []string) {
		notify.Print(
			notify.TypeInfo,
			"Starting setup process...",
			nil,
		)

		isSudo := isSudo(input.NewConsoleReader(os.Stdin))

		if isSudo {
			fmt.Println("🔐 Using sudo for setup...")
			tmux("sudo")
		} else {
			fmt.Println("🔓 Not using sudo for setup...")
			tmux("")
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func tmux(sudo string) {
	notify.Print(
		notify.TypeInfo,
		"Installing tmux...",
		nil,
	)

	cmd := exec.Command(sudo, "apt", "install", "tmux")

	err := cmd.Run()
	if err != nil {
		fmt.Println("❌ Error installing tmux:", err)

		return
	}
	fmt.Println("✅ Tmux setup complete!")
}

func isSudo(ir input.InputReader) bool {
	input := ir.ReadInput("📝 Do you want to use sudo? (y or n): ")

	cleanInput := strings.ToLower(strings.TrimSpace(input))

	if strings.HasPrefix(cleanInput, "y") {
		return true
	}

	if strings.HasPrefix(cleanInput, "n") {
		return false
	}

	notify.Print(
		notify.TypeError,
		notify.ErrorInvalidYesOrNoInput,
		nil,
	)
	return false
}
