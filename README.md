# go-jira-migrate
a golang tool to migrate jira project

Read jira issues from SRC, apply some transformations like google translate and create issues on DST.

# install

go get github.com/bregydoc/gtranslate
go mod init test
go mod tidy

go run migrate.go

