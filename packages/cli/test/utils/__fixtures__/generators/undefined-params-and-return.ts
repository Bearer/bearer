import { FetchData, TOAUTH2AuthContext } from '@bearer/functions'

export default class FunctionUndefinedStuff extends FetchData implements FetchData<any, any, TOAUTH2AuthContext> {
  async action(event: any): Promise<any> {
    return { data: [] }
  }
}
