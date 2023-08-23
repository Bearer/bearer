# GraphQL Parser

This is using https://github.com/bkegley/tree-sitter-graphql.git to parse GraphQL files and get the tree structure.

## How to develop

### Installation steps

```bash
git clone https://github.com/bkegley/tree-sitter-graphql.git
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
