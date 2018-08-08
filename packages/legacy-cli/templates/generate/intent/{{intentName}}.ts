import { {{intentType}}, T{{authType}}Context } from '@bearer/intents'
import Client from './client'

export default class {{intentName}}Intent {
  static intentName: string = '{{intentName}}'
  static intentType: any = {{intentType}}

  {{actionExample}}
}

