import * as React from 'react'
import { storiesOf } from '@storybook/react'

const kinds = [false, 'action', 'primary', 'secondary', 'success', 'danger', 'warning', 'info', 'light', 'dark']
const sizes = [false, 'medium', 'small', 'large']

storiesOf('Button', module)
  .addWithJSX('colors', () => kinds.map(k => [<bearer-button kind={k}>{k || 'default'}</bearer-button>, ' ']), {
    indent_size: 2
  }) // Add additional info text directly
  .addWithJSX('sizes', () => sizes.map(k => [<bearer-button size={k}>{k || 'default'}</bearer-button>, ' ']))
