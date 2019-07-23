# @bearer/js

[![Version](https://img.shields.io/npm/v/@bearer/js.svg)](https://npmjs.org/package/@bearer/js)
![npm bundle size (scoped)](https://img.shields.io/bundlephobia/minzip/@bearer/js.svg)
[![Downloads/week](https://img.shields.io/npm/dw/@bearer/js.svg)](https://npmjs.org/package/@bearer/js)
[![License](https://img.shields.io/npm/l/@bearer/js.svg)](https://github.com/Bearer/bearer/packages/cli/blob/master/package.json)

The hassle-free way to use bearer's integrations into any web application

## Getting started

Bearer lib can be used instantly in your page or with a package system.

### Directly in your page

```html
<script src="https://cdn.jsdelivr.net/npm/@bearer/js@latest/lib/bearer.production.min.js"></script>
<script>
  // you have now access to a global `bearer` function, initialize your code by passing the `clientId` as parameter
  const bearerClient = bearer('clientId')
</script>
```

### With a build system

```bash
yarn add @bearer/js
# or
npm install @bearer/js
```

In your app

```jsx
import bearer from '@bearer/js'

class MyApp {
  componentDidMount() {
    bearer('clientId')
  }
}
```

## Usage

## Calling any APIs

Out of the box, Bearer provides an added value proxy for any API. You can call any API (aka integration) endpoint as follow

```js
const bearerClient = bearer('BEARER_CLIENT_ID')

// initialize an API client to target a specific service (ex: slack, github, etc...)
const slack = bearerClient.integration('INTEGRATION_ID')

// re-use previously created API client
slack
  .auth(authId) // see #connect section for more information
  .get('/reminders.list') // all REST verbs are available here: .post, .put, .delete etc...
  .then(console.log)
  .catch(console.error)

// passing extra arguments
slack
  .auth(authId)
  .post('/reminders.add', { text: 'Remind me something', time: 'in 10 seconds'})
  .then(console.log)
  .catch(console.error)
```


### Invoke js

The Bearer SDK for JavaScript lets you invoke integration's function (if their execution is not restricted to server side usage).

```js
const bearerClient = bearer('BEARER_CLIENT_ID')

const myIntegration = bearerClient.integration('INTEGRATION_ID')
myIntegration
  .invoke('myFunction')
  .then(console.log)
  .catch(console.error)

// is equivalent to
bearerClient
  .invoke('INTEGRATION_ID', 'myFunction')
  .then(console.log)
  .catch(console.error)
```

Passing params to your function works as follow:

```js
bearerClient
  .invoke('INTEGRATION_ID', 'myFunction', {
    query: { foo: 'bar' }
  })
  .then(console.log)
  .catch(console.error)
```

### i18n

`@bearer/js` comes with an i18n module that let you deal with internationalization of Bearer's integrations

**bearer.i18n.locale**

Lets you change the locale

```js
bearer.i18n.locale = 'es'
```

**bearer.i18n.load**

Lets you load custom translation for integrations

```js
// with a simple dictionnary
const dictionnary = { titles: { welcome: 'Ola!' } }
bearer.i18n.load('integration-uuid', dictionnary)

// with a promise returning a dictionnary
const promiseReturningADictionnary = Promise.new((resolve, reject) => {
  // async stuff
  resolve({ titles: { welcome: 'Ola!' } })
})
bearer.i18n.load('integration-uuid', promiseReturningADictionnary)

// for a given locale
const dictionnary = { titles: { welcome: 'Guten Morgen' } }
bearer.i18n.load('integration-uuid', dictionnary, { locale: 'de' })

// for multiple integrations on a single page
const dictionnary = {
  ['integration-one-uuid']: { title: { welcome: 'Hello my friend' } },
  ['integration-two-uuid']: { message: { goodbye: 'Bye Bye' } }
}
bearer.i18n.load(null, dictionnary)
```

### Secure

If you want to add a level of security, you can switch to the secure mode:

```js
window.bearer.secured = true
// at the initialisation time
window.bearer('clientId', { secured: true })
```

Once this mode is turned on, all your values passed in the properties need to be encrypted using your `ENCRYPTION_KEY`.

#### CLI

Within the CLI, you can use `bearer encrypt` to get the

```js
bearer encrypt ENCRYPTION_KEY MESSAGE
```

#### NodeJS

```typescript
import Cipher from '@bearer/security'

const cipher = new Cipher(ENCRYPTION_KEY)
cipher.encrypt(MESSAGE)
```

### connect

`connect` lets you easily retrieve the `auth-id` for an integration using OAuth authentication. Before using it, you'll need to generate a `setup-id` with the setup component of your integration

```js
bearerClient
  .connect('integration-uuid', 'setup-id')
  .then(data => {
    // user has connected himself to the OAuth provider and you now have access to the authId
    console.log(data.authId)
  })
  .catch(() => {
    // user canceled the authentication
  })
```

you can pass your own `auth-id` within options parameters as follows

```js
bearerClient.connect('integration-uuid', 'setup-id', { authId: 'my-own-non-guessable-auth-id' })
```

### init options

_Soon_
