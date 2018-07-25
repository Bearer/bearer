export default ({ intents }) => {
  return intents
    .map(Object.keys)
    .map(
      intent => `
const ${intent} = require("./dist/${intent}").default;
module.exports[${intent}.intentName] = ${intent}.intentType.intent(${intent}.action);
`
    )
    .join('\n')
}
