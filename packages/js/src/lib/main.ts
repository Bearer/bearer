import Bearer, { TBearerOptions } from './bearer'
import debug from './logger'

import i18n, { I18n } from './i18n'

/**
 * @param  {string} clientId Client ID you'll find within the developer portal > Settings
 * @param  {Partial<TBearerOptions>} options? Fine tune bearer functionalities
 */
// @ts-ignore
const bearer: TBearer = (clientId: string, options?: Partial<TBearerOptions>) => {
  bearer.instance = new Bearer(clientId, options)
  return bearer.instance
}

export type TBearer = {
  (clientId: string, options?: Partial<TBearerOptions>): Bearer
  instance: Bearer
  version: string
  i18n: I18n
  secured: boolean
  _instance?: Bearer
}

export default bearer

// Handle non instantiated bearer client
Object.defineProperty(bearer, 'instance', {
  get: function(this: TBearer) {
    if (!this._instance) {
      logMissingInstance()
      this._instance = new Bearer(undefined)
    }
    return this._instance
  },
  set: function(this: TBearer, bearerInstance: Bearer) {
    this._instance = bearerInstance
  }
})

bearer.version = 'BEARER_VERSION'
bearer.i18n = i18n

// fake the presence of secured
bearer.secured = false
Object.defineProperty(bearer, 'secured', {
  get: function() {
    return this.instance && this.instance.secured
  },
  set: function(secured: boolean) {
    this.instance.secured = secured
  }
})

const MISSING_INSTANCE_MEESSAGE =
  '%c No bearer client has been initialized. please make sure you call bearer("YOUR_CLIENT") before using any component of backend function'
const MISSING_INSTANCE_STYLE = 'font-weight:bold;'

function logMissingInstance() {
  if (debug.enabled) {
    debug.extend('main')(MISSING_INSTANCE_MEESSAGE, MISSING_INSTANCE_STYLE)
  } else {
    console.warn(MISSING_INSTANCE_MEESSAGE, MISSING_INSTANCE_STYLE)
  }
}
