# @bearer/react

React tooling for bearer.sh components

## Installation

Install dependencies

```bash
yarn add @bearer/js @bearer/react
```

## Factory

```tsx
import { factory, Bearer } from '@bearer/react'
const { Connect } = factory('my-integration-name')
```

### `Connect` component

This component expect a `setupId` to authenticate the final user.

```tsx
import { factory } from '@bearer/react'
const { Connect } = factory('my-integration-name')

class MyComponent extends React.Component {
  handleSuccess = ({ authId }) => {
    // do whatever you want with the authId received
  }
  render() {
    return (
      <Bearer clientId="clientId">
        <Connect
          setupId="setupId"
          render={({ loading, connect, error }) => {
            if (loading) {
              return <Loading />
            }
            if (error) {
              return (
                <div>
                  Error, please retry <button onClick={connect}>Retry</button>
                </div>
              )
            }
            return <button onClick={connect}>Connect ...</button>
          }}
        />
      </Bearer>
    )
  }
}
```

## fromBearer

This component takes html tag name of a component as well as optional type information to produce a react component that can we re-used throughout an application.

```tsx
import * as React from 'react'
const MyText = fromBearer<{ text?: string; anotherPropery: string }>('bearer-my-text')

class ReactComponent extends React.Component {
  render() {
    return <MyText text="hello" anotherPropery="world" />
  }
}
```

Output events can be handled by simply adding the eventname as a prop on the component and adding a basic handler

```tsx
<Share bearer-uuid-feature-shared={this.onShared} />
```

For these tags to work correctly they must have a parent `Bearer`

## Bearer

This component maintains a shared state for a group of components as well as adding the bearer script tags to the page. This tag can be added at any level of the application as is convenient for the implementation

### Example Use

we will use [this integration](https://app.bearer.sh/integrations/6d29c4-share-slack-beta-4/preview) for our examples:

First on the preview page obtain a setup-id from the setup components so we do not need to include these in our page. Next lets define our components in a constants file:

`bearer.ts`

```TS
import { fromBearer } from '@bearer/react'

export const ChannelSelect = fromBearer('bearer-6d29c4-share-slack-beta-4-channel-action')
export const Share = fromBearer<{message?:string, text?:string}>('bearer-6d29c4-share-slack-beta-4-feature-action')
export const SlackConnect = fromBearer('bearer-6d29c4-share-slack-beta-4-connect-action')
```

We can include the `Bearer` at any level but in this example lets use it all together in the same component

`slack-share-component.tsx`

```tsx
import * as React from 'react'
import { Bearer } from '@bearer/react'
import { SlackConnect, ChannelSelect, Share } from '../../constants/bearer'

export default class SlackShareSetup extends React.Component {
  public render() {
    const intialContext = { 'setup-id': 'SETUP_ID_SAMPLE' }
    return (
      <Bearer clientId="YOUR_CLIENT_ID" initialContext={intialContext}>
        <SlackConnect />
        <ChannelSelect />
        <Share message="hello world!" text="Test!" bearer-6d29c4-share-slack-beta-4-feature-shared={this.onShared} />
      </Bearer>
    )
  }
}
```

This allows to share messages but what if we want to persist the users setup information. @bearer/react provides two methods to acomplish this.

### Using onUpdate

`Bearer` has an `onUpdate`callback we can hook into which is called every time data changes within the provider

`slack-share-component.tsx`

```tsx
import * as React from 'react'
import { Bearer } from '@bearer/react'
import { SlackConnect, ChannelSelect, Share } from '../../constants/bearer'

export default class SlackShareSetup extends React.Component {
  public render() {
    const intialContext = { 'setup-id': 'SETUP_ID_SAMPLE' }
    return (
      <Bearer
        clientId="YOUR_CLIENT_ID"
        initialContext={intialContext}
        onUpdate={(data: any) => {
          this.setState({ data })
        }}
      >
        <SlackConnect />
        <ChannelSelect />
        <Share message="hello world!" text="Test!" bearer-6d29c4-share-slack-beta-4-feature-shared={this.onShared} />
        <button onClick={this.handleSave}>Save Setup</button>
      </Bearer>
    )
  }
  private handleSave = () => console.log('handleSave', this.state.data)
}
```

### using a context consumer

Internally `Bearer` uses the [react context API](https://reactjs.org/docs/context.html). For a more advaned but flexable method we can access the consumer directly via `BearerContext` and then use the currently set details as we wish.

`slack-share-component.tsx`

```tsx
import * as React from 'react'
import { Bearer, BearerContext } from '@bearer/react'
import { SlackConnect, ChannelSelect, Share } from '../../constants/bearer'

export default class SlackShareSetup extends React.Component {
  public render() {
    const intialContext = { 'setup-id': 'SETUP_ID_SAMPLE' }
    return (
      <Bearer clientId="YOUR_CLIENT_ID" initialContext={intialContext}>
        <SlackConnect />
        <ChannelSelect />
        <Share message="hello world!" text="Test!" bearer-6d29c4-share-slack-beta-4-feature-shared={this.onShared} />
        <BearerContext.Consumer>
          {context => <button onClick={this.handleSave(context.state)}>Save Setup</button>}
        </BearerContext.Consumer>
      </Bearer>
    )
  }
  private handleSave = (data: any) => () => {
    console.log('handleSave', data)
  }
}
```
