package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/viniciuscg/survey/v2"
	"github.com/viniciuscg/vinux/internal/notify"
	"github.com/viniciuscg/vinux/internal/storage"
)

var repoAddCmd = &cobra.Command{
	Use:   "add [path]",
	Short: "Add your repository path",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addReposPath(args[0])
	},
}

var repoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List saved repositories",
	Run: func(cmd *cobra.Command, args []string) {
		if !storage.VerifyPathExists(storage.GetReposDir()) {
			notify.Print(
				notify.TypeError,
				"Storage path does not exist. Please initialize first.",
				nil,
			)

			return
		}

		listRepos()
	},
}

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "If type 'vinux repo repoName' you will enter in the repo shell if not it will show a list of repos to select",
	Run: func(cmd *cobra.Command, args []string) {
		if !storage.VerifyPathExists(storage.GetReposDir()) {
			notify.Print(
				notify.TypeError,
				"Storage path does not exist. Please initialize first.",
				nil,
			)

			return
		}

		if len(args) > 0 {
			goToRepo(args[0])
		}

		goToRepoWithoutArgs()
	},
}

func init() {
	rootCmd.AddCommand(repoCmd)

	repoCmd.AddCommand(repoAddCmd)
	repoCmd.AddCommand(repoListCmd)
}

func addReposPath(path string) {
	if !storage.VerifyPathExists(path) {
		notify.Print(
			notify.TypeError,
			"Path does not exist, please try again.",
			nil,
		)

		return
	}

	if err := storage.SaveReposPath(path); err != nil {
		notify.Print(
			notify.TypeError,
			"Saving repo path",
			err,
		)

		return
	}

	notify.Print(
		notify.TypeSuccess,
		fmt.Sprintf("Repos path saved %s", path),
		nil,
	)
}

func listRepos() {
	repos, err := storage.ListRepos()
	if err != nil {
		notify.Print(
			notify.TypeError,
			"Listing repos",
			err,
		)

		return
	}

	if len(repos) == 0 {
		notify.Print(
			notify.TypeInfo,
			"No repos saved",
			nil,
		)

		return
	}

	for _, path := range repos {
		notify.Print(
			notify.TypeFolder,
			path,
			nil,
		)
	}
}

func goToRepo(name string) {
	path, err := storage.GetRepo(name)
	if err != nil {
		if os.IsNotExist(err) {
			notify.Print(
				notify.TypeError,
				"Repo not found",
				err,
			)

			return
		}

		notify.Print(
			notify.TypeError,
			"Getting repo path",
			err,
		)

		return
	}

	cmd := exec.Command("zsh")
	cmd.Dir = path
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

}

func goToRepoWithoutArgs() {
	repos, err := storage.ListRepos()
	if err != nil {
		notify.Print(
			notify.TypeError,
			"Getting current directory",
			err,
		)

		return
	}

	goToRepo(selectedRepo(repos))
}

func selectedRepo(repos []string) string {
	var selected string
	prompt := &survey.Select{
		Message: "Select a repository:",
		Options: repos,
	}
	survey.AskOne(prompt, &selected, survey.WithPageSize(20))

	return selected
}
