# Bearer Core

Bearer Core contains helpers and business logic that you must use if you want create new components and interact with Bearer servers.

- Authentication
- Networks helpers
- Decorators [README](./src/README/md)

## Development

Build a local bundle

```bash
cp .env{.example,}

yarn build
yarn link
```

Somewhere else in a different repository

```bash
yarn link "@bearer/core"
```
