# Protocol Buffer Parser

This is using https://github.com/mitchellh/tree-sitter-proto to parse Protobuf files and get the tree structure.

## How to develop

### Installation steps

```bash
git clone https://github.com/mitchellh/tree-sitter-proto.git
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
