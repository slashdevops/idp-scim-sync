# AWS SAM Template

This project uses [AWS SAM](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-specification.html) to define the Lambda-based deployment of `idp-scim-sync`. The source template is [template.yaml](../template.yaml).

## What The Template Does

The template deploys the main `idpscim` program as a scheduled Lambda function that synchronizes Google Workspace groups and users into AWS IAM Identity Center through the SCIM API.

At a high level, the template does the following:

- Creates the Lambda function and schedules it with EventBridge.
- Stores sensitive credentials in Secrets Manager.
- Creates an encrypted S3 bucket for the sync state file.
- Creates the IAM role and permissions required by the Lambda function.
- Creates the CloudWatch log group used by the function.
- Adds metadata required to publish the application to the AWS Serverless Application Repository.

## Build And Publish Metadata

The template includes two important metadata blocks:

| Metadata | Purpose |
| --- | --- |
| `AWS::Serverless::Function` `BuildMethod: makefile` | Tells `sam build` to invoke `build-LambdaFunction` from [Makefile](../Makefile) |
| `AWS::ServerlessRepo::Application` | Provides the public application metadata used by the AWS Serverless Application Repository |

This means the SAM build for this project is not a generic Go Lambda build. It is repository-specific and intentionally compiles [cmd/idpscim/main.go](../cmd/idpscim/main.go) into a `bootstrap` binary.

## Parameters

The template parameters are grouped into three practical areas.

### Runtime And Sync Behavior

| Parameter | Purpose | Default |
| --- | --- | --- |
| `ScheduleExpression` | EventBridge schedule for the Lambda execution | `rate(15 minutes)` |
| `SyncMethod` | Sync strategy implemented by the application | `groups` |
| `SyncUserFields` | Optional user attributes to synchronize | empty |
| `GWSGroupsFilter` | Google Workspace group filter | empty |
| `LogLevel` | Application log level | `info` |
| `LogFormat` | Application log format | `json` |
| `MemorySize` | Lambda memory allocation | `256` |
| `Timeout` | Lambda timeout in seconds | `300` |
| `Runtime` | Lambda runtime | `provided.al2023` |
| `Architecture` | Lambda CPU architecture | `arm64` |
| `LambdaFunctionHandler` | Lambda handler entrypoint | `bootstrap` |
| `LambdaFunctionName` | Lambda function name | `idp-scim-sync` |
| `LogGroupName` | CloudWatch log group name | `/aws/lambda/idp-scim-sync` |
| `LogGroupRetentionDays` | CloudWatch log retention | `7` |
| `RoleNameSuffix` | Optional suffix to avoid IAM role name collisions | empty |

### State File Storage

| Parameter | Purpose | Default |
| --- | --- | --- |
| `BucketNamePrefix` | Prefix used to build the state bucket name | `idp-scim-sync-state` |
| `BucketKey` | S3 object key for the saved sync state | `data/state.json` |

The final bucket name is assembled as:

```text
<BucketNamePrefix>-<AWS Account ID>-<AWS Region>
```

### Credentials And Secret Names

| Parameter | Purpose | Default |
| --- | --- | --- |
| `GWSServiceAccountFile` | Google Workspace service account JSON contents | none |
| `GWSServiceAccountFileSecretName` | Secret name for that JSON | `IDPSCIM_GWSServiceAccountFile` |
| `GWSUserEmail` | Delegated Google Workspace user email | none |
| `GWSUserEmailSecretName` | Secret name for that email | `IDPSCIM_GWSUserEmail` |
| `SCIMEndpoint` | AWS IAM Identity Center SCIM endpoint | none |
| `SCIMEndpointSecretName` | Secret name for the SCIM endpoint | `IDPSCIM_SCIMEndpoint` |
| `SCIMAccessToken` | AWS IAM Identity Center SCIM access token | none |
| `SCIMAccessTokenSecretName` | Secret name for the SCIM token | `IDPSCIM_SCIMAccessToken` |

The sensitive input values are stored in Secrets Manager, and the Lambda function receives the secret names through environment variables. The function then resolves the actual secret values at runtime.

## Resources Created

After packaging and deployment, the template creates the following resources:

| # | Type | Logical ID | Purpose |
| --- | --- | --- | --- |
| 1 | `AWS::Serverless::Function` | `LambdaFunction` | Runs the sync logic from `idpscim` |
| 2 | `AWS::IAM::Role` | `LambdaFunctionRole` | Grants Lambda access to Secrets Manager, S3, KMS, logs, and X-Ray |
| 3 | `AWS::SecretsManager::Secret` | `AWSGWSServiceAccountFileSecret` | Stores Google Workspace credentials |
| 4 | `AWS::SecretsManager::Secret` | `AWSGWSUserEmailSecret` | Stores the delegated Google Workspace email |
| 5 | `AWS::SecretsManager::Secret` | `AWSSCIMEndpointSecret` | Stores the AWS SCIM endpoint |
| 6 | `AWS::SecretsManager::Secret` | `AWSSCIMAccessTokenSecret` | Stores the AWS SCIM token |
| 7 | `AWS::KMS::Key` | `KMSKey` | Encrypts the S3 state bucket |
| 8 | `AWS::KMS::Alias` | `KMSKeyAlias` | Stable alias for the KMS key |
| 9 | `AWS::S3::Bucket` | `Bucket` | Stores the sync state file |
| 10 | `AWS::S3::BucketPolicy` | `BucketPolicy` | Restricts access to the state bucket and enforces TLS |
| 11 | `AWS::Logs::LogGroup` | `LambdaFunctionLogGroup` | Stores Lambda logs |

The scheduled trigger is defined inside the `AWS::Serverless::Function` as a SAM `Schedule` event, which expands into the EventBridge rule and the matching Lambda invoke permission during transformation.

## Lambda Environment

The Lambda function receives configuration through environment variables mapped from template parameters and created secret names.

Important environment variables include:

- `IDPSCIM_LOG_LEVEL`
- `IDPSCIM_LOG_FORMAT`
- `IDPSCIM_SYNC_METHOD`
- `IDPSCIM_SYNC_USER_FIELDS`
- `IDPSCIM_GWS_GROUPS_FILTER`
- `IDPSCIM_AWS_S3_BUCKET_NAME`
- `IDPSCIM_AWS_S3_BUCKET_KEY`
- `IDPSCIM_GWS_USER_EMAIL_SECRET_NAME`
- `IDPSCIM_GWS_SERVICE_ACCOUNT_FILE_SECRET_NAME`
- `IDPSCIM_AWS_SCIM_ENDPOINT_SECRET_NAME`
- `IDPSCIM_AWS_SCIM_ACCESS_TOKEN_SECRET_NAME`

## Outputs

The template exports two useful outputs:

| Output | Description |
| --- | --- |
| `LambdaFunctionArn` | ARN of the deployed Lambda function |
| `LambdaFunctionName` | Name of the deployed Lambda function |

## Related Documentation

- [AWS-SAM.md](AWS-SAM.md) for deployment, update, and publish workflows
- [Configuration.md](Configuration.md) for application configuration behavior
- [README.md](../README.md) for high-level usage guidance
