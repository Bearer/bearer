import replace from 'rollup-plugin-replace'
import alias from 'rollup-plugin-alias'

export const plugins = () => {
  const basePLugins = [
    replace({
      BEARER_API_HOST: JSON.stringify(process.env.API_HOST),
      BEARER_SCENARIO_ID: process.env.BEARER_SCENARIO_ID,
      '@stencil/core': '@bearer/stencil-core'
    }),
    alias({
      '@stencil/core': './node_modules/@bearer/stencil-core'
    })
  ]

  return basePLugins
}

export default plugins
