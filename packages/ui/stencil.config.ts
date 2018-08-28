import { Config } from '@stencil/core'
import { postcss } from '@stencil/postcss'
import { sass } from '@stencil/sass'

import * as autoprefixer from 'autoprefixer'
import * as replace from 'rollup-plugin-replace'

export const config: Config = {
  namespace: 'bearer-ui',
  plugins: [
    sass({
      injectGlobalPaths: ['src/globals/base.scss']
    }),
    replace({
      'process.env.BUILD': JSON.stringify(process.env.BUILD)
    }),
    postcss({
      plugins: [autoprefixer()]
    })
  ],
  outputTargets: [
    { type: 'dist' },
    {
      type: 'www'
    }
  ]
}
