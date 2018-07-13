// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/shihanng/mdstory/github"
	"github.com/shihanng/mdstory/medium"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A helper to help you setup your Medium and GitHub credentials",
	Long: `A helper that asks for your Medium's integration token and
GitHub's personal access token and stores them in the config file.
It will overwrite the existing config file.`,
	RunE: loginRun,
	Args: cobra.ExactArgs(0),
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func require(v string) error {
	if v == "" {
		return errors.New("required")
	}
	return nil
}

var defaultPromptTemplates = &promptui.PromptTemplates{
	Prompt:          "{{ . }}: ",
	Valid:           "{{ . }}: ",
	Invalid:         "{{ . }}: ",
	Success:         "{{ . }}: ",
	ValidationError: "({{ . }})",
}

func loginRun(_ *cobra.Command, _ []string) error {
	promptMedium := promptui.Prompt{
		Label:     "Medium's integration token",
		Validate:  require,
		Templates: defaultPromptTemplates,
	}

	fmt.Println("Visit https://medium.com/me/settings to generate a new integration tokens.")
	mediumToken, err := promptMedium.Run()
	if err != nil {
		return errors.Wrap(err, "in getting Medium's integration token")
	}

	promptGitHub := promptui.Prompt{
		Label:     "GitHub's personal access token",
		Validate:  require,
		Templates: defaultPromptTemplates,
	}

	fmt.Println("\nVisit https://github.com/settings/tokens/new to generate a new personal access tokens.")
	fmt.Println("Scope: gist, Create gists")
	githubToken, err := promptGitHub.Run()
	if err != nil {
		return errors.Wrap(err, "in getting GitHub's personal access token")
	}

	fmt.Println("")

	mediumClient, err := medium.New(mediumToken)
	if err != nil {
		return err
	}
	fmt.Printf("will use Medium as: %s\n", mediumClient.Username())
	viper.Set("medium_access_token", mediumToken)

	githubClient, err := github.New(githubToken)
	if err != nil {
		return err
	}
	fmt.Printf("will use GitHub as: %s\n", githubClient.Username())
	viper.Set("github_access_token", githubToken)

	touch()

	if err := viper.WriteConfig(); err != nil {
		return errors.Wrap(err, "in writing config file")
	}

	return nil
}

// touch is kind of like /bin/touch.
func touch() {
	f, err := os.OpenFile(configName+"."+configType, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
}
