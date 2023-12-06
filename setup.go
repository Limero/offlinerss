package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/limero/offlinerss/models"
)

func setup() (models.Config, error) {
	answers := struct {
		Server   string
		Username string
		Password string
		Clients  []string
	}{}

	survey.ErrorTemplate = `{{color .Icon.Format }}{{ .Icon.Text }} Error: {{ .Error.Error }}{{color "reset"}}
`

	qs1 := []*survey.Question{
		{
			Name: "server",
			Prompt: &survey.Select{
				Message: "Pick server:",
				Options: []string{
					"NewsBlur",
					"Miniflux",
				},
			},
		},
	}

	if err := survey.Ask(qs1, &answers); err != nil {
		return models.Config{}, err
	}

	qs2 := []*survey.Question{
		{
			Name:   "username",
			Prompt: &survey.Input{Message: fmt.Sprintf("Enter username for %s:", answers.Server)},
			Validate: func(val interface{}) error {
				if str, ok := val.(string); !ok || len(str) == 0 {
					return errors.New("invalid username")
				}
				return nil
			},
		},
		{
			Name:   "password",
			Prompt: &survey.Password{Message: fmt.Sprintf("Enter password for %s:", answers.Server)},
			Validate: func(val interface{}) error {
				if str, ok := val.(string); !ok || len(str) == 0 {
					return errors.New("invalid password")
				}
				return nil
			},
		},
	}

	if err := survey.Ask(qs2, &answers); err != nil {
		return models.Config{}, err
	}

	serverConfig := models.ServerConfig{
		Name:     models.ServerName(strings.ToLower(answers.Server)),
		Username: answers.Username,
		Password: answers.Password,
	}

	server := getServer(serverConfig)
	fmt.Printf("Attempting to login to %s as user %q\n", serverConfig.Name, serverConfig.Username)
	if err := server.Login(); err != nil {
		return models.Config{}, err
	}
	fmt.Printf("Successfully logged in!\n\n")

	qs3 := []*survey.Question{
		{
			Name: "clients",
			Prompt: &survey.MultiSelect{
				Message: "Select clients:",
				Options: []string{
					"FeedReader",
					"Newsboat",
					"QuiteRSS",
				},
				VimMode: true,
			},
			Validate: func(val interface{}) error {
				if opts, ok := val.([]survey.OptionAnswer); !ok || len(opts) == 0 {
					return errors.New("you need to pick at least one client")
				}
				return nil
			},
		},
	}

	if err := survey.Ask(qs3, &answers); err != nil {
		return models.Config{}, err
	}

	clientConfigs := make([]models.ClientConfig, len(answers.Clients))
	for i, c := range answers.Clients {
		clientConfigs[i] = models.ClientConfig{
			Name: models.ClientName(strings.ToLower(c)),
		}
	}

	return models.Config{
		Server:  serverConfig,
		Clients: clientConfigs,
	}, nil
}
