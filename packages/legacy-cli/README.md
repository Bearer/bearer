# Bearer CLI

## Installation

```bash
$ npm install @bearer/bearer-cli
```

## Usage

```bash
$ echo OrgId=4l1c3 > ~/.bearerrc
$ bearer new attachPullRequest && cd attachPullRequest
$ bearer generate searchRepositories --type GetCollection
$ bearer generate getPullRequest --type GetResource
$ bearer deploy
```

## Commands list

### Generators

**Generate a new integration\***

```bash
$ bearer new attachPullRequest
```

**Generate a new func of type collection**

```bash
$ bearer generate getRepositories --type=GetCollection
```

### Deploy

**Deploy implemented integration**

```bash
$ bearer deploy
```

To deploy to dev env use:

```bash
$ BEARER_ENV=dev bearer deploy
```
