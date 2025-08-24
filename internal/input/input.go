package input

import (
	"strings"

	"github.com/viniciuscg/survey/v2"
	"github.com/viniciuscg/vinux/internal/notify"
)

func ReadInput(prompt string, validators ...survey.Validator) (string, error) {
	var input string

	err := survey.Ask(
		[]*survey.Question{
			{
				Name:      "name",
				Prompt:    &survey.Input{Message: prompt},
				Validate:  survey.ComposeValidators(validators...),
				Transform: survey.Title,
			},
		},
		&input,
	)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(input), nil
}

func YesOrNoCheck(message string) bool {
	_, err := ReadInput(
		message+" (y or n): ",
		survey.HasYesOrNoPrefix,
	)
	if err != nil {
		notify.Print(
			notify.TypeError,
			"Reading yes/no input.",
			err,
		)

		return false
	}

	s := strings.ToLower(message)
	return strings.Contains(s, "y")
}
