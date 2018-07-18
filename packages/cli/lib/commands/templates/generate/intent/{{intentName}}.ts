import { {{intentType}}, TContext } from '@bearer/intents'
import Client from './Client'

export default class {{intentName}}Intent {
  static intentName: string = '{{intentName}}'
  static intentType: any = {{intentType}}

  {{actionExample}}
}

