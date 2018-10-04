import { Prop, RootComponent } from '@bearer/core'

@RootComponent({
  group: 'complex-feature',
  role: 'display'
})
export class FeatDisplayRootComponent {
  @Prop()
  aStringProp: string = 'ok'
  @Prop()
  aNumberProp: number = 5
}
