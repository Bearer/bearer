const { plugins } = require('@bearer/core/dist/plugins')
const transformers = require('@bearer/core/dist/transformers/index').default

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
      baseUrl: '/prerender',
    }
  ],
  customTransformers: {
    prependBefore: transformers({ verbose: true })
  },
  plugins: [...plugins()]
}

exports.devServer = {
  root: 'www',
  watchGlob: '**/**'
}
