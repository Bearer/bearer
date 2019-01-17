import { FetchData } from '@bearer/intents'
// Uncomment this line if you need to use Client
// import Client from './client'

export default class Intent {
  static intentType = FetchData

  static async action(event: any): Promise<any> {
    // const token = event.context.authAccess.accessToken
    // Put your logic here
    return { data: [] }
  }
}
