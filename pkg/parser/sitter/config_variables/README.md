## How to develop

### Installation steps

```bash
brew install tree-sitter
```

### How to test

Update the `grammar.js`

Build the grammar:

```bash
yarn build
```

Run the tests (build first):

```bash
yarn test
```

Update Go library (build first):

```bash
yarn update
```
