import { Component, Method, State, Prop, Event, EventEmitter, Listen } from '@bearer/core'

@Component({
  tag: 'bearer-navigator-screen',
  styleUrl: 'NavigatorScreen.scss',
  shadow: true
})
export class BearerNavigatorScreen {
  @State() visible: boolean = false
  @State() data: any

  @Prop() navigationTitle: ((data: any) => string) | string
  @Prop()
  renderFunc: <T>(
    params: {
      next: (data: any) => void
      prev: () => void
      complete: () => void
      data: T
    }
  ) => void
  @Prop() name: string

  @Event() stepCompleted: EventEmitter
  @Event() scenarioCompleted: EventEmitter
  @Event() navigatorGoBack: EventEmitter

  @Method()
  willAppear(data) {
    this.data = data
    this.visible = true
  }

  @Method()
  willDisappear() {
    this.visible = false
  }

  @Method()
  getTitle() {
    if (typeof this.navigationTitle === 'function') {
      return this.navigationTitle(this.data)
    }
    return this.navigationTitle
  }

  @Listen('completeScreen')
  completeScreenHandler({ detail }) {
    this.next(detail)
  }

  next = data => {
    const payload = this.name ? { [this.name]: data } : data
    this.stepCompleted.emit(payload)
  }

  prev = () => {
    this.navigatorGoBack.emit()
  }

  complete = () => {
    this.scenarioCompleted.emit()
  }

  render() {
    if (!this.visible) {
      return false
    }
    return (
      <div class="screen">
        {this.renderFunc ? (
          this.renderFunc({
            data: this.data,
            next: this.next,
            prev: this.prev,
            complete: this.complete
          })
        ) : (
          <slot />
        )}
      </div>
    )
  }
}
