# idp-scim-sync

Sync your Google Workspace Groups and Users to AWS Single Sing-On using SCIM protocol.

## Available Commands

### gws

```cmd
idpscimcli gws groups list -u "user.email@google.com" -s "./credentials.json" -q "name:Admin*" -q "name:SuperAdmin*"
```
