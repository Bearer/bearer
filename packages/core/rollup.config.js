import typescript from 'rollup-plugin-typescript2'
import resolve from 'rollup-plugin-node-resolve'
import commonjs from 'rollup-plugin-commonjs'
import replace from 'rollup-plugin-replace'
import { terser } from 'rollup-plugin-terser'
import strip from 'rollup-plugin-strip'

import pkg from './package.json'

const isProduction = process.env.NODE_ENV === 'production'

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
      LIB_VERSION: pkg.version,
      'process.env.BUILD': JSON.stringify(process.env.BUILD)
    })
  ]
  if (process.env.BUILD === 'distribution') {
    base.push(
      strip({
        include: ['**/*.js']
      })
    )
  }
  return isProduction ? [...base, terser()] : base
}

const bundles = [
  { input: 'src/index.ts', output: { file: pkg.module, format: 'es' }, plugins: [...plugins()] },
  { input: 'src/index.ts', output: { file: pkg.main, format: 'cjs' }, plugins: [...plugins()] },
  { input: 'src/index.ts', output: { file: pkg.browser, name: 'core', format: 'iife' }, plugins: [...plugins()] },
  { input: 'src/plugins.ts', output: { file: 'lib/plugins.js', format: 'cjs' }, plugins: plugins() }
]

export default bundles
