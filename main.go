package main

import (
	"encoding/base64"
	"fmt"
	"github.com/HuascarBaptista/bitrise-step-update-version-upload-comment/jira"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
	"os"
	"strings"
)

// Config ...
type Config struct {
	UserName        string `env:"user_name,required"`
	APIToken        string `env:"api_token,required"`
	BaseURL         string `env:"base_url,required"`
	IssueKeys       string `env:"jira_tickets"`
	Message         string `env:"comments,required"`
	Version         string `env:"version"`
	AdditionalLabel string `env:"label"`
}

func main() {
	var cfg Config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Issue with input: %s", err)
	}
	if len(cfg.IssueKeys) == 0 {
		log.Infof("NO TICKETS:")
		os.Exit(0)
	}

	stepconf.Print(cfg)
	fmt.Println()

	encodedToken := generateBase64APIToken(cfg.UserName, cfg.APIToken)
	client := jira.NewClient(encodedToken, cfg.BaseURL)
	issueKeys := strings.Split(cfg.IssueKeys, `|`)

	var comments []jira.Comment
	for _, issueKey := range issueKeys {
		comments = append(comments, jira.Comment{Content: cfg.Message, IssuKey: issueKey})
	}

	if err := client.PostIssueCommentsAndVersionAndLabel(comments, cfg.Version, cfg.AdditionalLabel); err != nil {
		failf("Posting comments failed with error: %s", err)
	}
	os.Exit(0)
}

func generateBase64APIToken(userName string, apiToken string) string {
	v := userName + `:` + apiToken
	return base64.StdEncoding.EncodeToString([]byte(v))
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}
