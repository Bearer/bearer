import { Intent, IntentType, BearerFetch } from '@bearer/core'

export class IntentDecorated {
  @Intent('saveState', IntentType.RetrieveState)
  fetcher: BearerFetch
}
