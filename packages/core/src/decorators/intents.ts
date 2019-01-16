import { intentRequest } from '../requests'

/**
 * Declarations
 */

export enum IntentType {
  RetrieveState = 'RetrieveState',
  SaveState = 'SaveState',
  FetchData = 'FetchData'
}

type TFetchBearerResult<T = any> = { meta: { referenceId: string }; data: T; error?: any }

export type TFetchBearerData<T = any> = { data: T; referenceId?: string }

export type BearerFetch<T = any> = (...args: any[]) => Promise<TFetchBearerData<T>>

type IDecorator = (target: any, key: string) => void

type BearerComponent = {
  setupId: string
  SCENARIO_ID: string
  referenceId: string
  el?: HTMLElement
}

/**
 * Constants
 */

export const BearerContext = 'bearerContext'
export const setupId = 'setupId'
export const IntentSaved = 'BearerStateSaved'
export const BearerStateSavedEvent = 'bearer:StateSaved'

/**
 * Intent
 */

// Usage
// @Intent('intentName') propertyName: BearerFetch
// or
// @Intent('intentNameResource', IntentType.FetchData) propertyName: BearerFetch
export function Intent(intentName: string, type: IntentType = IntentType.FetchData): IDecorator {
  return function(target: BearerComponent, key: string): void {
    const getter = (): BearerFetch => {
      return function(this: BearerComponent, params = {}): Promise<TFetchBearerData> {
        // NOTE: here we have to use target. Not sure why
        const scenarioId = target.SCENARIO_ID
        if (!scenarioId) {
          return missingScenarioId()
        }

        const intent = intentRequest<TFetchBearerResult>({
          intentName,
          scenarioId,
          [setupId]: params[setupId] || retrieveSetupId(target)
        })

        // prepare params and body

        const referenceId = retrieveReferenceId(this)
        if (type === IntentType.RetrieveState && !referenceId) {
          return missingReferenceId()
        }

        const { body, ...queryParams } = params
        const baseQuery = referenceId ? { referenceId } : {}
        const query = { ...baseQuery, ...queryParams }
        const init = { method: 'POST', body: JSON.stringify(body || {}) }

        // Build promise
        return intentPromise(intent(query, init))
      }
    }

    defineIntentProp(target, key, getter)
  }
}

export function intentPromise(promise: Promise<TFetchBearerResult>): Promise<TFetchBearerData> {
  return new Promise((resolve, reject) => {
    promise
      .then((payload: TFetchBearerResult) => {
        if (payload.error) {
          reject({ error: payload.error })
        } else {
          const { data, meta: { referenceId } = { referenceId: null } } = payload
          resolve({ data, referenceId })
        }
      })
      .catch(error => {
        reject({ error })
      })
  })
}

/**
 * Helpers
 */

function missingScenarioId(): Promise<any> {
  console.info('[BEARER]', 'Missing scenarioId, skipping api call')
  return Promise.reject(new BearerMissingScenarioId())
}

function missingReferenceId(): Promise<any> {
  console.info('[BEARER]', 'Missing referenceId, skipping RetrieveState api call')
  return Promise.reject(new BearerMissingReferenceId())
}

function retrieveSetupId(target: BearerComponent): string {
  return target.setupId || (target[BearerContext] && target[BearerContext][setupId])
}

function retrieveReferenceId(target: BearerComponent): string {
  return target.referenceId
}

function defineIntentProp(target: BearerComponent, key: string, getter: any): void {
  const setter = () => {}

  if (delete target[key]) {
    Object.defineProperty(target, key, {
      get: getter,
      set: setter
    })
  }
}

/**
 * Custom Errors
 */

class BearerMissingReferenceId extends Error {
  message = `Attribute ${this.group}ReferenceId is missing. Cannot fetch data without any reference`
  constructor(private readonly group: string = 'feature') {
    super()
  }
}

class BearerMissingScenarioId extends Error {
  message = 'Scenario ID is missing. Please add @RootComponent decorator above your class definition'
}
