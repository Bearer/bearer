import { BearerRef, Input, RootComponent, State, Prop } from '@bearer/core'

@RootComponent({
  group: 'no-options',
  role: 'action'
})
class NoOptionsComponent {
  @Input()
  aString: BearerRef<string> = 'ok'

  @Input()
  aStringWithoutInitializer: BearerRef<string>

  @Input()
  object: BearerRef<{ title: string }> = { title: 'Guest' }

  @Prop({ mutable: true })
  aValueThatChanges: String

  @Prop()
  aStaticValue: String
}
