import * as React from 'react'
import { storiesOf } from '@storybook/react'

storiesOf('I18n', module).addWithJSX('colors', () => [
  <bearer-i18n key="missing.key" var='{"ok": "asdasdasddas"}' default="defautl value" />
])
