GO = go

test: ## Executes the tests
	${GO} test -race -cover ./...