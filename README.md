# aws-sso-gws-sync

Sync your Google Workspace Groups and Users to AWS Single Sing-On

## Available Commands

### gws

```cmd
ssocli gws groups list -u "user.email@google.com" -s "./credentials.json" -q "name:Admin*" -q "name:SuperAdmin*"
```
