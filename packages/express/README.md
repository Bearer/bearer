# `@bearer/express`

> TODO: description

## Usage

```tsx
// your server.ts
import express from 'express'
import bearerWebhooks from '@bearer/express'

const app = express()

// each value must be a function returning a promise
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
