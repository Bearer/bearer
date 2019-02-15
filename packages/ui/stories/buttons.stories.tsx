import * as React from 'react'
import { storiesOf } from '@storybook/react'

const kinds = [false, 'embed', 'primary', 'secondary', 'danger', 'error']

storiesOf('Buttons', module).addWithJSX(
  'colors',
  () => kinds.map(k => [<bearer-button kind={k}>{k || 'default'}</bearer-button>, ' ']),
  {
    indent_size: 2
  }
)
