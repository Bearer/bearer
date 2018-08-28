import typescript from 'rollup-plugin-typescript2'
import resolve from 'rollup-plugin-node-resolve'
import commonjs from 'rollup-plugin-commonjs'
import replace from 'rollup-plugin-replace'
import { terser } from 'rollup-plugin-terser'
import copy from 'rollup-plugin-copy'

import { version } from './package.json'

const { parsed: parsedConfig } = require('dotenv').config()

const isProduction = process.env.NODE_ENV === 'production'

if (isProduction) {
  const { parsed: parsedSample } = require('dotenv').config({
    path: '.env.example'
  })

  const requiredKeys = new Set(Object.keys(parsedSample || {}))

  const setEquality = (set1, set2) => set1.size === set2.size && Array.from(set1).every(item => set2.has(item))

  const configuredKeys = new Set(Object.keys(parsedConfig || {}).filter(key => parsedConfig[key]))
  if (!setEquality(requiredKeys, configuredKeys)) {
    console.warn('Missing configuration, please check .env.* files')
  }
}

function plugins() {
  const base = [
    commonjs(),
    typescript({
      exclude: ['*.d.ts', '**/*.d.ts', '**/plugins.ts', '**/node_modules/**']
    }),
    resolve({
      jsnext: true,
      main: true,
      browser: true
    }),
    replace({
      LIB_VERSION: version,
      'process.env.BUILD': JSON.stringify(process.env.BUILD)
    })
  ]
  return isProduction ? [...base, terser()] : base
}

const bundles = [
  {
    input: 'src/index.ts',
    output: [
      {
        file: 'lib/main.es.js',
        format: 'es'
      },
      {
        file: 'lib/main.js',
        format: 'cjs',
        exports: 'named'
      },
      {
        file: 'lib/main.browser.js',
        name: 'Bearer',
        format: 'iife'
      }
    ],
    plugins: plugins()
  },
  {
    input: 'src/plugins.ts',
    output: {
      file: 'lib/plugins.js',
      format: 'cjs'
    }
  }
]

export default bundles
