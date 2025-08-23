package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/viniciuscg/vinux/internal/input"
	"github.com/viniciuscg/vinux/internal/notify"
)

var selectMessagePattern = func(postFix string) string {
	return postFix + " (space to select, enter to confirm):"
}

var yesNoOptions = []string{
	"Yes",
	"No",
}

var commitPrefixOptions = []string{
	"build:",
	"chore:",
	"ci:",
	"docs:",
	"feat:",
	"fix:",
	"perf:",
	"refactor:",
	"revert:",
	"style:",
	"test:",
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commits changes to the repository",
	Run: func(cmd *cobra.Command, args []string) {
		notify.Print(
			notify.TypeInfo,
			"Starting commit process...",
			nil,
		)

		files := getFilesToCommit()

		selected := selectFiles(files)
		if len(selected) == 0 {
			notify.Print(
				notify.TypeWarning,
				"No files selected.",
				nil,
			)

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
		notify.Print(
			notify.TypeError,
			"Running git status.",
			err,
		)

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

	if len(files) == 0 {
		notify.Print(
			notify.TypeWarning,
			"No changes to commit.",
			nil,
		)

		return nil
	}

	return files
}

func selectFiles(files []string) []string {
	var selected []string
	prompt := &survey.MultiSelect{
		Message: selectMessagePattern("Select files to add on stash (based on Conventional Commits and commitlint config)"),
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
		notify.Print(
			notify.TypeError,
			"Adding files to git.",
			err,
		)

		return
	}

	notify.Print(
		notify.TypeSuccess,
		fmt.Sprintf(
			"✅ Added files: %v\n",
			files,
		),
		nil,
	)
}

func commit() {
	prefix := selectedPrefix()

	notify.Print(
		notify.TypeWrite,
		"Enter your commit message below",
		nil,
	)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	commitMessage := fmt.Sprintf("%s %s", prefix, input)
	cmd := exec.Command("git", "commit", "-m", commitMessage)
	err := cmd.Run()
	if err != nil {
		notify.Print(
			notify.TypeError,
			"Committing changes.",
			err,
		)

		return
	}

	notify.Print(
		notify.TypeSuccess,
		"✅ Changes committed successfully!",
		nil,
	)

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
	if !hasToPush() {
		notify.Print(
			notify.TypeBlock,
			"Skipping push.",
			nil,
		)

		return
	}

	cmd := exec.Command("git", "push")
	err := cmd.Run()
	if err != nil {
		notify.Print(
			notify.TypeError,
			"Pushing changes.",
			err,
		)

		return
	}

	notify.Print(
		notify.TypeSuccess,
		"Changes pushed successfully!",
		nil,
	)

	if hasOtherCommits() {
		commitAgain()
	}

}

func hasToPush() bool {
	return yesOrNoCheck(
		input.NewConsoleReader(os.Stdin),
		"Do you want to push the changes?",
	)
}

func hasOtherCommits() bool {
	return yesOrNoCheck(
		input.NewConsoleReader(os.Stdin),
		"Are there other commits to do?",
	)
}

func yesOrNoCheck(ir input.InputReader, message string) bool {
	input := ir.ReadInput(message + " (y or n): ")

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

func commitAgain() {
	getFilesToCommit()
	cmd := exec.Command("vinux", "commit")
	err := cmd.Run()
	if err != nil {
		notify.Print(
			notify.TypeError,
			"Running vinux commit.",
			err,
		)

		return
	}
}
