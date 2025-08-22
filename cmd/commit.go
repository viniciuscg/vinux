package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var selectMessagePattern = func(postFix string) string {
	return postFix + " (space to select, enter to confirm):"
}

var yesNoOptions = []string{
	"Yes",
	"No",
}

var commitPrefixOptions = []string{
	"feat: ",
	"fix: ",
	"docs: ",
	"style: ",
	"refactor: ",
	"perf: ",
	"test: ",
	"chore: ",
	"revert: ",
	"build: ",
	"ci: ",
	"config: ",
	"release: ",
	"other: ",
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commits changes to the repository",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🌊 Starting commit process...")

		files := getFilesToCommit()
		if len(files) == 0 {
			fmt.Println("⚠️ No changes to commit.")
			return
		}

		selected := selectFiles(files)
		if len(selected) == 0 {
			fmt.Println("⚠️ No files selected.")
			return
		}

		addFiles(selected)
		commit()
		push()
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}

func getFilesToCommit() []string {
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("❌ Error running git status:", err)
		return nil
	}

	lines := strings.Split(string(out), "\n")
	var files []string

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			files = append(files, parts[1])
		}
	}

	return files
}

func selectFiles(files []string) []string {
	var selected []string
	prompt := &survey.MultiSelect{
		Message: selectMessagePattern("Select files to add on stash"),
		Options: files,
	}
	survey.AskOne(prompt, &selected)
	return selected
}

func addFiles(files []string) {
	args := append([]string{"add"}, files...)
	cmd := exec.Command("git", args...)
	err := cmd.Run()
	if err != nil {
		fmt.Println("❌ Error adding files:", err)
		return
	}

	fmt.Printf("✅ Added files: %v\n", files)
}

func commit() {
	prefix := selectedPrefix()
	fmt.Print("📝 Enter your commit message: ")
	reader := bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)
	commitMessage := fmt.Sprintf("%s %s", prefix, userInput)
	cmd := exec.Command("git", "commit", "-m", commitMessage)
	err := cmd.Run()
	if err != nil {
		fmt.Println("❌ Error committing changes:", err)
		return
	}

	fmt.Println("✅ Changes committed successfully!")
}

func selectedPrefix() string {
	var selected string
	prompt := &survey.Select{
		Message: selectMessagePattern("Select commit Prefix"),
		Options: commitPrefixOptions,
	}
	survey.AskOne(prompt, &selected)
	return selected
}

func push() {
	var selected string
	prompt := &survey.Select{
		Message: "Do you want to push the changes?",
		Options: yesNoOptions,
	}
	survey.AskOne(prompt, &selected)

	if selected == "No" {
		fmt.Println("🚫 Skipping push.")
		return
	}

	cmd := exec.Command("git", "push")
	err := cmd.Run()
	if err != nil {
		fmt.Println("❌ Error pushing changes:", err)
		return
	}

	fmt.Println("🚀 Changes pushed successfully!")
}
