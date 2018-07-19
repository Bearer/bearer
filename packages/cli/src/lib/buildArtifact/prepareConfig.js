const vm = require('vm')
const globby = require('globby')
const path = require('path')
const fs = require('fs')
const { promisify } = require('util')
const readFileAsync = promisify(fs.readFile)

const AUTH_CONFIG_FILE = 'auth.config.json'

module.exports = (codePath, scenarioUuid) => {
  const fullPath = path.resolve(codePath)

  module.paths.push(path.join(fullPath, 'node_modules'))

  return globby([`${fullPath}/dist/*.js`]).then(files =>
    files
      .reduce(async (acc, f) => {
        const code = await readFileAsync(f)
        const context = vm.createContext({ module: {} })
        vm.runInNewContext(code.toString(), context)
        const intent = context.module.exports.default

        if (intent && intent.intentName)
          acc.then(config =>
            config.intents.push({
              [intent.intentName]: `index.${intent.intentName}`
            })
          )
        return acc
      }, Promise.resolve({ integration_uuid: scenarioUuid, intents: [] }))
      .then(async config => {
        try {
          config.auth = JSON.parse(
            await readFileAsync(path.join(fullPath, '..', AUTH_CONFIG_FILE))
          )
          return config
        } catch (e) {
          console.log(e)
        }
      })
  )
}
