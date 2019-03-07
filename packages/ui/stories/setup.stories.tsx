import * as React from 'react'
import { storiesOf } from '@storybook/react'

const fields1 = [
  { type: 'text', label: 'Client ID', controlName: 'clientID' },
  { type: 'password', label: 'Client Secret', controlName: 'clientSecret' }
]
storiesOf('Setup component', module)
  .addWithJSX('block', () => {
    return (
      <div>
        <bearer-setup integrationId="BEARER_INTEGRATION_ID" fields="oauth2" />
        <bearer-setup-display />
      </div>
    )
  })
  .addWithJSX('inline', () => {
    return (
      <div>
        <bearer-setup display="inline" integrationId="BEARER_INTEGRATION_ID" fields={fields1} />
        <bearer-setup-display />
      </div>
    )
  })
