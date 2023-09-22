# PHP Parser

This is using a different version of the PHP parser to the Go Tree Sitter
version, as we wanted newer features.

Source: https://github.com/tree-sitter/tree-sitter-php
Version SHA: 57f855461aeeca73bd4218754fb26b5ac143f98f

## How to develop

### Installation steps

```bash
git clone https://github.com/tree-sitter/tree-sitter-php.git
yarn install
```

### How to test

Update the `grammar.js`

Run the tests

```bash
yarn test
```

You can also use the playground provided by tree-sitter

```bash
yarn playground
```
