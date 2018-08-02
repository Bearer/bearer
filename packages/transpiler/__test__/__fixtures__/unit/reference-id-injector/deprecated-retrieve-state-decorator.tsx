import { RetrieveStateIntent, BearerFetch } from '@bearer/core'
export class WithDeprecatedRetrieveStateIntent {
  @RetrieveStateIntent() fetcher: BearerFetch
}
