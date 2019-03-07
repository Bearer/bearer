import { EventEmitter, EventSubscription } from 'fbemitter'
import postRobot from 'post-robot'

import BearerConfig from './bearer-config'
import Events from './event-names'
import { formatQuery } from './utils'
import debug from './logger'

const logger = debug.extend('Bearer')

const warn = (message: string) => {
  logger(`%c ${message}`, 'background: #222; color: #bada55')
}

const BEARER_WINDOW_INSTANCE_KEY = 'BEARER_INSTANCE'
const BEARER_EMITTER = 'BEARER_EMITTER'
const BEARER_WINDOW_KEY = 'BEARER'
const BEARER_CONFIG_KEY = 'BEARER_CONFIG'
const IFRAME_NAME = 'BEARER-IFRAME'
const LOG_LEVEL_KEY = 'LOG_LEVEL'

type TAuthorizationPayload = {
  data: {
    integrationId: string
    authIdentifier?: string
  }
}

export default class Bearer {
  private static set _instance(bearerInstance: Bearer) {
    if (window[BEARER_WINDOW_INSTANCE_KEY]) {
      warn('Replacing bearer instance')
    }
    window[BEARER_WINDOW_INSTANCE_KEY] = bearerInstance
  }

  private static get _instance(): Bearer | undefined {
    return window[BEARER_WINDOW_INSTANCE_KEY]
  }

  public static get emitter(): EventEmitter {
    if (!window[BEARER_EMITTER]) {
      window[BEARER_EMITTER] = new EventEmitter()
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
      warn('Waiting Bearer to be initialized')
    }
    return this._maybeInitialized
  }

  set maybeInitialized(promise) {
    this._maybeInitialized = promise
  }

  public static init(config: any = {}): Bearer {
    if (this._instance) {
      warn('One instance is already configured, reaplacing it')
    }
    this._instance = new Bearer({ ...config, ...(window[BEARER_CONFIG_KEY] || {}) }, window)

    return this._instance
  }

  static onAuthorized = (integrationId: string, callback: (authorize: boolean) => void): EventSubscription => {
    logger('onAuthorized register %s', integrationId)
    return Bearer.emitter.addListener(Events.AUTHORIZED, (data: TAuthorizationPayload) => {
      if (data.data.integrationId === integrationId) {
        logger('onAuthorized authorized %s', integrationId)
        callback(true)
      } else {
        logger('onAuthorized different integrationId %s', integrationId)
      }
    })
  }

  static onRevoked = (integrationId: string, callback: (authorize: boolean) => void): EventSubscription => {
    logger('onRevoked register %s', integrationId)
    return Bearer.emitter.addListener(Events.REVOKED, (data: TAuthorizationPayload) => {
      if (data.data.integrationId === integrationId) {
        logger('onRevoked revoked %s', integrationId)
        callback(false)
      } else {
        logger('onRevoked different integrationId %s', integrationId)
      }
    })
  }
  public allowIntegrationRequests: (initialize: true) => void

  private iframe: HTMLIFrameElement
  private isSessionInitialized = false
  private _maybeInitialized: Promise<boolean>

  constructor(args, private readonly window: Window) {
    this.bearerConfig = new BearerConfig(args || {})
    this.maybeInitialized = new Promise((resolve, _reject) => {
      this.allowIntegrationRequests = resolve
    })
    this.window[LOG_LEVEL_KEY] = this.bearerConfig.postRobotLogLevel
    this.initSession()
  }

  authorized = (data: TAuthorizationPayload) => {
    Bearer.emitter.emit(Events.AUTHORIZED, data)
  }

  revoked = (data: TAuthorizationPayload) => {
    Bearer.emitter.emit(Events.REVOKED, data)
  }

  hasAuthorized = (integrationId: string): Promise<boolean> =>
    new Promise((resolve, reject) => {
      postRobot
        .send(this.iframe, Events.HAS_AUTHORIZED, {
          integrationId,
          clientId: Bearer.config.clientId
        })
        .then(({ data, data: { authorized } }: { data: { authorized: boolean } }) => {
          logger('HAS_AUTHORIZED response %j', data)
          authorized ? resolve(true) : reject(false)
        })
        .catch(iframeError)
    })

  revokeAuthorization = (integrationId: string, authId: string): void => {
    postRobot
      .send(this.iframe, Events.REVOKE, {
        authId,
        integrationId,
        clientId: Bearer.config.clientId
      })
      .then(() => {
        logger('Signing out')
      })
      .catch(iframeError)
  }

  initSession() {
    if (this.window !== undefined && !document.querySelector(`#${IFRAME_NAME}`)) {
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
      this.iframe.style.display = 'none'
      document.body.appendChild(this.iframe)
    }
  }

  askAuthorizations = ({
    integrationId,
    setupId,
    authRefId: authId = ''
  }: {
    integrationId: string
    setupId: string
    authRefId?: string
  }): boolean => {
    if (this.isSessionInitialized) {
      const query = formatQuery({
        setupId,
        authId,
        secured: Bearer.config.secured,
        clientId: this.bearerConfig.clientId
      })
      const AUTHORIZED_URL = `${Bearer.config.integrationHost}v2/auth/${integrationId}?${query}`
      this.window.open(AUTHORIZED_URL, '', 'resizable,scrollbars,status,centerscreen=yes,width=500,height=600')
      return true
    }
    return false
  }

  private sessionInitialized(_event) {
    logger('session initialized')
    this.isSessionInitialized = true
    this.allowIntegrationRequests(true)
  }
}

function iframeError(e) {
  logger('Error contacting iframe')
}
