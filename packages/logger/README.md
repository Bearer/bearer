# @bearer/logger

## Usage

We recommend to create a logger file within you project

```ts
// my-logger.ts
import debug from '@bearer/logger'

export default () => debug('my-package-name')
```

then in your app

```ts
import debug from 'path/to/my-logger'

const logger = debug()
logger('message to debug')
logger({ object: 'something' })

// sub logger

const subLogger = logger.extend('sub-feature')
subLogger('blablabl')
```

You'll need to set `DEBUG=*` to see all logs or `DEBUG=bearer:my-package-name` to see logs produced by your application.

### Browser support

TODO
