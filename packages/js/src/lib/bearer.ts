import debounce from 'debounce'
import { TIntegration } from './types'

const prefix = 'bearer'
const DEFAULT_OPTIONS = {
  domObserver: true,
  refreshDebounceDelay: 200
}

export default class Bearer {
  // @ts-ignore for now
  private registeredIntegrations: Record<string, bolean> = {}
  private config: TBearerOptions = DEFAULT_OPTIONS
  private observer?: MutationObserver
  private debounceRefresh: () => void

  constructor(readonly clientId: string, options?: Partial<TBearerOptions>) {
    this.config = { ...options, ...DEFAULT_OPTIONS }
    this.debounceRefresh = debounce(this.loadMissingIntegrations, this.config.refreshDebounceDelay)
    this.initialIntegrationLoading()
    if (this.config.domObserver) {
      this.registerDomObserver()
    }
  }

  /**
   * Retrieve all dom elements starting by bearer- and ask for assets urls if
   */
  loadMissingIntegrations = () => {
    const elements = findElements(document.getElementsByTagName('*'))
    this.sendTags(elements.filter(this.registeredIntegration))
  }

  /**
   * check wether if an integration is resgistered
   */
  registeredIntegration = (tagName: string): boolean => {
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
  sendTags = async (tags: string[]): Promise<boolean> => {
    if (!tags.length) {
      return Promise.resolve(true)
    }
    try {
      const response = await fetch('BEARER_PARSE_TAGS_URI', {
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

type TBearerOptions = {
  domObserver: boolean
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
