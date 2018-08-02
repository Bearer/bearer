import { SaveStateIntent, BearerFetch } from '@bearer/core'
export class WithDeprecatedSaveStateIntent {
  @SaveStateIntent() fetcher: BearerFetch
}
