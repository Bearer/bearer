import { FetchData, TOAUTH2AuthContext } from '@bearer/intents'

export default class IntentUndefinedStuff extends FetchData implements FetchData<any, any, TOAUTH2AuthContext> {
  async action(event: any): Promise<any> {
    return { data: [] }
  }
}
