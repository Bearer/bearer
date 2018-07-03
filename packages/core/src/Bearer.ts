import BearerConfig from './BearerConfig'
import fbemitter from 'fbemitter'
import Events from './EventNames'
import * as postRobot from 'post-robot'

const BEARER_WINDOW_KEY = 'BEARER'
const IFRAME_NAME = 'BEARER-IFRAME'
const AUTHORIZED = 'true'

class Bearer {
  static emitter: fbemitter.EventEmitter = new fbemitter.EventEmitter()
  private static _instance: Bearer
  public static init(config: any): Bearer {
    if (process.env.NODE_ENV === 'development') {
      if (this._instance) {
        console.warn('One instance is already configured, reaplacing it')
      }
    }

    this._instance = new Bearer(config || {})

    return this.instance
  }

  public static get instance(): Bearer {
    return this._instance || this.init({})
  }

  public static get config(): BearerConfig {
    return this.instance.bearerConfig
  }

  static get version() {
    return 'LIB_VERSION'
  }

  private iframe: HTMLIFrameElement
  private isSessionInitialized: boolean = false
  private allowIntegrationRequests: () => void
  private _maybeInitialized: Promise<boolean>

  constructor(args) {
    this.bearerConfig = new BearerConfig(args || {})
    this.maybeInitialized = new Promise((resolve, reject) => {
      this.allowIntegrationRequests = resolve
    })
    this.initSession()
  }

  get bearerConfig(): BearerConfig {
    return window[BEARER_WINDOW_KEY]
  }

  set bearerConfig(config: BearerConfig) {
    window[BEARER_WINDOW_KEY] = config
  }

  get maybeInitialized(): Promise<boolean> {
    if (!this.isSessionInitialized) {
      console.warn('Waiting Bearer to be initialized')
    }
    return this._maybeInitialized
  }

  set maybeInitialized(promise) {
    this._maybeInitialized = promise
  }

  bearerSessionInitialized = () => {
    console.log('[BEARER]', 'session initialized')
    this.isSessionInitialized = true
    this.allowIntegrationRequests()
  }

  bearerAuthorized = ({ data: { scenarioId } }) => {
    console.log('[BEARER]', 'Scenario authorized', scenarioId)
    localStorage.setItem(authorizedKey(scenarioId), AUTHORIZED)
    Bearer.emitter.emit(Events.SCENARIO_AUTHORIZED, {
      authorized: true,
      scenarioId
    })
  }

  bearerCookieSetup = event => {
    console.log('[BEARER]', 'cookie setup')
    document.cookie = event.data.cookie
    event.data.syncCookies(document.cookie)
  }

  hasAuthorized = scenarioId =>
    localStorage.getItem(authorizedKey(scenarioId)) === AUTHORIZED

  revokeAuthorization = (scenarioId: string): void => {
    localStorage.setItem(authorizedKey(scenarioId), undefined)
    Bearer.emitter.emit(Events.SCENARIO_AUTHORIZED, {
      authorized: this.hasAuthorized(scenarioId),
      scenarioId
    })
  }

  initSession() {
    if (
      typeof window !== 'undefined' &&
      !document.querySelector(`#${IFRAME_NAME}`)
    ) {
      postRobot.on(
        Events.BEARER_SESSION_INITIALIZED,
        this.bearerSessionInitialized
      )
      postRobot.on(Events.BEARER_AUTHORIZED, this.bearerAuthorized)
      postRobot.on(Events.BEARER_COOKIE_SETUP, this.bearerCookieSetup)
      this.iframe = document.createElement('iframe')
      this.iframe.src = `${this.bearerConfig.integrationHost}v1/user/initialize`
      this.iframe.id = IFRAME_NAME
      this.iframe.width = '0'
      this.iframe.height = '0'
      this.iframe.frameBorder = '0'
      // this.iframe.addEventListener('load', this.sessionInitialized, true)
      document.body.appendChild(this.iframe)
    }
  }

  askAuthorizations({ scenarioId, setupId }) {
    if (this.isSessionInitialized) {
      const AUTHORIZED_URL = `${
        Bearer.config.integrationHost
      }v1/auth/${scenarioId}?reference_id=${setupId}`
      window.open(
        AUTHORIZED_URL,
        '',
        'resizable,scrollbars,status,centerscreen=yes,width=500,height=600'
      )
      return true
    }
    return false
  }
}

function authorizedKey(scenarioId) {
  return `authorized_${scenarioId}`
}

export default Bearer
