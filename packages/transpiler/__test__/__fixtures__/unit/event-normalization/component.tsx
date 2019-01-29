import { Event, Listen } from '@bearer/core'

class Class {
  @Event() event: any
  @Event({ eventName: 'body:SomeThinIDontwant_to_beeeeeeenoramlized123-', somethingElseuntouched: 'Something' })
  bodyscopeNotNormalized: any
  @Event({ eventName: 'SomeThinIDontwant_to_beeeeeeenoramlized123-', somethingElseuntouched: 'Something' })
  totNormalized: any
  @Event({ eventName: 'body:bearer-SomeThinIWant_to_beeeeeeenoramlized123-', somethingElseuntouched: 'Something' })
  bodyscopedNormalized: any
  @Event({ eventName: 'bearer-SomeThinIWant_to_beeeeeeenoramlized123-', somethingElseuntouched: 'Something' })
  normalized: any

  @Listen('anUntouchedEvent') anEvent: any
  @Listen('body:bearer-somethingSdsfjkh_-') bodyScopednormalizedEvent: any
  @Listen('bearer-somethingNeedsTobe-normalized') normalizedEvent: any
}
