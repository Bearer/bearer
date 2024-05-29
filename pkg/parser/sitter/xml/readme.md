# XML Parser

An XML parser as the Go Tree Sitter library doesn't include one.

Source: https://github.com/dorgnarg/tree-sitter-xml
Version SHA: 36dd54f701e3b3e030412a854295af971cf74ad1

## Updating

At the time of writing, the current version of the grammar doesn't work unless
it's rebuilt with the latest tree sitter CLI. First, clone the grammar repo and
then (in the grammar repo):

```bash
$ tree-sitter generate
```
