import { BearerState, Prop } from '@bearer/core'

class HasBearerContextInjected {
  @BearerState attachPullRequests: any[] = []
  @Prop({ context: 'bearer' })
  bearerContext: any
}
