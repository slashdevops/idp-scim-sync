# AWS SAM Template

This project use [AWS SAM](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-specification.html), a tool for creating serverless applications.  This facilitates the creation of serverless applications that can be deployed to [AWS Lambda](https://aws.amazon.com/lambda/) and others different resources using [AWS CloudFormation.](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/Welcome.html).

## Resources

The  file [template.yaml](https://github.com/slashdevops/idp-scim-sync/blob/main/template.yaml) after being transformed by the [AWS SAM package command](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-cli-command-reference-sam-package.html) in `packaged.yaml` creates the following resources:

| #   | Type                        | CloudFormation Logical ID                  |
| --- | --------------------------- | ------------------------------------------ |
| 1   | AWS::SecretsManager::Secret | AWSGWSServiceAccountFileSecret             |
| 2   | AWS::SecretsManager::Secret | AWSGWSUserEmailSecret                      |
| 3   | AWS::SecretsManager::Secret | AWSSCIMAccessTokenSecret                   |
| 4   | AWS::SecretsManager::Secret | AWSSCIMEndpointSecret                      |
| 5   | AWS::S3::Bucket             | Bucket                                     |
| 6   | AWS::S3::BucketPolicy       | BucketPolicy                               |
| 7   | AWS::KMS::Key               | KMSKey                                     |
| 8   | AWS::KMS::Alias             | KMSKeyAlias                                |
| 9   | AWS::Lambda::Function       | LambdaFunction                             |
| 10  | AWS::Logs::LogGroup         | LambdaFunctionLogGroup                     |
| 11  | AWS::IAM::Role              | LambdaFunctionRole                         |
| 12  | AWS::Events::Rule           | LambdaFunctionSyncScheduledEvent           |
| 13  | AWS::Lambda::Permission     | LambdaFunctionSyncScheduledEventPermission |
