package main

import (
	"fmt"
	"os"
	"path/filepath"
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
	var filePath string
	userCache = map[string]string{}

	if len(os.Args) > 1 {
		filePath = os.Args[1]
	} else {
		filePath = os.Getenv(EnvFilepath)
	}
	if filePath == "" {
		fmt.Fprintln(os.Stderr, "File path is required")
		os.Exit(1)
	}

	slackToken := os.Getenv(EnvSlackToken)
	if slackToken == "" {
		fmt.Fprintln(os.Stderr, "Slack Token is required")
		os.Exit(1)
	}

	filePath = filepath.Clean(filePath)
	yamlFile, err := os.ReadFile(filePath)
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
				Name:   "Team " + strings.Title(strings.ToLower(group.Name)),
				Handle: strings.ToLower(group.Name),
			}
			ugr, err := api.CreateUserGroup(ug)
			if err != nil {
				fmt.Printf("Error creating the Slack UserGroup: %v\n", err)
				return
			}
			slackGroupId = ugr.ID
		}

		// 2. Update the members
		if len(memberIDs) > 0 {
			mlist := strings.Join(memberIDs, ",")

			fmt.Printf("Updating Slack UserGroup %v with members: %s\n", group.Name, mlist)
			_, err = api.UpdateUserGroupMembers(slackGroupId, mlist)
			if err != nil {
				fmt.Printf("Error updating members for the Slack group %s: %v\n", group.Name, err)
				return
			}
		}
	}

	fmt.Println("Slack has been updated successfully")
}

func getListUserIDs(api *slack.Client, userList []string) []string {
	userIDs := []string{}

	for _, user := range userList {

		if cachedID, ok := userCache[user]; ok {
			userIDs = append(userIDs, cachedID)
			continue
		}

		userInfo, err := api.GetUserByEmail(user)
		if err != nil {
			fmt.Printf("Error getting user %s from Slack API: %v\n", user, err)
			continue
		}
		userIDs = append(userIDs, userInfo.ID)
		userCache[user] = userInfo.ID
	}

	return userIDs
}
