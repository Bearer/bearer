const prepareConfig = require('./prepareConfig')

xtest('config preparation', () => {
  let handlerContent = `
    const fs = require('fs');
    module.exports.intentOne = () => {};
    module.exports.intentTwo = () => {};
  `

  const config = prepareConfig(handlerContent)
  expect(config.intents).toEqual([
    { intentOne: 'index.intentOne' },
    { intentTwo: 'index.intentTwo' }
  ])
  expect(config.integration_uuid).toBeDefined()
})
