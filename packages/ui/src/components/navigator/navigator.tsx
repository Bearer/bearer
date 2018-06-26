import { Component, State, Listen, Element } from '@bearer/core'

@Component({
  tag: 'bearer-navigator',
  styleUrl: 'navigator.scss',
  shadow: true
})
export class BearerNavigator {
  @Element() el: HTMLElement
  @State() screens: Array<any> = []
  @State() screenData: Object = {}
  @State() _visibleScreen: number
  @State() navigationTitle: string

  set visibleScreen(index) {
    if (this._visibleScreen >= 0) {
      const currentScreen = this.screens[this._visibleScreen]
      if (currentScreen) {
        currentScreen.willDisappear()
        currentScreen.classList.remove('in')
      }
    }
    this._visibleScreen = index
    const newScreen = this.screens[this._visibleScreen]
    if (newScreen) {
      newScreen.willAppear(this.screenData)
      this.navigationTitle = newScreen.getTitle()
      newScreen.classList.add('in')
    }
  }

  get visibleScreen(): number {
    return this._visibleScreen
  }

  @Listen('scenarioCompleted')
  scenarioCompletedHandler() {
    this.screenData = {}
    this.visibleScreen = 0
  }

  @Listen('stepCompleted')
  stepCompletedHandler(event) {
    event.preventDefault()
    event.stopImmediatePropagation()
    this.screenData = {
      ...this.screenData,
      ...event.detail
    }
    this.next(null)
  }

  @Listen('navigatorGoBack')
  prev(e) {
    if (e) {
      e.preventDefault()
      e.stopPropagation()
    }
    if (this.hasPrevious()) {
      this.visibleScreen = this.visibleScreen - 1
    }
  }

  next = e => {
    if (e) {
      e.preventDefault()
      e.stopPropagation()
    }
    if (this.hasNext()) {
      this.visibleScreen = this.visibleScreen + 1
    }
  }

  hasNext = () => this.visibleScreen < this.screens.length - 1
  hasPrevious = () => this.visibleScreen > 0

  componentDidLoad() {
    if (this.el.shadowRoot) {
      this.screens = this.el.shadowRoot
        .querySelector('slot:not([name])')
        ['assignedNodes']()
        .filter(node => node.willAppear)
    }
    this.visibleScreen = 0
  }

  render() {
    return (
      <div>
        <h3 class="title">
          <bearer-navigator-back disabled={!this.hasPrevious()} />
          <slot name="header-name">{this.navigationTitle}</slot>
        </h3>
        <slot />
      </div>
    )
  }
}
