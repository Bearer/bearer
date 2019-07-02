# `@bearer/node`

[![Version](https://img.shields.io/npm/v/@bearer/logger.svg)](https://npmjs.org/package/@bearer/logger)
![npm bundle size (scoped)](https://img.shields.io/bundlephobia/minzip/@bearer/logger.svg)
![node (scoped)](https://img.shields.io/node/v/@bearer/node.svg)
[![Downloads/week](https://img.shields.io/npm/dw/@bearer/logger.svg)](https://npmjs.org/package/@bearer/logger)
[![License](https://img.shields.io/npm/l/@bearer/logger.svg)](https://github.com/Bearer/bearer/packages/logger/blob/master/package.json)

Node client to interact with [Bearer.sh](https://www.bearer.sh)

## Usage

Get your Bearer's [credentials](https://app.bearer.sh/keys) and setup Bearer as follow:

### Calling an API using Bearer

```tsx
import Bearer from '@bearer/node'

const bearerClient = Bearer(process.env.BEARER_API_KEY) // copy and paste the `API key`
const github = bearerClient.integration('INTEGRATION_ID') // you'll find it on the Bearer's dashboard

github
  .get('/users/repos')
  .then(console.log)
  .catch(console.error)
```

### Calling a custom backend function

You'll need to push your function to Bearer first ([discover how](https://docs.bearer.sh/working-with-bearer/manipulating-apis)):

```tsx
import Bearer from '@bearer/node'

const bearerClient = Bearer(process.env.BEARER_API_KEY) // copy and paste the `API key`
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
```

_Note 1_: we are using axios as the http client. Each `.get()`, `.post()`, `.put()`, ... or `.invoke()` returns an Axios Promise. https://github.com/axios/axios

_Note 2_: If you are using ExpressJS, have a look at the [@bearer/express](https://github.com/Bearer/bearer/tree/master/packages/express) client
