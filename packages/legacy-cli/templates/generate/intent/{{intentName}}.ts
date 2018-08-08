import { {{intentType}}, T{{authType}}Context, {{callbackType}} } from '@bearer/intents'
// Uncomment this line if you need to use Client
// import Client from './client'

export default class {{intentName}}Intent {
  static intentName: string = '{{intentName}}'
  static intentType: any = {{intentType}}

  {{actionExample}}
}

