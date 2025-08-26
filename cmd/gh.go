package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/viniciuscg/survey/v2"
	"github.com/viniciuscg/vinux/internal/input"
	"github.com/viniciuscg/vinux/internal/notify"
)

type PR struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	URL    string `json:"url"`
}

var gitHubCmd = &cobra.Command{
	Use:   "gh",
	Short: "GitHub related commands",
}

var gitHubSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		setupGitHub()
	},
}

var gitHubPrListCmd = &cobra.Command{
	Use:   "pr",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		listOpenPrs()
	},
}

var gitHubPrCreateCmd = &cobra.Command{
	Use:   "pr-create",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		listOpenPrs()
	},
}

func init() {
	rootCmd.AddCommand(gitHubCmd)

	gitHubCmd.AddCommand(gitHubSetupCmd)
	gitHubCmd.AddCommand(gitHubPrListCmd)
	gitHubCmd.AddCommand(gitHubPrCreateCmd)
}

func setupGitHub() {
	notify.Print(
		notify.TypeInfo,
		"Setting up GitHub...",
		nil,
	)

	installGitHub()
	authenticateGitHub()
	testingAuthenticateGitHub()

	notify.Print(
		notify.TypeSuccess,
		"GitHub setup completed!",
		nil,
	)
}

func listOpenPrs() {
	cmd := exec.Command("gh", "pr", "list", "--json", "number,title,url")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		notify.Print(
			notify.TypeError,
			"Opening PR list.",
			err,
		)
		return
	}

	var prs []PR
	if err := json.Unmarshal(out, &prs); err != nil {
		notify.Print(
			notify.TypeError,
			"Parsing PR list JSON.",
			err,
		)
		return
	}

	if len(prs) == 0 {
		notify.Print(
			notify.TypeInfo,
			"Nenhum PR aberto",
			nil,
		)
		return
	}

	for _, pr := range prs {
		id := color.New(color.FgHiGreen, color.Bold).Sprintf("%d:", pr.Number)
		title := color.New(color.FgHiYellow).Sprintf("%s", pr.Title)

		fmt.Printf(
			"%s %s\n",
			id,
			hyperlink(pr.URL, title),
		)
	}
}

func createPr() {
	title, err := input.ReadInput(
		"Enter the PR title:",
		survey.Required,
		nil,
	)
	if err != nil {
		notify.Print(
			notify.TypeError,
			"Reading PR title.",
			err,
		)

		return
	}

	var body string

	if input.YesOrNoCheck("Do you want to add a body to the PR?") {
		inputBody, err := input.ReadInput(
			"Enter the PR body:",
			nil,
			nil,
		)
		if err != nil {
			notify.Print(
				notify.TypeError,
				"Reading PR body.",
				err,
			)

			return
		}

		body = inputBody
	}

	cmd := exec.Command(
		"gh", "pr", "create",
		"--title", title,
		"--body", body,
		"--draft",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func installGitHub() {
	_, err := exec.Command("sudo", "apt", "install", "gh", "-y").Output()
	if err != nil {
		notify.Print(
			notify.TypeError,
			"Running git status.",
			err,
		)

		return
	}
	notify.Print(
		notify.TypeSuccess,
		"GitHub CLI installed successfully!",
		nil,
	)
}

func authenticateGitHub() {
	cmd := exec.Command("gh", "auth", "login")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		notify.Print(
			notify.TypeError,
			"Authenticating GitHub CLI.",
			err,
		)

		return
	}

	notify.Print(
		notify.TypeSuccess,
		"GitHub authentication completed!",
		nil,
	)
}

func testingAuthenticateGitHub() {
	out, err := exec.Command("gh", "auth", "status").Output()
	if err != nil {
		notify.Print(
			notify.TypeError,
			"Testing GitHub authentication.",
			err,
		)

		return
	}

	notify.Print(
		notify.TypeSuccess,
		"GitHub authentication status:\n"+string(out),
		nil,
	)
}

func hyperlink(url, text string) string {
	return fmt.Sprintf("\x1b]8;;%s\x07%s\x1b]8;;\x07", url, text)
}
