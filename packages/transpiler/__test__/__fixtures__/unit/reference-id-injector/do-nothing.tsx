import { Function, BearerFetch } from '@bearer/core'
class DoNothing {
  @OtherDecorator other: any
  @Function('MyFunction') fetcher: BearerFetch
}
