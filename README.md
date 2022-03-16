# go-jira-migrate

**A golang tool to migrate Jira issues from one server to another.**

## How it works

1. For each issue in `Jira.csv`:
1.1. Connect to `Jira source server` and get all details of issue.
1.2. Use *Google Translate API* to translate the *Summary* and *Description* from `pt-br` to `en`.
1.3. Fix the *IssueType* to the type available on `Jira destination server`.
1.4. Copy custom fields like *Story Points*.
1.5. Copy *Labels* field.
1.6. Add a label with the *Key of source issue*. That help you to find one issue using the source Key.
1.7. Add a label with the *name of Components of source issue*. Since the destination server have other components we add them as labels to preserve information.
1.8. Create a new issue on `Jira destination server`.

## Setup

Generate a list of issues using a `Jira SQL query` like:
```sql
project = MSS AND Sprint in futureSprints() ORDER BY Key ASC
```
Export the issues to a *csv file* as `Jira.csv`.

Create two *Jira configuration files* as bellow:

`secrets_src.json`
```json
{
	"JiraUrl": "https://your_src_jira_domain.atlassian.net",
	"JiraUserName":"your_id@mmm.com",
	"JiraUserToken":"dlfhsaldflsdafg9245"
}
```

`secrets_dst.json`
```json
{
	"JiraUrl": "https://your_dst_jira_domain.atlassian.net",
	"JiraUserName":"your_id@mmm.com",
	"JiraUserToken":"779776869796gjvui6giu6n"
}
```
Notes:
- For atlassian.net cloud create a token and use it.
- For on premisse use your password on JiraUserToken.

Check and modify the function createFromIssue() on migrate.go to apply Issue transformations that you need.

*Tip: you can get the complete issue structure of your server just by downloading the json from an Issue.Self URL and using vscode to format it*


# install

```bash
go get github.com/bregydoc/gtranslate
go mod init test
go mod tidy
```



# run

```bash
go run migrate.go
```

# example

Issue **read** from source server:
```text
** Processing Issue:  MSS-216
Self :https://xxxxxxxxx.atlassian.net/rest/api/2/issue/90980
Key :MSS-216
Summary: [feature] Estudo e documentação para captura de métricas de utilização do app - Parte 1
Type: Story
Priority: Major
Description: Pesquisar soluções de codigo aberto para a captura de métricas de uso do App.
E Documentar.
Se houver mais do que um comparar os recursos para seleção.
Rodar em timebox até onde couber dentro de 5 ptos (1 semana).
a refencia é o Firebase (não é codigo aberto).
```

Issue **written** to destination server:
```text
CREATED:Self :https://xxxx.com/rest/api/2/issue/1627192
Key :DMOB-1364
Summary: [feature] Study and documentation to capture app usage metrics - Part 1
Type: Story
Priority: Lowest
Description: Research open source solutions for capturing App usage metrics.
And Document.
If there is more than one compare the features for selection.
Run in timebox as far as it will fit within 5 points (1 week).
the reference is Firebase (not open source).
----
*[pt-br original]*
Pesquisar soluções de codigo aberto para a captura de métricas de uso do App.
E Documentar.
Se houver mais do que um comparar os recursos para seleção.
Rodar em timebox até onde couber dentro de 5 ptos (1 semana).
a refencia é o Firebase (não é codigo aberto).

*Imported* from Jira Certi: https://xxxxxxxxx.atlassian.net/browse/MSS-216
```
