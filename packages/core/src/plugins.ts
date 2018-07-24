import replace from 'rollup-plugin-replace'

export const plugins = () => {
  const withVariables = {
    BEARER_SCENARIO_ID: process.env.BEARER_SCENARIO_ID,
    BEARER_INTEGRATION_HOST: process.env.BEARER_INTEGRATION_HOST || 'https://int.staging.bearer.sh/',
    BEARER_AUTHORIZATION_HOST: process.env.BEARER_AUTHORIZATION_HOST || 'https://int.staging.bearer.sh/'
  }

  if (process.env.BEARER_DEBUG) {
    console.log('[BEARER]', 'withVariables', withVariables)
  }

  return [replace(withVariables)]
}

export default plugins
