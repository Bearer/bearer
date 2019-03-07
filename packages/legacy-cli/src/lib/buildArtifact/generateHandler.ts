export default ({ functions }) => {
  return functions
    .map(Object.keys)
    .map(
      func => `
const ${intent} = require("./${intent}").default;
module.exports[${intent}.intentName] = ${intent}.intentType.intent(${intent}.action);
`
    )
    .join('\n')
}
