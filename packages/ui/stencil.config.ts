import { Config } from '@stencil/core'
import { postcss } from '@stencil/postcss'
import { sass } from '@stencil/sass'

import * as autoprefixer from 'autoprefixer'
import * as replace from 'rollup-plugin-replace'
import * as strip from 'rollup-plugin-strip'

const plugins: any = [
  sass({
    injectGlobalPaths: ['src/globals/base.scss']
  }),
  replace({
    'process.env.BUILD': JSON.stringify(process.env.BUILD)
  }),
  postcss({
    plugins: [autoprefixer()]
  })
]

if (process.env.BUILD === 'distribution') {
  plugins.push(strip({ include: ['**/*.js', '**/*.ts'] }))
}

export const config: Config = {
  namespace: 'bearer-ui',
  plugins,
  outputTargets: [
    { type: 'dist', dir: 'lib' },
    {
      type: 'www'
    }
  ]
}
