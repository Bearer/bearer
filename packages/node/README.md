# `@bearer/node`

> TODO: description

## Usage

### Call a Bearer function

```tsx
// somewhere in your application, we'll use an express route here
import clientFactory from '@bearer/node/lib/client'

const bearerClient = clientFactory(process.env.BEARER_SECRET_TOKEN)
// You can pass query or body parameter depending on Function requirement
const options = { query: { status: 'open' }, body: { title: 'title' } }

bearerClient
  .invoke('1234-integration-to-invoke', 'functionName', options)
  .then(() => {
    console.log('Successfully invokeed function')
  })
  .catch(() => {
    console.log('Something wrong happened')
  })

// or async/await
try {
  const response = await bearerClient.invoke('1234-integration-to-invoke', 'functionName', options)

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

const integrationClient = new IntegrationClient(process.env.BEARER_SECRET_TOKEN, 'a-integration-uuid')

const reponse = await integrationClient.invoke('functionName', options)
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

### Use Bearer express webhook middleware

```tsx
// your server.ts
import express from 'express'
import bearerWebhooks from '@bearer/node/lib/express'

const app = express()

// each valueS must be a fonction returning a promise
const webhookHandlers = {
  ['integration-name-to_handle']: req =>
    new Promise(() => {
      // you logic goes here
      if (something) resolve()
      else {
        reject()
      }
    }),
  ['with-async-await']: async req => {
    // you logic goes here
    // ex: console.log(req.body)
    const reponse = await somethingYouWantToWaitFor
    if (response.success) {
      return whatever
    } else {
      throw new Error('An error occured')
    }
  }
}
// Without options
app.use('/whaterver_path_you_want/webhhoks', bearerWebhooks(webhookHandlers))

// With options
app.use('/whaterver_path_you_want/webhhoks', bearerWebhooks(webhookHandlers), { token: 'YOU_SECRET_TOKEN' })
```
