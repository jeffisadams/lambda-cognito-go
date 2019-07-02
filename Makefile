# Set defaults if the Env vars aren't set
ifeq ($(YOUR_EMAIL),)
YOUR_EMAIL := 'test@example.com'
endif

ifeq ($(STACK_NAME),)
STACK_NAME := 'cognito-lambda-test-stack'
endif

ifeq ($(STACK_BUCKET),)
STACK_BUCKET := 'cognito-test-stack-plumbing-bucket'
endif

.PHONY: test
test:
	aws cloudformation validate-template --template-body file://template.yaml

.PHONY: clean
clean:
	rm -rf ./dist
	rm -rf template_deploy.yaml

.PHONY: deps
deps: clean
	go get github.com/aws/aws-lambda-go/events
	go get github.com/aws/aws-lambda-go/lambda

.PHONY: build
build: deps
	GOOS=linux go build -o dist/authenticated ./src/authenticated.go
	GOOS=linux go build -o dist/open ./src/open.go

.PHONY: api
api: build
	sam local start-api

.PHONY: deploy
deploy: build
	aws cloudformation package \
		--template-file template.yaml \
		--output-template template_deploy.yaml \
		--s3-bucket $(STACK_BUCKET)

	# aws s3 cp ./swagger.yaml s3://$(STACK_BUCKET)/lambda-cognito-go-api-def.yaml
	aws cloudformation deploy \
		--no-fail-on-empty-changeset \
		--template-file template_deploy.yaml \
		--stack-name $(STACK_NAME) \
		--parameter-overrides "ResourceBucket=$(STACK_BUCKET)" "YourEmail=$(YOUR_EMAIL)"

.PHONY: teardown
teardown:
	aws cloudformation delete-stack --stack-name $(STACK_NAME)
