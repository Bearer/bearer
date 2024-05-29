# PHP Parser

This is using a different version of the PHP parser to the Go Tree Sitter
version, as we wanted newer features.

Source: https://github.com/tree-sitter/tree-sitter-php
Version SHA: 92a98adaa534957b9a70b03e9acb9ccf9345033a

## Updating

You must rename `tree_sitter_php*` symbols to `tree_sitter_php2*` to avoid
conflicts with the Go Tree Sitter bundled parser.
