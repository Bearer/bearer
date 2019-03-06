import replace from 'rollup-plugin-replace'

import debug from './logger'

export const plugins = () => {
  const withVariables = {
    BEARER_INTEGRATION_ID: process.env.BEARER_INTEGRATION_ID,
    BEARER_INTEGRATION_HOST: process.env.BEARER_INTEGRATION_HOST || 'https://int.staging.bearer.sh/',
    BEARER_AUTHORIZATION_HOST: process.env.BEARER_AUTHORIZATION_HOST || 'https://int.staging.bearer.sh/'
  }

  if (process.env.BEARER_DEBUG) {
    debug.extend('plugins')('withVariables %j', withVariables)
  }

  return [replace(withVariables)]
}

export default plugins
