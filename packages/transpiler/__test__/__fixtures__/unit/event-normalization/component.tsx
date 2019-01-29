import { Event, Listen } from '@bearer/core'

class Class {
  @Event() event: any
  @Event({ eventName: 'SomeThinIwant_to_noramlize123-', somethingElseuntouched: 'Something' }) normalized: any

  @Listen('anUntouchedEvent') anEvent: any
  @Listen('body:bearer-somethingSdsfjkh_-') bodyScopednormalizedEvent: any
  @Listen('bearer-somethingNeedsTobe-normalized') normalizedEvent: any
}
