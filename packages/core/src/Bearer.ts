import fbemitter, { EventSubscription } from 'fbemitter'
import postRobot from 'post-robot'

import BearerConfig from './BearerConfig'
import Events from './EventNames'

const BEARER_WINDOW_INSTANCE_KEY = 'BEARER_INSTANCE'
const BEARER_EMITTER = 'BEARER_EMITTER'
const BEARER_WINDOW_KEY = 'BEARER'
const BEARER_CONFIG_KEY = 'BEARER_CONFIG'
const IFRAME_NAME = 'BEARER-IFRAME'

type TAuthorizationPayload = {
  scenarioId: string
  authIdentifier?: string
}

class Bearer {
  private static set _instance(bearerInstance: Bearer) {
    if (window[BEARER_WINDOW_INSTANCE_KEY]) {
      console.warn('[BEARER]', 'Replacing bearer instance')
    }
    window[BEARER_WINDOW_INSTANCE_KEY] = bearerInstance
  }

  private static get _instance(): Bearer | undefined {
    return window[BEARER_WINDOW_INSTANCE_KEY]
  }

  public static get emitter(): fbemitter.EventEmitter {
    if (!window[BEARER_EMITTER]) {
      window[BEARER_EMITTER] = new fbemitter.EventEmitter()
    }
    return window[BEARER_EMITTER]
  }

  public static get instance(): Bearer {
    return window[BEARER_WINDOW_INSTANCE_KEY] || this.init({})
  }

  public static get config(): BearerConfig {
    return this.instance.bearerConfig
  }

  static get version() {
    return 'LIB_VERSION'
  }

  get bearerConfig(): BearerConfig {
    return window[BEARER_WINDOW_KEY]
  }

  set bearerConfig(config: BearerConfig) {
    window[BEARER_WINDOW_KEY] = config
  }

  get maybeInitialized(): Promise<boolean> {
    if (!this.isSessionInitialized) {
      console.warn('[BEARER]', 'Waiting Bearer to be initialized')
    }
    return this._maybeInitialized
  }

  set maybeInitialized(promise) {
    this._maybeInitialized = promise
  }
  public static init(config: any = {}): Bearer {
    if (this._instance) {
      console.warn('One instance is already configured, reaplacing it')
    }
    this._instance = new Bearer({ ...config, ...(window[BEARER_CONFIG_KEY] || {}) })
    return this._instance
  }

  static onAuthorized = (scenarioId: string, callback: (authorize: boolean) => void): EventSubscription => {
    console.log('[BEARER]', 'register onAuthorized', scenarioId)
    return Bearer.emitter.addListener(Events.AUTHORIZED, () => {
      // TODO : listen only for the scenarioId (+ setupId ?)
      console.log('[BEARER]', 'calling authorized listener', scenarioId)
      callback(true)
    })
  }

  static onRevoked = (scenarioId: string, callback: (authorize: boolean) => void): EventSubscription => {
    console.log('[BEARER]', 'register onRevoked', scenarioId)
    return Bearer.emitter.addListener(Events.REVOKED, () => {
      // TODO : listen only for the scenarioId (+ setupId ?)
      console.log('[BEARER]', 'calling revoke listener', scenarioId)
      callback(false)
    })
  }

  private iframe: HTMLIFrameElement
  private isSessionInitialized: boolean = false
  private allowIntegrationRequests: (initialize: true) => void
  private _maybeInitialized: Promise<boolean>

  constructor(args) {
    this.bearerConfig = new BearerConfig(args || {})
    console.info('[BEARER]', 'config initialized with', args)
    this.maybeInitialized = new Promise((resolve, _reject) => {
      this.allowIntegrationRequests = resolve
    })
    this.initSession()
  }

  authorized = (data: TAuthorizationPayload) => {
    console.log('[BEARER]', 'Bearer:emitting ', { event: Events.AUTHORIZED, ...data })
    Bearer.emitter.emit(Events.AUTHORIZED, data)
  }

  revoked = (data: TAuthorizationPayload) => {
    console.log('[BEARER]', 'Bearer:emitting ', { event: Events.REVOKED, ...data })
    Bearer.emitter.emit(Events.REVOKED, data)
  }

  hasAuthorized = (scenarioId): Promise<boolean> =>
    new Promise((resolve, reject) => {
      postRobot
        .send(this.iframe, Events.HAS_AUTHORIZED, {
          scenarioId,
          clientId: Bearer.config.clientId
        })
        .then(({ data, data: { authorized } }) => {
          console.log('[BEARER]', 'HAS_AUTHORIZED response', data)
          authorized ? resolve(true) : reject(false)
        })
        .catch(iframeError)
    })

  revokeAuthorization = (scenarioId: string): void => {
    postRobot
      .send(this.iframe, Events.REVOKE, {
        scenarioId,
        clientId: Bearer.config.clientId
      })
      .then(() => {
        console.log('[BEARER]', 'Signing out')
      })
      .catch(iframeError)
  }

  initSession() {
    if (window !== undefined && !document.querySelector(`#${IFRAME_NAME}`)) {
      postRobot.on(Events.SESSION_INITIALIZED, event => {
        this.sessionInitialized(event)
      })
      postRobot.on(Events.AUTHORIZED, this.authorized)
      postRobot.on(Events.REVOKED, this.revoked)

      this.iframe = document.createElement('iframe')
      this.iframe.src = `${this.bearerConfig.authorizationHost}v1/user/initialize`
      this.iframe.id = IFRAME_NAME
      this.iframe.width = '0'
      this.iframe.height = '0'
      this.iframe.frameBorder = '0'
      document.body.appendChild(this.iframe)
    }
  }

  askAuthorizations = ({ scenarioId, setupId }): boolean => {
    if (this.isSessionInitialized) {
      const AUTHORIZED_URL = `${Bearer.config.integrationHost}v1/auth/${scenarioId}?setupId=${setupId}`
      window.open(AUTHORIZED_URL, '', 'resizable,scrollbars,status,centerscreen=yes,width=500,height=600')
      return true
    }
    return false
  }

  private sessionInitialized(_event) {
    console.log('[BEARER]', 'session initialized')
    this.isSessionInitialized = true
    this.allowIntegrationRequests(true)
  }
}

function iframeError(e) {
  console.error('[BEARER]', 'Error contacting iframe', e)
}

export default Bearer
