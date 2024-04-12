PROJECT_NAME := "notify-me"

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test: pre-commit-all ## Run all tests

setup: ## install required modules
	python -m pip install -U -r requirements-dev.txt
	pre-commit install

pre-commit-all: ## run pre-commit on all files
	pre-commit run --all-files

pre-commit: ## run pre-commit
	pre-commit run
