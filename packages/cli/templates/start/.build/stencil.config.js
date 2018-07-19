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
  plugins: [...plugins()],
  globalScript: 'src/global/context.ts'
}
