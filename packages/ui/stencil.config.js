const sass = require('@stencil/sass')
const replace = require('rollup-plugin-replace')

exports.config = {
  namespace: 'bearer-ui',
  plugins: [
    sass({
      injectGlobalPaths: ['src/globals/base.scss']
    }),
    replace({
      'process.env.BUILD': JSON.stringify(process.env.BUILD)
    })
  ],
  outputTargets: [
    {
      type: 'dist'
    },
    {
      type: 'www',
      serviceWorker: false
    }
  ]
}

exports.devServer = {
  root: 'www',
  watchGlob: '**/**'
}
