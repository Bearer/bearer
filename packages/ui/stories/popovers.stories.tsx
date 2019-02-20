import * as React from 'react'
import { storiesOf } from '@storybook/react'

storiesOf('Popover', module)
  .addWithJSX('direction', () => {
    return [
      <bearer-popover direction="right">
        <span slot="popover-button">Toggle me right</span>
        <div>I'm it's content.</div>
      </bearer-popover>,
      <bearer-popover direction="top">
        <span slot="popover-button">Toggle me up</span>
        <div>I'm it's content.</div>
      </bearer-popover>,
      <bearer-popover direction="bottom">
        <span slot="popover-button">Toggle me down</span>
        <div>I'm it's content.</div>
      </bearer-popover>,
      <bearer-popover direction="left">
        <span slot="popover-button">Toggle me left</span>
        <div slot="popover-header">I'm it's header.</div>
        <div>I'm it's content.</div>
      </bearer-popover>
    ]
  })
  .addWithJSX('slotted', () => {
    return [
      <bearer-popover opened="true">
        <span slot="popover-toggler">
          <button>
            I'm <code>slot-toggler</code>
          </button>
        </span>
        <div slot="popover-header">
          I'm <code>slot-header</code>.
        </div>
        <div>
          I'm the <code>default slot</code>.
        </div>
        <div slot="popover-footer">
          I'm <code>slot-footer</code>.
        </div>
      </bearer-popover>,
      <br />,
      <br />,
      <br />,
      <br />,
      <br />,
      <bearer-popover direction="bottom" aligned="right">
        <span slot="popover-button">
          I'm <code>slot-button</code>
        </span>
        <div slot="popover-header">
          I'm <code>slot-header</code>.
        </div>
        <div>
          I'm the <code>default slot</code>.
        </div>
        <div slot="popover-footer">
          I'm <code>slot-footer</code>.
        </div>
      </bearer-popover>
    ]
  })
  .addWithJSX('opened by default', () => {
    return (
      <bearer-popover opened="true">
        <span slot="popover-button">I'm a popover</span>
        <div slot="popover-container" style={{ padding: '20px' }}>
          I'm opened by default.
        </div>
      </bearer-popover>
    )
  })
  .addWithJSX('properties', () => {
    return (
      <bearer-popover opened="true">
        <span slot="popover-button">Toggle me</span>
        <div slot="popover-container" style={{ padding: '20px' }}>
          I'm opened by default.
        </div>
      </bearer-popover>
    )
  })
