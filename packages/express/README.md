# `@bearer/express`


[![Version](https://img.shields.io/npm/v/@bearer/express.svg)](https://npmjs.org/package/@bearer/express)
![npm bundle size (scoped)](https://img.shields.io/bundlephobia/minzip/@bearer/express.svg)
[![Downloads/week](https://img.shields.io/npm/dw/@bearer/express.svg)](https://npmjs.org/package/@bearer/express)
[![License](https://img.shields.io/npm/l/@bearer/express.svg)](https://github.com/Bearer/bearer/packages/cli/blob/master/package.json)

> TODO: description

## Usage

Get your Bearer's [credentials](https://app.bearer.sh/keys) and setup Bearer as follow:

```tsx
// your server.ts
import express from 'express'
import bearerWebhooks from '@bearer/express'

const app = express()

// each value must be a function returning a promise
const webhookHandlers = {
  ['INTEGRATION_UUID']: req =>
    new Promise(() => {
      // you logic goes here
      if (something) resolve()
      else {
        reject()
      }
    }),
  // With ASync
  ['INTEGRATION_UUID']: async req => {
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

app.use('/webhooks', bearerWebhooks(webhookHandlers, { token: 'ENCRYPTION_KEY' })) // Copy and Paste you Encryption Key
```

