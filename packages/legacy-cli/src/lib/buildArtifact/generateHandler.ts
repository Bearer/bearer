export default ({ intents }) => {
  return intents
    .map(Object.keys)
    .map(
      func => `
const ${intent} = require("./${intent}").default;
module.exports[${intent}.intentName] = ${intent}.intentType.intent(${intent}.action);
`
    )
    .join('\n')
}
