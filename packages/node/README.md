# `@bearer/node`

This is the Node client to interact with Bearer's integrations.

## Usage

Get your Bearer's [credentials](https://app.bearer.sh/keys) and setup Bearer as follow:

### Call a Bearer function

```tsx
import bearer from '@bearer/node'

const bearerClient = bearer(process.env.BEARER_API_KEY) // copy and paste the `API key`
// You can pass query or body parameter depending on Function requirement
const options = { query: { status: 'open' }, someData: 'anything' }

bearerClient
  .invoke('INTEGRATION_UUID', 'myFunction', options)
  .then(() => {
    console.log('Successfully invoked function')
  })
  .catch(() => {
    console.log('Something went wrong')
  })

// or async/await
try {
  const response = await bearerClient.invoke('INTEGRATION_UUID', 'myFunction', options)

  // play with response here
} catch (e) {
  // handle error
}
```

_Note_: we are using axios a http client. Each `.invoke()` returns an Axios Promise. https://github.com/axios/axios

### Integration Client

Integration Client allows you to pass the `INTEGRATION_UUID` only once:

```tsx
import { IntegrationClient } from '@bearer/node/lib/client'

const integrationClient = new IntegrationClient(process.env.BEARER_API_KEY, 'INTEGRATION_UUID') // copy and paste the `API key`

const reponse = await integrationClient.invoke('myFunction', options)
```

If you are a TypeScript user, you can provide a list of functions to use for an integration:

```tsx
const integrationClient = new IntegrationClient<'functionName' | 'other-function'>(
  process.env.BEARER_API_KEY, // copy and paste the `API key`
  'INTEGRATION_UUID'
)

integrationClient.invoke('functionName', options) // OK
integrationClient.invoke('other-function', options) // OK
integrationClient.invoke('unknow-function', options) // Argument of type '"unknow-function"' is not assignable to parameter of type 'TIntegrationFunctionNames'.
```

_NB: If you are using ExpressJS, have a look at the [ExpressJS](https://github.com/Bearer/bearer/tree/master/packages/express) client_
