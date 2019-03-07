import { Function, FunctionType, BearerFetch } from '@bearer/core'

export class FunctionDecorated {
  @Function('saveState', FunctionType.SaveState)
  fetcher: BearerFetch
}
