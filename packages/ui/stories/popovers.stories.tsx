import * as React from 'react'
import { storiesOf } from '@storybook/react'

storiesOf('Popover', module)
  .addWithJSX('direction', () => {
    return [
      <bearer-popover direction="right">
        <span slot="popover-toggler">Toggle me right</span>
        <div slot="popover-container">I'm it's container.</div>
      </bearer-popover>,
      <bearer-popover direction="top">
        <span slot="popover-toggler">Toggle me up</span>
        <div slot="popover-container">I'm it's container.</div>
      </bearer-popover>,
      <bearer-popover direction="bottom">
        <span slot="popover-toggler">Toggle me down</span>
        <div slot="popover-container">I'm it's container.</div>
      </bearer-popover>,
      <bearer-popover direction="left">
        <span slot="popover-toggler">Toggle me left</span>
        <div slot="popover-container">I'm it's container.</div>
      </bearer-popover>
    ]
  })
  .addWithJSX('slotted', () => {
    return [
      <bearer-popover opened="true">
        <span slot="popover-toggler">
          I'm <code>slot-toggler</code>
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
      <bearer-popover opened="true" direction="bottom" aligned="right">
        <span slot="popover-toggler">
          I'm <code>slot-toggler</code>
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
        <span slot="popover-toggler">I'm a popover</span>
        <div slot="popover-container" style={{ padding: '20px' }}>
          I'm opened by default.
        </div>
      </bearer-popover>
    )
  })
