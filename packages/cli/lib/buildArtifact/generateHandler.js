module.exports = ({ intents }) => {
  return intents
    .map(Object.keys)
    .map(
      intent => `
const ${intent} = require("../../intents/${intent}");
module.exports[${intent}.intentName] = ${intent}.intentType.intent(${intent}.action);
`
    )
    .join('\n')
}
