import { Component } from '@bearer/core'
import Bearer from '@bearer/core'

@Component({
  tag: 'bearer-loading',
  styleUrl: 'Loader.scss',
  shadow: true
})
export class Button {
  render() {
    const { loadingComponent } = Bearer.config
    if (loadingComponent) {
      const Tag = loadingComponent
      return <Tag />
    }
    return (
      <div id="root">
        <div id="loader">
          <div id="d1" />
          <div id="d2" />
          <div id="d3" />
          <div id="d4" />
          <div id="d5" />
        </div>
      </div>
    )
  }
}
