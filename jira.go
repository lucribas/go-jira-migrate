package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/andygrunwald/go-jira"
	"github.com/tkanos/gonfig"
)

type Configuration struct {
	JiraUrl       string
	JiraUserName  string
	JiraUserToken string
}

func main() {

	configuration := Configuration{}
	err := gonfig.GetConf("secrets.json", &configuration)
	if err != nil {
		fmt.Printf(err.Error(), "")
		return
	}

	fmt.Printf("Start Jira with username: %s\n", configuration.JiraUserName)

	base := configuration.JiraUrl
	tp := jira.BasicAuthTransport{
		Username: configuration.JiraUserName,
		Password: configuration.JiraUserToken,
	}

	jiraClient, err := jira.NewClient(tp.Client(), base)
	if err != nil {
		fmt.Printf(err.Error(), "")
		return
	}

	issue, _, _ := jiraClient.Issue.Get("MSS-260", nil)
	fmt.Printf("------------------------------\n")
	fmt.Printf("Id :%s\n", issue.ID)
	fmt.Printf("Key :%s\n", issue.Key)
	fmt.Printf("Summary: %+v\n", issue.Fields.Summary)
	fmt.Printf("Type: %s\n", issue.Fields.Type.Name)
	fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)
	fmt.Printf("Description: %s\n", issue.Fields.Description)
	fmt.Printf("AccountID: %s\n", issue.Fields.Assignee.EmailAddress)
	fmt.Printf("Labels: %s\n", issue.Fields.Labels)

	// ctx := context.Background()
	// ctx = context.WithValue(ctx, "update", "{\"labels\": [{\"add\":\"triaged\"}}")
	// jiraClient.Issue.UpdateWithContext(ctx, issue)
	fmt.Printf("================================\n")

	type Labels struct {
		Add string `json:"add" structs:"add"`
	}

	type Update struct {
		Labels []Labels `json:"labels" structs:"labels"`
	}

	c := map[string]interface{}{
		"update": Update{
			Labels: []Labels{
				{
					Add: "my_label_HA",
				},
			},
		},
	}

	data, _ := json.Marshal(c)
	fmt.Printf("----->%v\n", string(data))

	resp, err := jiraClient.Issue.UpdateIssue(issue.ID, c)

	if err != nil {
		fmt.Println(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	issue, _, _ = jiraClient.Issue.Get("MSS-260", nil)

	fmt.Printf("------------------------------\n")
	fmt.Printf("------------------------------\n")
	fmt.Printf("Key :%s\n", issue.Key)
	fmt.Printf("Summary: %+v\n", issue.Fields.Summary)
	fmt.Printf("Type: %s\n", issue.Fields.Type.Name)
	fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)
	fmt.Printf("Description: %s\n", issue.Fields.Description)
	fmt.Printf("Labels: %s\n", issue.Fields.Labels)
}
