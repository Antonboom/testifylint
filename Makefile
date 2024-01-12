GOBIN= ${GOPATH}/bin

# https://taskfile.dev/installation/
${GOBIN}/task:
	go install github.com/go-task/task/v3/cmd/task@latest

.PHONY : install
install: ${GOBIN}/task

.PHONY : dev/tools
dev/tools: install
	task tools:install

.PHONY : clean
clean: 
	rm -f ${GOBIN}/task
