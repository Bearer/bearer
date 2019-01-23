import * as fs from 'fs-extra'

import * as ts from 'typescript'

import { getIntentName, IntentCodeProcessor, isIntentClass } from './generators'

type TConfig = {
  intents: string[]
  integration_uuid: string
  auth?: any
}

export const transformer = (intents: string[]) => (context: ts.TransformationContext) => {
  return (tsSourceFile: ts.SourceFile) => {
    function visit(tsNode: ts.Node) {
      if (isIntentClass(tsNode)) {
        const intentName = getIntentName(tsSourceFile)
        intents.push(intentName)
      }
      return tsNode
    }
    return ts.visitEachChild(tsSourceFile, visit, context)
  }
}

export default (authConfigFile: string, scenarioUuid: string, intentsDir: string): Promise<TConfig> => {
  return new Promise((resolve, reject) => {
    const intents: string[] = []
    new IntentCodeProcessor(intentsDir, transformer(intents))
      .run()
      .then(() => {
        const content = fs.readFileSync(authConfigFile, { encoding: 'utf8' })
        const config: TConfig = { intents, integration_uuid: scenarioUuid, auth: JSON.parse(content) }
        resolve(config)
      })
      .catch(error => {
        reject(error)
      })
  })
}
