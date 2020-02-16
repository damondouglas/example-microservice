export project=$$(gcloud config get-value project)
export default-repo=gcr.io/${project}
export TAG=v0.0.1

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

stub: ## Stub proto
	protoc -I proto proto/job.proto --go_out=plugins=grpc:job/pkg/jobpb

build: ## Build containers
	skaffold build --default-repo=${default-repo}

deploy: ## Deploy artifacts
	skaffold run --default-repo ${default-repo}

expose: ## Port forwards controller
	kubectl port-forward svc/controller 8080:80
