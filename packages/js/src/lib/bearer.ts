import debounce from 'debounce'
// must be the same version as the one used within the integation service
import postRobot from 'post-robot'
import { TIntegration } from './types'
import debug from './logger'
import { formatQuery } from './utils'

const logger = debug.extend('Bearer')
const prefix = 'bearer'
const DEFAULT_OPTIONS = {
  secured: undefined,
  integrationHost: 'INTEGRATION_HOST_URL',
  domObserver: true,
  refreshDebounceDelay: 200
}

export default class Bearer {
  secured?: boolean
  config: TBearerOptions = DEFAULT_OPTIONS
  private registeredIntegrations: Record<string, boolean> = {}
  private observer?: MutationObserver
  private debounceRefresh: () => void
  private authorizedListener!: postRobot.Cancellable
  private rejectedListener!: postRobot.Cancellable

  constructor(readonly clientId: string, options: Partial<TBearerOptions> = {}) {
    this.config = { ...DEFAULT_OPTIONS, ...cleanOptions(options) }
    this.secured = this.config.secured
    logger('init bearer instance clientId: %s with config: %j', clientId, this.config)
    this.debounceRefresh = debounce(this.loadMissingIntegrations, this.config.refreshDebounceDelay)
    this.initialIntegrationLoading()
    if (this.config.domObserver) {
      this.registerDomObserver()
    }
  }
  // TODO: move to a dedicated file
  /**
   * `connect` lets you easily retrieve `auth-id` for an integration using OAuth authentication. Before using it, you'll need to generate a `setup-id` with the setup component of your integration
   * @argument {string} integration the identifier of the Integration you want to connect to ex: 12345-attach-github-pull-request
   * @argument {string} setupId Setup's identifier you received earlier, a Bearer reference containing all required information about auth mechanism
   * @argument {Object} options Optional parameters like authId if you already have one
   */
  connect = (integration: string, setupId: string, { authId }: { authId?: string } = {}) => {
    const query = formatQuery({
      setupId,
      authId,
      secured: this.config.secured,
      clientId: this.clientId
    })
    const AUTHORIZED_URL = `${this.config.integrationHost}/v2/auth/${integration}?${query}`
    // TODO: get rid of post robot, too heqvy for our needs
    const promise = new Promise<{ integration: string; authId: string }>((resolve, reject) => {
      // TODO: use constants
      if (this.authorizedListener) {
        debug('canceling previous listener')
        this.authorizedListener.cancel()
        this.rejectedListener.cancel()
      }
      debug('add authorization listeners')
      this.authorizedListener = postRobot.on('BEARER_AUTHORIZED', ({ data }) => {
        debug('Authorized: %s => %j', integration, data)
        resolve({ ...data, integration })
      })
      this.rejectedListener = postRobot.on('BEARER_REJECTED', ({ data }) => {
        debug('Rejected: %s => %j', integration, data)
        reject({ ...data, integration })
      })
    }).then()
    window.open(AUTHORIZED_URL, '', 'resizable,scrollbars,status,centerscreen=yes,width=500,height=600')
    return promise
  }

  private _jsonRequest = async (path: string, { query = {}, params = {} } = {}) => {
    const url = [this.config.integrationHost, path].join('')
    const queryParams = { ...query, clientId: this.clientId, secured: this.config.secured }
    const queryString = buildQuery(cleanOptions(queryParams))

    logger('json request: path %s', path)

    return fetch(`${url}?${queryString}`, {
      method: 'POST',
      body: JSON.stringify(params || {}),
      headers: {
        'content-type': 'application/json'
      },
      credentials: 'include'
    }).then(async response => {
      if (response.status > 399) {
        logger('failing request %j', await response.clone().json())
      }
      return response
    })
  }

  functionFetch = async <DataPayload = any>(
    integrationId: string,
    functionName: string,
    { query = {}, ...params }: { query?: Record<string, string>; [key: string]: any } = {}
  ): Promise<TFetchBearerData<DataPayload>> => {
    const path = `/api/v3/functions/${integrationId}/${functionName}`

    try {
      const response = await this._jsonRequest(path, { query, params })
      const payload = await response.json()
      logger('successful request %j', payload)

      if (!payload.error) {
        const { data, meta: { referenceId } = { referenceId: null } } = payload
        return { data, referenceId }
      } else {
        throw { error: payload.error }
      }
    } catch (error) {
      logger('functionFetch failed %j', error, error.message)
      throw { error }
    }
  }

