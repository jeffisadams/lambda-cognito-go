# Supporting codebase for blog tutorial
Visit (Here)[https://tenmilesquare.com/aws-sam-api-with-cognito/] to view the context.


## To run
- Set the following environment variables
  - STACK_NAME
  - STACK_BUCKET (this bucket must exist in your AWS environment)
  - YOUR_EMAIL (The template will create a user at this email and send you a temp password)
- run `make deploy`
- run `./scripts/login_first.sh {{User Pool ID}} {{User Pool Client ID}} {{Your Email}} {{Temp password that was sent to you}}`
  - Get the first two values from the cloudformation outputs dashboard
- run `curl {{Url to your api}}/open`
  - Will work before logging in
- run `curl -H "Authorization: {{Auth Token from script above}}" {{Url to your api}}`
  - Will only work if you add the Authorization Header

PRs and suggestions are welcome.