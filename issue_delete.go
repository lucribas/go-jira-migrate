package main

import (
	// "context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/andygrunwald/go-jira"
	"github.com/bregydoc/gtranslate"
	"github.com/gocarina/gocsv"
	"github.com/tkanos/gonfig"
	"github.com/trivago/tgo/tcontainer"
)

type Configuration struct {
	JiraUrl       string
	JiraUserName  string
	JiraUserToken string
}

type CsvIssue struct {
	Label string `csv:"Label"` // .csv column headers
	Key   string `csv:"Issue_key"`
	Id    string `csv:"Issue_id"`
}

func main() {

	// open Jira Connections
	// jiraClientSrc := openJira("secrets_src.json")
	jiraClientDst := openJira("secrets_dst.json")

	// // Initial Test
	// issue, _, _ := jiraClientSrc.Issue.Get("MSS-260", nil)
	// printIssue(issue)
	// issue, _, _ = jiraClientDst.Issue.Get("DMOB-1202", nil)
	// printIssue(issue)

	// Read csv
	issuesCSV := readCSV("Jira.csv")

	for _, issue_csv := range issuesCSV {
		fmt.Printf("------------------------------\n")
		fmt.Println("** Processing Issue: ", issue_csv.Key)
		// get Issue from SRC
		issue, _, _ := jiraClientDst.Issue.Get(issue_csv.Key, nil)
		printIssue(issue)

		issue.Fields.Summary = "delete"
		// ctx := context.Background()
		// ctx = context.WithValue(ctx, "update", "{\"labels\": [{\"add\":\"triaged\"}}")
		// ctx = context.WithValue(ctx, "update", "{\"summary\":[{\"set\":\"delete\"}]]}")
		i := make(map[string]interface{})
		fields :=map[string]interface{}{"summary": "delete"}
		i["fields"] = fields
		jiraClientDst.Issue.UpdateIssue(issue_csv.Key, i)
		issue, _, _ = jiraClientDst.Issue.Get(issue_csv.Key, nil)
		printIssue(issue)
	}
}

func openJira(conf_file string) *jira.Client {

	configuration := Configuration{}
	err := gonfig.GetConf(conf_file, &configuration)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Start Jira with username: %s\n", configuration.JiraUserName)

	base := configuration.JiraUrl
	tp := jira.BasicAuthTransport{
		Username: configuration.JiraUserName,
		Password: configuration.JiraUserToken,
	}

	jiraClient, err := jira.NewClient(tp.Client(), base)
	if err != nil {
		panic(err)
	}

	return jiraClient

}

func readCSV(csv_file string) []*CsvIssue {
	in, err := os.Open(csv_file)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	issues := []*CsvIssue{}

	if err := gocsv.UnmarshalFile(in, &issues); err != nil {
		panic(err)
	}

	return issues
}

func printIssue(issue *jira.Issue) {
	fmt.Printf("Self :%s\n", issue.Self)
	fmt.Printf("Key :%s\n", issue.Key)
	fmt.Printf("Summary: %+v\n", issue.Fields.Summary)
	fmt.Printf("Type: %s\n", issue.Fields.Type.Name)
	fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)
	fmt.Printf("Description: %s\n", issue.Fields.Description)
	fmt.Printf("Labels: %s\n", issue.Fields.Labels)
}

func translateSummary(text string) string {
	translated, err := gtranslate.TranslateWithParams(
		text,
		gtranslate.TranslationParams{
			From: "pt",
			To:   "en",
		},
	)
	if err != nil {
		panic(err)
	}

	return translated
}
func translateDescription(text string) string {
	translated, err := gtranslate.TranslateWithParams(
		text,
		gtranslate.TranslationParams{
			From: "pt",
			To:   "en",
		},
	)
	if err != nil {
		panic(err)
	}

	header := "\n----\n*[pt-br original]*\n"
	return translated + header + text + "\n"
}

func createFromIssue(client *jira.Client, issueFrom *jira.Issue) {

	// Story points: import from custom field
	unknowns := tcontainer.NewMarshalMap()
	unknowns["customfield_10002"] = issueFrom.Fields.Unknowns["customfield_10002"]

	// Labels: original labels plus compent name
	labels := issueFrom.Fields.Labels
	labels = append(labels, issueFrom.Key) // add orign key
	for _, c := range issueFrom.Fields.Components {
		labels = append(labels, c.Name)
	}

	// IssueType mapping
	issueType := issueFrom.Fields.Type.Name
	if issueType == "Improvement" {
		issueType = "New Feature"
	}
	if issueType == "Sub-task" {
		issueType = "Task"
	}

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Assignee: &jira.User{
				Name: "lribas",
			},
			Description: issueFrom.Fields.Description,
			Type: jira.IssueType{
				Name: issueType,
			},
			Project: jira.Project{
				Key: "DMOB",
			},
			Summary:  issueFrom.Fields.Summary,
			Unknowns: unknowns,
			Labels:   labels,
		},
	}

	issue, resp, err := client.Issue.Create(&i)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf(string(body))
		}
		panic(err)
	}
	fmt.Printf("CREATED:")
	issue, _, _ = client.Issue.Get(issue.Key, nil)
	printIssue(issue)
}
