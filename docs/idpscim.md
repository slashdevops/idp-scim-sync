# idpscim

This is in charge of keeping Google Workspace Groups and Users sync with AWS Single Sing-On service using SCIM protocol, to do that in an efficient way it stores the last sync data into an `AWS S3 Bucket` and periodically it will check if there are changes.

This program could work in three ways:

1. As an [AWS Lambda function](https://aws.amazon.com/lambda/?nc1=h_ls) deployed via [AWS SAM](https://aws.amazon.com/serverless/sam/) or consumed directly from the [AWS Serverless Application Repository](https://aws.amazon.com/serverless/serverlessrepo/?nc1=h_ls)
2. As a command line tool
3. As a Docker container