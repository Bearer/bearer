# `@bearer/node`

> TODO: description

## Usage

### Call a Bearer function

```tsx
import bearer from '@bearer/node'

const bearerClient = bearer(process.env.BEARER_API_KEY)
// You can pass query or body parameter depending on Function requirement
const options = { query: { status: 'open' }, someData: 'anything' }

bearerClient
  .invoke('INTEGRATION_ID', 'myFunction', options)
  .then(() => {
    console.log('Successfully invoked function')
  })
  .catch(() => {
    console.log('Something went wrong')
  })

// or async/await
try {
  const response = await bearerClient.invoke('INTEGRATION_ID', 'myFunction', options)

  // play with response here
} catch (e) {
  // handle error
}
```

_Note_: we are using axios a http client. Each .invoke() returns an Axios Promise. https://github.com/axios/axios

### Integration client

Integration client facilitates func invokes and prevent you to pass integration name on every invoke

```tsx
import { IntegrationClient } from '@bearer/node/lib/client'

const integrationClient = new IntegrationClient(process.env.BEARER_API_KEY, 'a-integration-uuid')

const reponse = await integrationClient.invoke('myFunction', options)
```

If you are a Typescript user, you can provide a list of functions to use for a integration:

```tsx
const integrationClient = new IntegrationClient<'functionName' | 'other-function'>(
  process.env.BEARER_SECRET_TOKEN,
  'a-integration-uuid'
)

integrationClient.invoke('functionName', options) // OK
integrationClient.invoke('other-function', options) // OK
integrationClient.invoke('unknow-function', options) // Argument of type '"unknow-function"' is not assignable to parameter of type 'TIntegrationFunctionNames'.
```

