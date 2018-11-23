# @bearer/transpiler

## What

A Typescript compiler with Bearer custom transformers built in.

## Why

Typescript compiler gives control on AST over processed Typescript files. Bearer leverage this feature to facilitate scenario writing.

Scenario frontend build process:

1. Bearer transpiler: uses scenarios views to generate Stencil compatible component
2. Stencil transpiler: use Bearer pre-processed files to generate web component

## Links

- Using the Typescript Compiler API https://github.com/Microsoft/TypeScript/wiki/Using-the-Compiler-API
- Built with StencilJS - [doc](https://stenciljs.com)
