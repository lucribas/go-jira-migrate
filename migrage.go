package main

import (
	"fmt"
	"os"

	"github.com/andygrunwald/go-jira"
	"github.com/bregydoc/gtranslate"
	"github.com/gocarina/gocsv"
	"github.com/tkanos/gonfig"
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
	jiraClientSrc := openJira("secrets_src.json")
	// jiraClientDst := openJira("secrets_dst.json")

	issue, _, _ := jiraClientSrc.Issue.Get("MSS-260", nil)
	printIssue(issue)
	// issue, _, _ = jiraClientDst.Issue.Get("DMOB-1202", nil)
	// printIssue(issue)

	// Read csv
	issues := readCSV("Jira.csv")

	for _, issue := range issues {
		fmt.Printf("------------------------------\n")
		fmt.Println("** Processing Issue: ", issue.Key)
		// get Issue from SRC
		issue, _, _ := jiraClientSrc.Issue.Get(issue.Key, nil)
		// translate issue fields
		issue.Fields.Summary = translateSummary(issue.Fields.Summary)
		if issue.Fields.Description != "" {
			issue.Fields.Description = translateDescription(issue.Fields.Description)
		}
		printIssue(issue)
		// create Issue on DST
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

	return translated + "\n\n[pt-br original]\n" + text
}
