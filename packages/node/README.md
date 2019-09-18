# `@bearer/node`

[![Version](https://img.shields.io/npm/v/@bearer/logger.svg)](https://npmjs.org/package/@bearer/logger)
![npm bundle size (scoped)](https://img.shields.io/bundlephobia/minzip/@bearer/logger.svg)
![node (scoped)](https://img.shields.io/node/v/@bearer/node.svg)
[![Downloads/week](https://img.shields.io/npm/dw/@bearer/logger.svg)](https://npmjs.org/package/@bearer/logger)
[![License](https://img.shields.io/npm/l/@bearer/logger.svg)](https://github.com/Bearer/bearer/packages/logger/blob/master/package.json)

Node client to query any APIs and custom functions using [Bearer.sh](https://www.bearer.sh)

## Usage

Get your Bearer's [credentials](https://app.bearer.sh/keys) and setup Bearer as follow:

### Calling any APIs

```tsx
import bearer from '@bearer/node'

const client = bearer(process.env.BEARER_SECRET_KEY) // find it on https://app.bearer.sh/keys
const github = client.integration('INTEGRATION_ID') // you'll find it on the Bearer's dashboard

github
  .get('/repositories')
  .then(console.log)
  .catch(console.error)
```

More advanced examples:

```tsx
// With query parameters
github
  .get('/repositories', { query: { since: 364 } })
  .then(console.log)
  .catch(console.error)

// Making an authenticated POST
github
  .auth(authId) // Create an authId for GitHub on https://app.bearer.sh
  .post('/user/repos', { body: { name: 'Just setting up my Bearer.sh' } })
  .then(console.log)
  .catch(console.error)
```

Using `async/await`:

```tsx
const response = await github
  .auth(authId) // Create an authId for GitHub on https://app.bearer.sh
  .post('/user/repos', { body: { name: 'Just setting up my Bearer.sh' } })

console.log(response)
```

### Calling custom functions

```tsx
import bearer from '@bearer/node'

const client = bearer(process.env.BEARER_SECRET_KEY) // copy and paste the `API key`
const github = client.integration('INTEGRATION_ID')

github
  .invoke('myFunction')
  .then(console.log)
  .catch(console.error)
```

[Learn more](https://docs.bearer.sh/working-with-bearer/manipulating-apis) on how to use custom functions with Bearer.sh.

## Notes

_Note 1_: we are using [axios](https://github.com/axios/axios) as the http client. Each `.get()`, `.post()`, `.put()`, ... or `.invoke()` returns an Axios Promise.

_Note 2_: If you are using ExpressJS, have a look at the [@bearer/express](https://github.com/Bearer/bearer/tree/master/packages/express) client
