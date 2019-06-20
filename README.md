# Intro
It's official! AWS has decided that Lambdas are our hammer, and we're all wandering around looking for nails.  It's a compelling use case.  It's cheap to run, easy-ish to maintain, no infrastructure, and you can run scalable code as a function in the cloud.  Lots of tutorials exist to get a hello world function running using various tools both coding and AWS UI related.  These tutorials often leave out the ability to create a central API Gateway for a set of functions, and leave out how to protect your API with a basic Authentication layer.  My goal is to show a "hello *Whole Wide* world" example that includes some of these details.

## TLDR
Basic Repo is [here](https://github.com/jeffisadams/lambda-cognito-go).  PRs and suggestions welcome.

# What we are building
- An API Gateway
- A Cognito User Pool to restrict access to one of our functions
    - A Cognito User
- A simple funcion that is *NOT* protected by our auth layer
- A simple funcion that is protected by our created auth layer


# Pre-requisites
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html)
- [AWS Serverless Application Model CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
- [Go](https://golang.org/doc/install)
- [Command line access to your environment](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_users_create.html)

# Dive In
You can download the repo, set the needed variables (STACK_NAME, STACK_BUCKET, YOUR_EMAIL) and run `make deploy` to see this in action.  We want to be able to build from zero to stack.  I have found that AWS is a sensitive beast and will require continual iteration around subtle details.  I cannot stress enough the need to have code that you can run repeatedly in order to step through these iterations methodically. Take it from someone who has lost days of work to this phenomenon that it is worth setting up code from the start.

## UserPool / Client / User
```
UserPool:
    Type: AWS::Cognito::UserPool
    Properties:
        AdminCreateUserConfig:
        AllowAdminCreateUserOnly: false
        UserPoolName: TestingUsers
        UsernameAttributes:
        - email
        AutoVerifiedAttributes:
        - email
        Policies:
        PasswordPolicy:
            MinimumLength: 6
            RequireLowercase: true
            RequireNumbers: false
            RequireSymbols: false
            RequireUppercase: true
UserPoolTokenClient:
    Type: AWS::Cognito::UserPoolClient
    Properties:
      UserPoolId: !Ref UserPool
      GenerateSecret: false
      ExplicitAuthFlows:
      - USER_PASSWORD_AUTH
```
We create a userpool and a user pool client.  The pool is the abstract collection of users and their info.  The client is the ability to login using the SDK or the CLI.  You may need additional clients (We don't yet have Oauth) and additional properties, but this is a working minimum set that works.  Take a look at the [Cloudformation Reference Docs](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-cognito-userpool.html) for more details.  For logging in, we will be using the AWS CLI.
```
UserPoolUser:
    Type: AWS::Cognito::UserPoolUser
    Properties:
      DesiredDeliveryMediums:
        - EMAIL
      Username: !Ref YourEmail
      UserPoolId: !Ref UserPool
```
We also create your first user using a set Environment variable YOUR_EMAIL.  This will send your email a temporary password on stack creation.  Below are instructions for how we will login (spoiler, it's with the CLI).

## API
Technically we don't need this.  Sam will create one for us.  But then when you have two functions, you have two full APIs.  We can do better.
```
ServiceApi:
    DependsOn: UserPool
    Type: AWS::Serverless::Api
    Properties:
      Name: ServiceApi
      StageName: !Ref Version
      Cors:
        AllowMethods: "'*'"
        AllowHeaders: "'*'"
        AllowOrigin: "'*'"
      Auth:
        Authorizers:
          CognitoAuthorizer:
            UserPoolArn: !GetAtt "UserPool.Arn"
```
Once created, we use the API ID to attach the created functions in one logical group.  I have also set Cors headers leaving this wide open.  Do not do this unless you understand the implications.  Now when we create our functions we can pool them together under this API and have a more organized Microservice instead of a collection of functions.

## Functions
Pretty basic declaration.  Here is the Unauthenticated Function

```
OpenLambdaFunction:
    Type: AWS::Serverless::Function
    Properties:
      Description: Handles the basic request with no need for authentication
      Runtime: go1.x
      Handler: ./dist/open
      Events:
        Get:
          Type: Api
          Properties:
            Path: /open
            RestApiId: !Ref ServiceApi
            Method: GET
```
And the Authenticated.  They look very similar, but I wanted different code to handle each.  There are good use cases for both merging into events vs separate functions.
```
LambdaFunction:
    Type: AWS::Serverless::Function
    Properties:
      Description: Handles the basic request
      Runtime: go1.x
      Handler: ./dist/authenticated
      Events:
        Get:
          Type: Api
          Properties:
            Path: /
            RestApiId: !Ref ServiceApi
            Method: GET
            Auth:
              Authorizer: CognitoAuthorizer
```
The Function specifies the API Gateway to file under, the Authorizer to use, and the path / method to respond to.

## Code
This is arguably the simplest part.  Just send back a 200.  This code is basically the same for both, but with payload content tweaks.
```
return events.APIGatewayProxyResponse{
    StatusCode: 200,
    Body:       "{\"message\":\"Marshalling a return body is a problem for another day.  But the request was successful\",\"structure\":\"See this is actually not an error\"}",
    // This is important as part of the CORS config.
    // Again you should know the security implications of CORS before implementing this
    Headers: map[string]string{
        "Content-Type":                 "application/json",
        "Access-Control-Allow-Origin":  "*",
        "Access-Control-Allow-Methods": "*",
        "Access-Control-Allow-Headers": "*",
    },
}, nil
```


# Build and Deploy
Since we are using GO, we need a build and compile process. I am using a basic Makefile to compile using GO and run the AWS SAM Cli commands.

You will need to set the following variables:
- STACK_NAME
    - What to call the cloudformation Stack
- STACK_BUCKET
    - Sam uploads your compiled code resources to a bucket.  I'm leaving it to you to create a bucket and set this environment variable
- YOUR_EMAIL
    - Put in a valid email address for the first user to create.  You can do this through the AWS portal, but the focus is again on a full code auth example.


Once set, run `make deploy`.  Assuming you have access to your AWS environment, you'll see the build process compile the code, upload it to the bucket while transpiling the SAM template into an AWS cloudformation template, and deploying the stack.  This will work for updates as well.  You will want to view the outputs from the stack creation in order to get the ids needed for login, and the API url to call.  I use the web portal for this purpose, but you can also access the output with the CLI.

# I have an API, now what?
## Let's call the open endpoint.
```
curl https://{{Your API ID}.execute-api.us-east-1.amazonaws.com/v1/open
{ message: "This endpoint does not require any authentication", structure: "This field was added just to prove it's not an error" }
```
This works!  Now let's call the Authorized endpoint
```
curl https://{{Your API ID}.execute-api.us-east-1.amazonaws.com/v1
{"message":"Unauthorized"}
```
Technically this is a good thing, but we can do better.  We need to login.  We created a token client that will respond to SDK / CLI requests to log in.  This is arguably less secure, but allows us to login without additional infrastructure.

## Login and get a token
We will use the AWS cli to login.  The first login will require changing the password and follow a challenge workflow.  To make this slightly less painful, I crteated a script you can call that will log in and run the password challenge response.
```
AUTH_CHALLENGE_SESSION=`aws cognito-idp initiate-auth \
--auth-flow USER_PASSWORD_AUTH \
--auth-parameters "USERNAME={{YOUR_EMAIL}},PASSWORD={{password from the email AWS sent you}}" \
--client-id {{Token Client ID}} \
--query "Session" \
--output text`

# Then respond to the challenge
aws cognito-idp admin-respond-to-auth-challenge \
--user-pool-id $USERPOOLID \
--client-id {{Token Client ID}} \
--challenge-responses "NEW_PASSWORD=Testing1,USERNAME={{YOUR_EMAIL}}" \
--challenge-name NEW_PASSWORD_REQUIRED \
--session $AUTH_CHALLENGE_SESSION
```

To use this script, get the output values from your cloudformation stack and run the following command:
`./scripts/login_first.sh {{User Pool ID}} {{Token Client ID}} {{Your Email}} {{Password in the Email AWS sent you}}`

This will change your password to 'Testing1' and log you in.  You will get back a JSON Web Token or JWT token you can now use to finally call the damn API.

# Bring it home
Now that we have the auth token, we can add it to the headers and call the 
```
curl -H "Authorization: {{AUTH_TOKEN (The output from the login_first script)}}" https://{{Your API ID}}.execute-api.us-east-1.amazonaws.com/v1

{"message":"Marshalling a return body is a problem for another day.  But the request was successful","structure":"See this is actually not an error"}
```

# What have we accomplished?
So we now went from zero to:
- Cognito UserPool Authorizer
- API Gateway where we can put multiple functions
- A function that does not require authorization at path /open
- A function that requires authorization at path /
We can login using the AWS CLI / the login script `./scripts/login.sh {{UserPool Client ID}} {{Your Email}} Testing1` and add the output IdToken to our request in order to call our API.

## Recommended Reading
- [Sam Examples](https://github.com/awslabs/serverless-application-model/tree/master/examples/2016-10-31)
- [Cognito Docs](https://docs.aws.amazon.com/cognito/latest/developerguide/what-is-amazon-cognito.html)


