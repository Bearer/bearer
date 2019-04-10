import { Config } from '@stencil/core'
const { plugins } = require('@bearer/core/lib/plugins')

export const config: Config = {
  namespace: 'bearer-' + process.env.BEARER_INTEGRATION_ID,
  enableCache: false,
  copy: [
    {
      src: 'bearer-manifest.json'
    }
  ],
  outputTargets: [
    {
      type: 'dist'
    },
    {
      type: 'www',
      serviceWorker: null,
      resourcesUrl: process.env.CDN_HOST
    }
  ],
  plugins: [...plugins()],
  globalScript: 'src/global/context.ts'
}
