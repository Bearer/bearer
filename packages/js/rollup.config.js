import typescript from 'rollup-plugin-typescript2'
import resolve from 'rollup-plugin-node-resolve'
import commonjs from 'rollup-plugin-commonjs'
import replace from 'rollup-plugin-replace'
import { terser } from 'rollup-plugin-terser'
import filesize from 'rollup-plugin-filesize'

import pkg from './package.json'

const isProduction = process.env.NODE_ENV === 'production'

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

const variables = {
  production: {
    INTEGRATION_HOST_URL: 'https://int.bearer.sh'
  },
  development: {
    INTEGRATION_HOST_URL: 'http://localhost:3000/'
  }
}

const createConfiguration = ({ mode = 'production', input, output }) => ({
  input,
  output: {
    ...output,
    sourcemap: true
  },
  plugins: [
    ...plugins,
    replace({
      BEARER_VERSION: pkg.version,
      INTEGRATION_HOST_URL: variables[mode].INTEGRATION_HOST_URL || 'https://int.bearer.sh',
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
  })
]
const distributionBundles = [
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
  }),
  createConfiguration({
    mode: 'development',
    input: 'src/index.ts',
    output: {
      file: pkg.unpkgDev,
      name: 'bearer',
      format: 'umd'
    }
  })
].concat(bundles)
export default (isProduction ? distributionBundles : bundles)
