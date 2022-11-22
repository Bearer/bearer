# Curio Documentation

Welcome to the Curio docs. To view the full documentation, visit the [official docs site](https://bearer.github.io/curio).

## Contributing

Thanks for your interest in contributing to the docs! It's a great way to participate in the project without diving head-first into Curio's source. Check our [contributing guide](https://github.com/Bearer/curio/blob/main/CONTRIBUTING.md) and [code of conduct](https://github.com/Bearer/curio/blob/main/CODE_OF_CONDUCT.md) for general details on contributing.

You can make changes to the documentation without setting up and running the docs site by editing markdown files directly. If you want to add a more advanced change, see the development details below.

## Development

The documentation is with on [eleventy](https://www.11ty.dev) and [tailwindcss](https://tailwindcss.com/), and written in markdown. To get started, you'll need NodeJS 12 or newer.

Navigate to the docs directory in the Curio repo, and `npm install` to bootstrap the dependencies.

```bash
# github.com/bearer/curio/docs
npm install
```

With dependencies installed, you can start a local dev server by running the `start` command:

```bash
npm start
```

This will provide you with a local version of the doc site that refreshes whenever changes are made.