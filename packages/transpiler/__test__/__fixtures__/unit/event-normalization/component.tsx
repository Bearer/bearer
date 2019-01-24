import { Event } from '@bearer/core'

class Class {
  @Event() event: any
  @Event({ eventName: 'SomeThinIwant_to_noramlize123-', somethingElseuntouched: 'Something' }) normalized: any
}