  /**
   * Retrieve all dom elements starting by bearer- and ask for assets urls if
   */
  loadMissingIntegrations = () => {
    const elements = findElements(document.getElementsByTagName('*'))
    const requestedElements = elements.filter(t => !this.registeredIntegration(t))
    logger(this.registeredIntegrations, elements, requestedElements)
    this.sendTags(requestedElements)
  }

  /**
   * check wether if an integration is resgistered
   */
  registeredIntegration = (tagName: string): boolean => {
    // NOTE:
    // .constructor !== HTMLElement looks weird but it does not work the other way   ¯\_(ツ)_/¯
    // constructor is supposed to be class extends HTMLElement{}
    this.registeredIntegrations[tagName] =
      this.registeredIntegrations[tagName] || document.createElement(tagName).constructor !== HTMLElement
    return this.registeredIntegrations[tagName]
  }

  /**
   * load integration asset or wait until dom is loaded
   */
  private initialIntegrationLoading = () => {
    if (document.readyState === 'complete' || document.readyState === 'interactive') {
      this.debounceRefresh()
    } else {
      document.addEventListener('DOMContentLoaded', this.debounceRefresh)
    }
  }

  /**
   * Register a DOM observer so that we can load integration assets only when we need them
   */
  private registerDomObserver = () => {
    if ('MutationObserver' in window) {
      this.disconnectObserver()

      const container = document.documentElement || document.body
      const config = { childList: true, subtree: true }

      this.observer = new MutationObserver(this.observerCallback)
      this.observer.observe(container, config)
    }
  }

  private observerCallback = (mutations: MutationRecord[]) => {
    for (const mutation of mutations) {
      if (mutation.type == 'childList') {
        if (mutation.addedNodes.length) {
          this.debounceRefresh()
        }
      }
    }
  }

  /**
   * remove dom observer
   */
  private disconnectObserver = () => {
    if (this.observer) {
      this.observer.disconnect()
      delete this.observer
    }
  }

  /**
   * retrieve corresponding integration asset url
   */
  private sendTags = async (tags: string[]): Promise<boolean> => {
    if (!tags.length) {
      return Promise.resolve(true)
    }
    try {
      const response = await fetch(`${this.config.integrationHost}/v1/parse-tags`, {
        headers: { 'content-type': 'application/json' },
        body: JSON.stringify({ tags, clientId: this.clientId }),
        method: 'POST'
      })
      if (response.status > 299) {
        throw new Error(`Error while fetching integration tag names: ${tags}`)
      }
      const integrations: TIntegration[] = await response.json()

      integrations.map(integration => {
        if (!document.querySelector(`#${getScriptId(integration.uuid)}`)) {
          document.body.appendChild(getScriptDOM(this.clientId, integration))
        }
      })
      return true
    } catch (e) {
      return false
    }
  }
}

export type TBearerOptions = {
  secured?: boolean
  domObserver: boolean
  integrationHost: string
  refreshDebounceDelay: number
}

/**
 * Extract/format element tag names given a regexp
 * @param elements
 * @param filter
 */
export function findElements(elements: HTMLCollection | NodeListOf<Element>, filter: RegExp = /^bearer-/i): string[] {
  return Array.from(elements)
    .filter(el => filter.test(el.tagName))
    .map(el => el.tagName.toLowerCase())
}

/**
 * Return the bearer script id
 * @param uuid
 */
function getScriptId(uuid: string): string {
  return `${prefix}-${uuid}` // id must start with a letter
}

/**
 * create a script tag for a given integration
 * @param clientId
 * @param integration
 */
function getScriptDOM(clientId: string, integration: TIntegration): HTMLScriptElement {
  const s = document.createElement('script')
  s.type = 'text/javascript'
  s.async = true
  const separator = integration.asset.indexOf('?') > -1 ? '&' : '?'
  s.src = [integration.asset, [`clientId=${clientId}`].join('&')].join(separator)
  s.id = getScriptId(integration.uuid)
  return s
}

function buildQuery(params: any) {
  function encode(k: string) {
    return encodeURIComponent(k) + '=' + encodeURIComponent(params[k])
  }
  return Object.keys(params)
    .map(encode)
    .join('&')
}

function cleanOptions(obj: Record<string, any>) {
  return Object.keys(obj).reduce(
    (acc, key: string) => {
      if (obj[key] !== undefined) {
        acc[key] = obj[key]
      }
      return acc
    },
    {} as Record<string, any>
  )
}

export type TFetchBearerData<T = any> = { data: T; referenceId?: string }
