import typescript from 'rollup-plugin-typescript2'
import resolve from 'rollup-plugin-node-resolve'
import commonjs from 'rollup-plugin-commonjs'
import replace from 'rollup-plugin-replace'
import { terser } from 'rollup-plugin-terser'
import filesize from 'rollup-plugin-filesize'

import pkg from './package.json'

const isProduction = process.env.NODE_ENV === 'production'

const mode = isProduction ? 'production' : 'development'

const plugins = [
  resolve({
    browser: true,
    preferBuiltins: false,
    extensions: ['.mjs', '.js', '.ts', '.jsx', '.tsx', '.json']
  }),
  commonjs(),
  typescript({
    exclude: ['**/__tests__/**', '**/node_modules/**']
  }),
  filesize({
    showMinifiedSize: isProduction,
    showGzippedSize: isProduction
  })
]

const createConfiguration = ({ input, output }) => ({
  input,
  output: {
    ...output,
    sourcemap: true
  },
  plugins: [
    ...plugins,
    replace({
      BEARER_VERSION: pkg.version,
      INTEGRATION_HOST_URL: process.env.INTEGRATION_HOST_URL || 'https://int.bearer.sh',
      __DEV__: mode === 'development',
      'process.env.NODE_ENV': JSON.stringify('production')
    }),
    mode === 'production' && terser()
  ].filter(Boolean)
})

const bundles = [
  createConfiguration({
    input: 'src/index.ts',
    output: {
      file: pkg.main,
      format: 'cjs'
    }
  }),
  createConfiguration({
    input: 'src/index.ts',
    output: {
      file: pkg.module,
      format: 'es'
    }
  }),
  createConfiguration({
    input: 'src/index.ts',
    output: {
      file: pkg.unpkg,
      name: 'bearer',
      format: 'umd'
    }
  })
]
export default bundles
