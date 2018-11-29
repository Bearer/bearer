import { RetrieveState } from '@bearer/intents'

type AType = { aliasedParams: string }

// @ts-ignore
type TOAUTH2AuthContext = {}
// @ts-ignore
type TFetchDataCallback = (p: any) => void

export class TypeAliasClass {
  // @ts-ignore
  static intentType = RetrieveState
  static intentName = 'TypeAliasParams'

  static action(context: TOAUTH2AuthContext, params: AType, body: any, callback: TFetchDataCallback) {
    callback({ data: [] })
  }
}
