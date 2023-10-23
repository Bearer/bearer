---
title: Enable Completion Script
---

# Enable shell completion

Below is example steps to enable shell completion feature for `bearer` cli:

## 1. Know your current shell

```bash
$ echo $SHELL
/bin/zsh # For this example it is zsh, but will be vary depend on your $SHELL, maybe /bin/bash or /bin/fish
```

## 2. Run `completion` command to get sub-commands

``` bash
bearer completion -h
```

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

```bash
echo "autoload -U compinit; compinit" >> ~/.zshrc
```

To load completions in your current shell session:

```bash
source <(bearer completion zsh); compdef _bearer bearer
```

To load completions for every new session, execute once:

### Linux

```bash
bearer completion zsh > "${fpath[1]}/_bearer"
```

### MacOS

```bash
bearer completion zsh > $(brew --prefix)/share/zsh/site-functions/_bearer
```

You will need to start a new shell for this setup to take effect.

## 3. Run the sub-commands following the instruction

```bash
bearer completion zsh > "${fpath[1]}/_bearer"
```

## 4. Start a new shell and you can see the shell completion

```bash
bearer [tab]
```
