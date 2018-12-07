import { RetrieveState } from '@bearer/intents'

//@ts-ignore
type TOAUTH2AuthContext = {}
//@ts-ignore
type TFetchDataCallback = (p: any) => void

export class ObjectLiteralParamsClass {
  // @ts-ignore
  static intentType = RetrieveState

  static action(
    context: TOAUTH2AuthContext,
    params: { inlineParam: string; optional?: number },
    body: any,
    callback: TFetchDataCallback
  ) {
    callback({ data: [] })
  }
}
