import {
  Component,
  Method,
  State,
  Prop,
  Event,
  EventEmitter,
  Listen
} from '@bearer/core'

@Component({
  tag: 'bearer-navigator-screen',
  styleUrl: 'NavigatorScreen.scss',
  shadow: true
})
export class BearerNavigatorScreen {
  @State() visible: boolean = false
  @State() data: any

  @Prop() navigationTitle: any
  @Prop() renderFunc: (data: any) => void
  @Prop() name: string

  @Event() stepCompleted: EventEmitter
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
    if (typeof this.navigationTitle === 'string') {
      return this.navigationTitle
    }
    return this.navigationTitle(this.data)
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

  render() {
    if (!this.visible) {
      return false
    }
    return (
      <div class="screen">
        {this.renderFunc ? (
          this.renderFunc({ data: this.data, next: this.next, prev: this.prev })
        ) : (
          <slot />
        )}
      </div>
    )
  }
}
