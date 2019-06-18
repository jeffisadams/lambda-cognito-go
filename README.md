# lambda-cognito-go
- A basic example of a lambda with AWS Cognito authentication and Sam Deployment process

## Intro
It's official! AWS has decided that Lambdas are our hammer, and we've been wandering around looking for nails ever since.  It's a compelling use case.  It's cheap to run, easy-ish to maintain, no infrastructure, and you can run scalable code as a function in the cloud.  In my quest to learn this new pattern, I came across many tutorials with part of the API Gateway story, but few that told the whole story from code to build to deploy.  Even fewer still that told this story collected with Authentication.  So let's build and document an API Gateway with Auth using AWS.  My intent is to focus on templates and code as infrastructure to stand up a Cognito User pool, add the authorizer to the protect the endpoint.

## Pre-requisites
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html)
- [AWS Serverless Application Model CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
- [Go](https://golang.org/doc/install)
- [Command line access to your environment](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_users_create.html)


## Some qualifications
- This is not the Front end, so retreiving a token and logging in will need to be done elsewhere.
- The plan is to use the baked in Authorizer as the lambda function.  Making a custom usage of Cognito in order to do analytics or other customized functionality is possible, but that's for another day.

## The Template
We want to be able to build from zero to stack repeatedly.  One frustration I have seen are templates showcasing the Ease of lambda deployment with cloudformation and SAM but no connection to the Authorization or other parts of the stack.  My goal is this is the entire back end ready to run.



## Recommended Reading
- [Sam Policies](https://github.com/awslabs/serverless-application-model/blob/master/docs/policy_templates.rst)

