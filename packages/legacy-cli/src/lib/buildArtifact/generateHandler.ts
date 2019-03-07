export default ({ functions }) => {
  return functions
    .map(Object.keys)
    .map(
      func => `
const ${func} = require("./${func}").default;
module.exports[${func}.intentName] = ${func}.intentType.call(${func}.action);
`
    )
    .join('\n')
}
