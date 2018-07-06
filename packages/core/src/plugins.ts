import replace from 'rollup-plugin-replace'

export const plugins = () => {
  const basePLugins = [
    replace({
      BEARER_API_HOST: JSON.stringify(process.env.API_HOST),
      BEARER_SCENARIO_ID: process.env.BEARER_SCENARIO_ID
    })
  ]

  return basePLugins
}

export default plugins
