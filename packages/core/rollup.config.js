import typescript from 'rollup-plugin-typescript2'
import resolve from 'rollup-plugin-node-resolve'
import commonjs from 'rollup-plugin-commonjs'
import replace from 'rollup-plugin-replace'
import { terser } from 'rollup-plugin-terser'
import strip from 'rollup-plugin-strip'

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
    // tslint:disable-next-line
    console.warn('Missing configuration, please check .env.* files')
  }
}

function plugins() {
  const base = [
    resolve({
      jsnext: true,
      main: true,
      preferBuiltins: false
    }),
    commonjs({
      namedExports: {
        // left-hand side can be an absolute path, a path
        // relative to the current directory, or the name
        // of a module in node_modules
        '../../node_modules/fbemitter/index.js': ['EventEmitter'],
        '../../node_modules/path/path.js': ['extname', 'resolve', 'sep']
      }
    }),
    typescript({
      exclude: ['*.d.ts', '**/*.d.ts', '**/plugins.ts', '**/node_modules/**']
    }),
    strip(),
    replace({
      LIB_VERSION: version,
      'process.env.BUILD': JSON.stringify(process.env.BUILD)
    })
  ]
  if (process.env.BUILD === 'distribution') {
    base.push(
      strip({
        include: ['**/*.js', '**/*.ts']
      })
    )
  }
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
        exports: 'named',
        file: 'lib/main.js',
        format: 'cjs'
      },
      {
        file: 'lib/main.browser.js',
        format: 'iife',
        name: 'Bearer'
      }
    ],
    plugins: plugins()
  },
  {
    input: 'src/plugins.ts',
    output: {
      file: 'lib/plugins.js',
      format: 'cjs'
    },
    plugins: plugins()
  }
]

export default bundles
