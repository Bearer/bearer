import typescript from 'rollup-plugin-typescript2'
import resolve from 'rollup-plugin-node-resolve'
import commonjs from 'rollup-plugin-commonjs'
import { terser } from 'rollup-plugin-terser'
import * as dotenv from 'dotenv'

const { parsed: parsedConfig } = dotenv.config()

const isProduction = process.env.NODE_ENV === 'production'

function plugins() {
  const base = [
    commonjs(),
    typescript({}),
    resolve({
      jsnext: true,
      main: true,
      browser: true
    })
  ]
  return isProduction ? [...base, terser()] : base
}

const bundles = [
  {
    input: 'src/index.ts',
    output: [
      {
        file: 'dist/main.es.js',
        format: 'es'
      },
      {
        file: 'dist/main.js',
        format: 'cjs',
        exports: 'named'
      },
      {
        file: 'dist/main.browser.js',
        name: 'Bearer',
        format: 'iife'
      }
    ],
    plugins: plugins()
  }
]

export default bundles
