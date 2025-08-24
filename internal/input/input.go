package input

import (
	"strings"

	"github.com/viniciuscg/survey/v2"
)

func ReadInput(prompt string, validators ...survey.Validator) (string, error) {
	var input string
	err := survey.Ask(
		[]*survey.Question{
			{
				Name:      "name",
				Prompt:    &survey.Input{Message: prompt},
				Validate:  survey.,
				Transform: survey.Title,
			},
		},

		&input,
	)

	validation := survey.ComposeValidators(validators...)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(input), nil
}
