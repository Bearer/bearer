import { Config } from '@stencil/core'
const { plugins } = require('@bearer/core/dist/plugins')

export const config: Config = {
  namespace: '{{componentTagName}}',
  outputTargets: [
    {
      type: 'dist'
    },
    {
      type: 'www',
      serviceWorker: null,
      resourcesUrl: process.env.CDN_HOST,
    }
  ],
  plugins: [...plugins()],
  globalScript: 'src/global/context.ts'
}