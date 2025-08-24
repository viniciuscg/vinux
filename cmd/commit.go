package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/viniciuscg/survey/v2"
	"github.com/viniciuscg/vinux/internal/input"
	"github.com/viniciuscg/vinux/internal/notify"
)

var selectMessagePattern = func(postFix string) string {
	return postFix + " (space to select, enter to confirm):"
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
		startCommitFlow()
	},
}

func startCommitFlow() {
	notify.Print(
		notify.TypeInfo,
		"Starting commit process...",
		nil,
	)

	files := getFilesToCommit()
	if len(files) == 0 {
		notify.Print(
			notify.TypeWarning,
			"No changes to commit.",
			nil,
		)

		return
	}

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
}

func init() {
	rootCmd.AddCommand(commitCmd)
}

func getFilesToCommit() []string {
	out, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		notify.Print(
			notify.TypeError,
			"Running git status.",
			err,
		)

		return nil
	}

	lines := strings.Split(string(out), "\n")
	files := []string{"add all files"}
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
		Message: selectMessagePattern("Select files to add on stash (based on Conventional Commits and commitlint config)"),
		Options: files,
	}
	survey.AskOne(prompt, &selected)

	return selected
}

func addFiles(files []string) {
	if files[0] == "add all files" {
		files = []string{"."}
	}

	args := append([]string{"add"}, files...)

	if err := exec.Command("git", args...).Run(); err != nil {
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
			"Added files: %v\n",
			files,
		),
		nil,
	)
}

func commit() {
	prefix := selectedPrefix()

	message, err := input.ReadInput(
		"Enter your commit message:",
		survey.Required,
	)
	if err != nil {
		notify.Print(
			notify.TypeError,
			"Reading commit message.",
			err,
		)

		return
	}

	commitMessage := fmt.Sprintf("%s %s", prefix, message)
	if err := exec.Command("git", "commit", "-m", commitMessage).Run(); err != nil {
		notify.Print(
			notify.TypeError,
			"Committing changes.",
			err,
		)

		return
	}

	notify.Print(
		notify.TypeSuccess,
		"Changes committed successfully!",
		nil,
	)

}

func selectedPrefix() string {
	var selected string
	prompt := &survey.Select{
		Message: selectMessagePattern("Select commit Prefix"),
		Options: commitPrefixOptions,
	}
	survey.AskOne(prompt, &selected, survey.WithPageSize(20))

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

	if err := exec.Command("git", "push").Run(); err != nil {
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
		recommit()
	}
}

func hasToPush() bool {
	return input.YesOrNoCheck(
		"Do you want to push the changes?",
	)
}

func hasOtherCommits() bool {
	return input.YesOrNoCheck(
		"Are there other commits to do?",
	)
}

func recommit() {
	files := getFilesToCommit()
	if len(files) == 0 {
		notify.Print(
			notify.TypeWarning,
			"No changes to commit.",
			nil,
		)

		return
	}

	startCommitFlow()
}
