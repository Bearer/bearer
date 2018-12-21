# tsconfig

> Shared [TypeScript config](https://www.typescriptlang.org/docs/handbook/tsconfig-json.html) for bearer projects

## Install

```
$ yarn add -D  @bearer/tsconfig
```

## Usage

`tsconfig.json`

```json
{
  "extends": "@bearer/tsconfig",
  "compilerOptions": {
    "outDir": "dist",
    "lib": ["es2018"]
  }
}
```
