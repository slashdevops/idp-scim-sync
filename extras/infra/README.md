# Infrastructure

These Cloudformation templates (`#_cfn_*`) are used to create the infrastructure necessary to deploy the poroject into the AWS.

Why these is needed?

1. Allow users to consume the `Serverless Lambda function` from the public [AWS Serverless Application Repository](https://aws.amazon.com/es/serverless/serverlessrepo/)
2. Allow users to consume the public project [Docker image](https://gallery.ecr.aws/l2n7y5s7/slashdevops/idp-scim-sync) in [Amazon ECR](https://aws.amazon.com/es/ecr/) whitout the limits of public Docker repositories
3. Use the [DevOps Culture](https://aws.amazon.com/devops/what-is-devops/) to manage the infrastructure

## The cost of the infrastructure

As ususal this is hard to estimate, but the cost of the infrastructure is the same as the cost of the AWS resources.

## Resources used in AWS

Long live resources:

* IAM Identity providers, OIDC integratin with Github to get advantage of the CI/CD pipeline (Github Actions)
* AWS IAM for all the roles and policies of the Infrastructure we will describe here
* AWS S3 Bucket to store the CloudFormation templates and the AWS SAM Serverless Application artifacts
* AWS ECR Repository to store the Docker images
* AWS CloudFormation Stack to deploy the infrastructure described above

During the tests:

* AWS SSO and SCIM API calls
* AWS Lambda, metrics, logs, etc.
* AWS S3 to store the `state file` of the execution of the Lambda
* KMS to encrypt the `state file` stored in AWS S3
* AWS IAM Roles and Policies to allow the Lambda to access the resources

Fortunatelly most of the resources are free according to the [AWS pricing](https://aws.amazon.com/pricing/) and also I'll apply for [AWS Credits for Open Source Projects](https://pages.awscloud.com/AWS-Credits-for-Open-Source-Projects) to support the infrastructure.
