package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/slack-go/slack"
	yaml "gopkg.in/yaml.v2"
)

const (
	EnvSlackToken = "SLACK_TOKEN"
	EnvFilepath   = "INPUT_FILEPATH"
)

var userCache map[string]string

type UserGroups struct {
	Groups []struct {
		Name    string   `yaml:"name"`
		Members []string `yaml:"members"`
	}
}

func main() {
	userCache = map[string]string{}

	filepath := os.Getenv(EnvFilepath)
	if filepath == "" {
		fmt.Fprintln(os.Stderr, "File path is required")
		os.Exit(1)
	}

	slackToken := os.Getenv(EnvSlackToken)
	if slackToken == "" {
		fmt.Fprintln(os.Stderr, "Slack Token is required")
		os.Exit(1)
	}

	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error opening the file: #%v ", err)
	}

	ug := UserGroups{}
	err = yaml.Unmarshal(yamlFile, &ug)
	if err != nil {
		fmt.Printf("Error unmarshalling the file: %v", err)
	}

	api := slack.New(slackToken)

	slackGroups, err := api.GetUserGroups(slack.GetUserGroupsOptionIncludeUsers(true))
	if err != nil {
		fmt.Printf("Error getting Slack UserGroups: %v\n", err)
		return
	}

	for _, group := range ug.Groups {
		memberIDs := getListUserIDs(api, group.Members)

		// 1. Check if groups exist otherwise create it
		slackGroupId := ""
		for _, sg := range slackGroups {
			if sg.Handle == group.Name {
				slackGroupId = sg.ID
			}
		}

		if slackGroupId == "" {
			ug := slack.UserGroup{
				Name:   "App Squad " + group.Name,
				Handle: group.Name,
			}
			ugr, err := api.CreateUserGroup(ug)
			if err != nil {
				fmt.Printf("Error creating the Slack UserGroup: %v\n", err)
				return
			}
			slackGroupId = ugr.ID
		}

		// 2. Update the members
		_, err = api.UpdateUserGroupMembers(slackGroupId, strings.Join(memberIDs, ","))
		if err != nil {
			fmt.Printf("Error updating members for the Slack UserGroups: %v\n", err)
			return
		}
	}

	fmt.Println(`::set-output name=output::"Slack has been updated successfully"`)
}

func getListUserIDs(api *slack.Client, userList []string) []string {
	userIDs := []string{}

	for _, user := range userList {
		if cachedID, ok := userCache[user]; ok {
			userCache[user] = cachedID
			continue
		}

		userInfo, err := api.GetUserByEmail(user)
		if err != nil {
			fmt.Printf("Error getting users from Slack: %v\n", err)
			return userIDs
		}
		userIDs = append(userIDs, userInfo.ID)
		userCache[user] = userInfo.ID
	}

	return userIDs
}
