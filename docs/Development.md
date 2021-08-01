# Development references

- [Testing Code that depends on google.golang.org/api](https://github.com/googleapis/google-api-go-client/blob/master/testing.md)
- [Semantic Commit Messages](https://gist.github.com/joshbuchea/6f47e86d2510bce28f8e7f42ae84c716)

## Semantic Commit Messages

See how a minor change to your commit message style can make you a better programmer.

Format: `<type>(<scope>): <subject>`

`<scope>` is optional

### Example

```text
feat: add hat wobble
^--^  ^------------^
|     |
|     +-> Summary in present tense.
|
+-------> Type: chore, docs, feat, fix, refactor, style, or test.
```

More Examples:

- `feat`: new feature for the user, not a new feature for build script
- `fix`: bug fix for the user, not a fix to a build script
- `docs`: changes to the documentation
- `style`: formatting, missing semi colons, etc; no production code change
- `refactor`: refactoring production code, eg. renaming a variable
- `perf`: A code change that improves performance
- `test`: adding missing tests, refactoring tests; no production code change
- `build`: Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)
- `ci`: Changes to our CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs)
- `chore`: updating grunt tasks etc; no production code change
- `license`: Edits regarding licensing; no production code change
- `revert`: Reverts a previous commit
- `bump`: Increase the version of something e.g. dependency
- `make`: Change the build process, or tooling, or infra
- `localize`: Translations update

References:

- https://www.conventionalcommits.org/
- https://seesparkbox.com/foundry/semantic_commit_messages
- http://karma-runner.github.io/1.0/dev/git-commit-msg.html

## Development

### Mocks

```bash
mockgen -package=sync -destination=internal/sync/service_mocks.go -source=internal/sync/service.go IdentityProviderService,SCIMService
```
