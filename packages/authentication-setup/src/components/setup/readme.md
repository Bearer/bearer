# bearer-save-setup

<!-- Auto Generated Below -->

## Properties

| Property      | Attribute     | Description                            | Type                                                       | Default     |
| ------------- | ------------- | -------------------------------------- | ---------------------------------------------------------- | ----------- |
| `clientId`    | `client-id`   | Your Bearer clientId                   | `string`                                                   | `undefined` |
| `integration` | `integration` | Integration identifier                 | `string`                                                   | `undefined` |
| `noError`     | `no-error`    |                                        | `boolean`                                                  | `undefined` |
| `setupId`     | `setup-id`    | Optionally provide your custom setupId | `string`                                                   | `undefined` |
| `type`        | `type`        | Authentication Type of the integration | `"APIKEY" \| "BASIC" \| "OAUTH1" \| "OAUTH2" \| "unknown"` | `'unknown'` |

## Events

| Event   | Description | Type                               |
| ------- | ----------- | ---------------------------------- |
| `saved` |             | `CustomEvent<{ setupId: string }>` |

---

_Built with [StencilJS](https://stenciljs.com/)_
