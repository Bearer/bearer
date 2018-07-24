import { Component, Element, Listen, Prop, State } from '@bearer/core'

const NAVIGATOR_AUTH_SCREEN_NAME = 'BEARER-NAVIGATOR-AUTH-SCREEN'

@Component({
  tag: 'bearer-navigator',
  shadow: true
})
export class BearerPopoverNavigator {
  @Element() el: HTMLElement
  @State() screens: Array<any> = []
  @State() screenData: Object = {}
  @State() _visibleScreen: number
  @State() navigationTitle: string

  @Prop() direction: string = 'top'
  @Prop() btnProps: JSXElements.BearerButtonAttributes = { content: 'Activate' }

  @Prop() display = 'popover'

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
    this.visibleScreen = this.hasAuthScreen() ? 1 : 0
    this.el.shadowRoot.querySelector('#button')['toggle'](false)
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

  get screenNodes() {
    return this.el.shadowRoot
      ? this.el.shadowRoot
          .querySelector('slot:not([name])')
          ['assignedNodes']()
          .filter(node => node.willAppear)
      : []
  }

  hasNext = () => this.visibleScreen < this.screens.length - 1
  hasPrevious = () => this.visibleScreen > 0
  hasAuthScreen = () => this.screenNodes.filter(node => node['tagName'] === NAVIGATOR_AUTH_SCREEN_NAME).length

  componentDidLoad() {
    this.screens = this.screenNodes
    this.visibleScreen = 0
  }

  render() {
    return (
      <bearer-button-popover
        btnProps={this.btnProps}
        id="button"
        direction={this.direction}
        header={this.navigationTitle}
        backNav={this.hasPrevious()}
      >
        <slot />
      </bearer-button-popover>
    )
  }
}
