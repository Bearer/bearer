import typescript from 'rollup-plugin-typescript2'
import pkg from './package.json'
import commonjs from 'rollup-plugin-commonjs'
import resolve from 'rollup-plugin-node-resolve'
import { terser } from 'rollup-plugin-terser'
import builtins from 'rollup-plugin-node-builtins'

export default [
  // browser-friendly UMD build
  {
    input: 'src/logger.ts',
    output: {
      name: 'logger',
      file: pkg.browser,
      format: 'iife'
    },
    plugins: [typescript(), resolve(), commonjs(), terser()]
  },

  // CommonJS (for Node) and ES module (for bundlers) build.
  // (We could have three entries in the configuration array
  // instead of two, but it's quicker to generate multiple
  // builds from a single configuration where possible, using
  // an array for the `output` option, where we can specify
  // `file` and `format` for each target)
  {
    input: 'src/logger.ts',
    output: { file: pkg.module, format: 'es' },
    plugins: [typescript(), builtins(), resolve(), commonjs(), terser()]
  },
  {
    input: 'src/logger.ts',
    output: { file: pkg.main, format: 'cjs' },
    plugins: [typescript(), resolve(), commonjs(), terser()]
  }
]
