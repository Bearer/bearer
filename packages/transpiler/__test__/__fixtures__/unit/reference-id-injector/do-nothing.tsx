import { Intent, BearerFetch } from '@bearer/core'
class DoNothing {
  @OtherDecorator other: any
  @Intent('MyIntent') fetcher: BearerFetch
}
