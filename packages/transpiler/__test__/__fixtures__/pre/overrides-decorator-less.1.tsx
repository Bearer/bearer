import { Component } from '@bearer/core'

@Component({
  tag: 'overrides-decorator-less'
})
class OverridesDecorarotLess {
  constructor() {}

  componentDidLoad() {
    console.log('componentDidLoad')
  }

  componentWillLoad() {
    console.log('componentWillLoad')
  }

  componentDidUnload() {
    console.log('componentDidUnload')
  }
}
