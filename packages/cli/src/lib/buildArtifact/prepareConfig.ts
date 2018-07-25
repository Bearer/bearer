import * as globby from 'globby'
import * as vm from 'vm'
import * as fs from 'fs'
import { promisify } from 'util'
const readFileAsync = promisify(fs.readFile)

export default (
  authConfigFile: string,
  distPath: string,
  scenarioUuid: string,
  nodeModulesPath: string
): Promise<{ intents: Array<{ [key: string]: string }>; integration_uuid: string }> => {
  module.paths.push(nodeModulesPath)

  return globby([`${distPath}/*.js`]).then(files =>
    files
      .reduce(async (acc, f) => {
        const code = await readFileAsync(f)
        const context = vm.createContext({ module: {} })

        vm.runInNewContext(code.toString(), context)
        const intent = context['module'].exports.default

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
          const content = await readFileAsync(authConfigFile, { encoding: 'utf8' })
          config.auth = JSON.parse(content)
          return config
        } catch (e) {
          throw new Error(`Unable to read ${authConfigFile} : ${e.toString()}`)
        }
      })
  )
}
