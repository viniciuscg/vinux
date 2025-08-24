package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/viniciuscg/survey/v2"
	"github.com/viniciuscg/vinux/internal/input"
	"github.com/viniciuscg/vinux/internal/notify"
)

type SetupType int

const (
	TypeTmux SetupType = iota

	TypeGogh

	TypeVsCode

	TypeGithub
)

type SetupOption struct {
	Label  string
	Type   SetupType
	Action func()
}

var setupOptionsList = []SetupOption{
	{Label: "tmux", Type: TypeTmux, Action: tmux},
	{Label: "gogh", Type: TypeGogh, Action: gogh},
	{Label: "vsCode", Type: TypeVsCode, Action: vsCode},
	{Label: "github", Type: TypeGithub, Action: github},
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Sets up the linux based on VINI standards",
	Run: func(cmd *cobra.Command, args []string) {
		startSetupFlow()
	},
}

func startSetupFlow() {
	notify.Print(
		notify.TypeInfo,
		"Starting setup process...",
		nil,
	)

	options := selectSetupOptions()
	if len(options) == 0 {
		notify.Print(
			notify.TypeWarning,
			"No setup options selected. Exiting.",
			nil,
		)

		return
	}

	for _, opt := range setupOptionsList {
		for _, selected := range options {
			if opt.Type == selected {
				opt.Action()
			}
		}
	}

	notify.Print(
		notify.TypeSuccess,
		"Setup process completed!",
		nil,
	)
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func selectSetupOptions() []SetupType {
	labelToTypeMap := make(map[string]SetupType)
	var labels []string

	for _, opt := range setupOptionsList {
		labels = append(labels, opt.Label)
		labelToTypeMap[opt.Label] = opt.Type
	}

	var selectedLabels []string
	prompt := &survey.MultiSelect{
		Message: "Select what you want to install or configure",
		Options: labels,
	}
	survey.AskOne(prompt, &selectedLabels)

	var selectedTypes []SetupType
	for _, label := range selectedLabels {
		selectedTypes = append(selectedTypes, labelToTypeMap[label])
	}

	return selectedTypes
}

func tmux() {
	notify.Print(
		notify.TypeInfo,
		"Installing tmux...",
		nil,
	)

	if err := exec.Command("sudo", "apt", "install", "tmux").Run(); err != nil {
		fmt.Println("Error installing tmux:", err)

		return
	}

	fmt.Println("✅ Tmux setup complete!")
}

func gogh() {
	notify.Print(
		notify.TypeInfo,
		"Setting up gogh...",
		nil,
	)

	cmd := exec.Command("sudo", "wget", "-O", "/usr/local/bin/gogh", "https://git.io/vQgMr")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		notify.Print(notify.TypeError, "Downloading gogh script.", err)

		return
	}

	cmd = exec.Command("sudo", "chmod", "+x", "/usr/local/bin/gogh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		notify.Print(notify.TypeError, "Making gogh executable.", err)

		return
	}

	cmd = exec.Command("gogh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		notify.Print(notify.TypeError, "Error running gogh:", err)

		return
	}

	fmt.Println("Gogh setup complete!")
}

func vsCode() {
	notify.Print(
		notify.TypeInfo,
		"Downloading VS Code .deb package...",
		nil,
	)

	debPath := "/tmp/vscode.deb"

	cmd := exec.Command("wget", "-qO", debPath, "https://code.visualstudio.com/sha/download?build=stable&os=linux-deb-x64")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		notify.Print(notify.TypeError, "Downloading VS Code failed.", err)
		return
	}

	notify.Print(
		notify.TypeInfo,
		"Installing VS Code...",
		nil,
	)

	cmd = exec.Command("sudo", "apt", "install", "-y", debPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		notify.Print(notify.TypeError, "Installing VS Code failed.", err)
		return
	}

	notify.Print(
		notify.TypeSuccess,
		"VS Code installed successfully!",
		nil,
	)

	if err := os.Remove(debPath); err != nil {
		notify.Print(
			notify.TypeWarning,
			"Cleaning up VS Code .deb file failed.",
			err,
		)
	}
}

func github() {
	notify.Print(
		notify.TypeInfo,
		"Setting up GitHub...",
		nil,
	)

	email, err := input.ReadInput(
		"Enter your GitHub email:",
		survey.EmailValidator,
	)
	if err != nil {
		notify.Print(
			notify.TypeError,
			"Reading GitHub email.",
			err,
		)

		return
	}

	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-C", email)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		notify.Print(
			notify.TypeError,
			"Generating SSH key.",
			err,
		)

		return
	}

	cmd = exec.Command("eval", "$(ssh-agent -s)")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		notify.Print(
			notify.TypeError,
			"Starting ssh-agent.",
			err,
		)

		return
	}

	cmd = exec.Command("ssh-add", "~/.ssh/id_ed25519")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		notify.Print(
			notify.TypeError,
			"Adding SSH key to ssh-agent.",
			err,
		)

		return
	}

	data, err := os.ReadFile("~/.ssh/id_ed25519.pub")
	if err != nil {
		fmt.Println("Error reading SSH key:", err)
		return
	}

	notify.Print(
		notify.TypeSuccess,
		fmt.Sprintf(
			"SSH key generated! Add it to your GitHub account: %s\n",
			data,
		),
		nil,
	)

	for {
		cmd := exec.Command("ssh", "-T", "git@github.com")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			notify.Print(
				notify.TypeSuccess,
				"SSH connection successful!",
				nil,
			)

			break
		}

		retry := input.YesOrNoCheck("SSH connection failed. Do you want to retry?")
		if !retry {
			notify.Print(
				notify.TypeWarning,
				"SSH setup aborted by user.",
				nil,
			)

			break
		}

		notify.Print(
			notify.TypeRetry,
			"Retrying SSH connection...",
			nil,
		)
	}

	notify.Print(
		notify.TypeSuccess,
		"GitHub setup complete!",
		nil,
	)
}
