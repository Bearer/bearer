import { getPropertyValue, IntentCodeProcessor, isIntentClass } from '@bearer/cli/src/utils/generators'

import * as fs from 'fs-extra'
import * as ts from 'typescript'

type TConfig = {
  intents: Array<{ [key: string]: string }>
  integration_uuid: string
  auth?: any
}

const INTENT_NAME_IDENTIFIER = 'intentName'

export default (
  authConfigFile: string,
  _distPath: string,
  scenarioUuid: string,
  _nodeModulesPath: string,
  intentsDir: string
): Promise<TConfig> => {

  const intents: Array<any> = []

  const transformer = (context: ts.TransformationContext) => {
    function visit(tsNode: ts.Node) {
      if (isIntentClass(tsNode)) {
        const intentName = getPropertyValue(tsNode as ts.ClassDeclaration, INTENT_NAME_IDENTIFIER)
        intents.push(
          {
            [intentName]: `index.${intentName}`
          }
        )
      }
      return tsNode
    }
    return (tsSourceFile: ts.SourceFile) => {
      return ts.visitEachChild(tsSourceFile, visit, context)
    }
  }

  return new Promise((resolve, _reject) => {
    new IntentCodeProcessor(intentsDir, transformer).run()
    const content = fs.readFileSync(authConfigFile, { encoding: 'utf8' })
    let config: TConfig = { intents, integration_uuid: scenarioUuid, auth: JSON.parse(content) }
    resolve(config)
  })
}
