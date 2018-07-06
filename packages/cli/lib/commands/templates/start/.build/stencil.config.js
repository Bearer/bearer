const { plugins } = require('@bearer/core/dist/plugins')

exports.config = {
  namespace: '{{componentTagName}}',
  outputTargets: [
    {
      type: 'dist'
    },
    {
      type: 'www',
      serviceWorker: false,
      resourcesUrl: process.env.CDN_HOST,
      baseUrl: '/prerender'
    }
  ],
  plugins: [...plugins()]
}

exports.devServer = {
  root: 'www',
  watchGlob: '**/**'
}
