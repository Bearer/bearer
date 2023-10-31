---
title: Enable Completion Script
---

# Enabling shell completions

Completions make it easier to use bearer on the command line. In this guide we will cover how to set everything up.

## Discover your current shell

```bash
$ echo $SHELL
/bin/zsh
```
In our example it's zsh, but will be vary depend on your $SHELL.

## Test the completion scripts

We currently support `zsh` `bash` and `fish`. To load completions in your current shell session you can run the following:

```bash
source <(bearer completion zsh); compdef _bearer bearer
```

Now you can test completions:

```bash
bearer [tab]
```

**Note:** if completions are not already enabled in your environment you will need
to enable it. In zsh you can execute the following once:

```bash
echo "autoload -U compinit; compinit" >> ~/.zshrc
```

## Final setup

To load completions for every new session you can add the test command to your `~/.zshrc` or similar. Alternatively if you are using zsh you can execute the following once:

### Linux

```bash
bearer completion zsh > "${fpath[1]}/_bearer"
```

### MacOS

```bash
bearer completion zsh > $(brew --prefix)/share/zsh/site-functions/_bearer
```

You will need to start a new shell for this setup to take effect.

