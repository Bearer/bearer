import { BearerRef, Input, RootComponent, State, Prop } from '@bearer/core'

@RootComponent({
  name: 'complete-options'
})
class NoOptionsComponent {
  @Input()
  aStringInput: BearerRef<string> = 'ok'

  @Prop({ mutable: true })
  aValueThatChanges: String

  @Prop()
  aStaticValue: String
}
