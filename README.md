# go-jira-migrate

**A golang tool to migrate Jira issues from one server to another.**

![Gopher](./gopher.png)

## How it works

For each issue in `Jira.csv`:
1. Connect to `Jira source server` and **get all details** of issue.
2. Use *Google Translate API* to **translate** the *Summary* and *Description* from `pt-br` to `en`.
3. **Fix** the *Project Name* to one predefined.
4. **Fix** the *Assignee* to one predefined.
5. **Map** the *IssueType* to the types available on `Jira destination server`.
6. **Copy** custom fields like *Story Points*.
7. **Copy** *Labels* field.
8. **Add a label** with the *Key of source issue*. This helps you to find a problem using the source key.
9. **Add a label** with the *name of Components of source issue*. As the target server has other components, we added them as labels to preserve the information.
10. **Create** a new issue on `Jira destination server`.

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
- For **atlassian.net cloud** create a token and use it on JiraUserToken.
- For **on premisse** use your password on JiraUserToken.

Check and modify the function **createFromIssue()** on `migrate.go` to apply issue transformations as need.

*Tip: you can get the **complete issue structure** of your server just by **downloading the json** from an `Issue.Self` URL and using **vscode** to format it*


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

Issue **read** from `Jira source server`:
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

Issue **written** to `Jira destination server`:
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
