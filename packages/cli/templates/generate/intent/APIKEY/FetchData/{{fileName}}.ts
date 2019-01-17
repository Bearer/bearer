import { {{intentType}}, T{{authType}}AuthContext, {{callbackType}} } from '@bearer/intents'
// Uncomment this line if you need to use Client
// import Client from './client'

export default class {{intentClassName}}Intent {
  static intentType: any = {{intentType}}

  {{actionExample}}
}
