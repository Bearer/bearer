# `@bearer/node`

> TODO: description

## Usage

### Call a Bearer intent

```tsx
// somewhere in your application, we'll use an express route here
import clientFactory from '@bearer/node/lib/client'

const bearerClient = clientFactory(process.env.BEARER_SECRET_TOKEN)
// You can pass query or body parameter depending on Intent requirement
const options = { query: { status: 'open' }, body: { title: 'title' } }

bearerClient
  .call('1234-scenario-to-call', 'intentName', options)
  .then(() => {
    console.log('Successfully called intent')
  })
  .catch(() => {
    console.log('Something wrong happened')
  })

//async await wait
try {
  const reponse = await bearerClient.call('1234-scenario-to-call', 'intentName', options)
} catch (e) {
  // handler error
}
// play with response here
```

_Note_: we are using axios a http client. Each .call() returns an Axios Promise. https://github.com/axios/axios

### Scenario client

Scenario client facilitates intent calls and prevent you to pass scenario name on every call

```tsx
import { ScenarioClient } from '@bearer/node/lib/client'

const scenarioClient = new ScenarioClient(process.env.BEARER_SECRET_TOKEN, 'a-scenario-uuid')

const reponse = await scenarioClient.call('intentName', options)
```

If you are a Typescript user, you can provide a list of intents to use for a scenario:

```tsx
const scenarioClient = new ScenarioClient<'intentName' | 'other-intent'>(
  process.env.BEARER_SECRET_TOKEN,
  'a-scenario-uuid'
)

scenarioClient.call('intentName', options) // OK
scenarioClient.call('other-intent', options) // OK
scenarioClient.call('unknow-intent', options) // Argument of type '"unknow-intent"' is not assignable to parameter of type 'TScenarioIntentNames'.
```

### Use Bearer express webhook middleware

```tsx
// your server.ts
import express from 'express'
import bearerWebhooks from '@bearer/node/lib/express'

const app = express()

// each valueS must be a fonction returning a promise
const webhookHandlers = {
  ['scenario-name-to_handle']: req =>
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
