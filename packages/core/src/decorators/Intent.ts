import { intentRequest } from '../requests'

/**
 * Declarations
 */
enum IntentNames {
  RetrieveState = 'RetrieveState',
  SaveState = 'SaveState'
}

export enum IntentType {
  GetCollection = 'GetCollection',
  GetResource = 'GetResource'
}

const IntentMapper = {
  [IntentType.GetCollection]: GetCollectionIntent,
  [IntentType.GetResource]: GetResourceIntent
}

type TPayloadResource = { meta: { referenceId?: string }; data: any }
type TPayloadCollection = { meta: { referenceId: string }; data: Array<any> }
type TFetchBearerResult = TPayloadResource | TPayloadCollection

export type TCollectionData = { items: Array<any>; referenceId?: string }
export type TResourceData = { object: any; referenceId?: string }
export type TFetchBearerData = TCollectionData | TResourceData

export interface BearerFetch {
  (...args: any[]): Promise<any>
}

interface IDecorator {
  (target: any, key: string): void
}

type BearerComponent = {
  BEARER_ID: string
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
 * Intents
 */

// Usage
// @Intent('intentName') propertyName: BearerFetch
// or
// @Intent('intentNameResource',IntentType.GetResource ) propertyName: BearerFetch
export function Intent(intentName: string, type: IntentType = IntentType.GetCollection): IDecorator {
  return function(target: BearerComponent, key: string): void {
    const getter = (): BearerFetch => {
      return function(this: BearerComponent, params = {}, init = {}) {
        // NOTE: here we have to use target. Not sure why
        const scenarioId = target.SCENARIO_ID
        if (!scenarioId) {
          return missingScenarioId()
        }
        const intent = intentRequest<TFetchBearerResult>({
          intentName,
          scenarioId,
          [setupId]: retrieveSetupId(target)
        })
        const referenceId = retrieveReferenceId(this)
        const baseQuery = referenceId ? { referenceId } : {}
        const promise: Promise<TFetchBearerResult> = intent.apply(null, [{ ...params, ...baseQuery }, init])
        return IntentMapper[type](promise)
      }
    }

    defineIntentProp(target, key, getter)
  }
}

type TEventPayload = {
  referenceId: string
  object: { referenceId: string; ReadAllowed: boolean; [key: string]: any }
}
// Usage
// @SaveStateIntent() propertyName: BearerFetch
// or
// @SaveStateIntent(IntentType.GetResource ) propertyName: BearerFetch
export function SaveStateIntent(): IDecorator {
  return function(target: BearerComponent, key: string): void {
    const getter = (): BearerFetch => {
      return function(this: BearerComponent, params: { body?: any; [key: string]: any } = {}, init: Object = {}) {
        const scenarioId = this.SCENARIO_ID
        if (!scenarioId) {
          return missingScenarioId()
        }

        const { body, ...query } = params
        const intent = intentRequest<TFetchBearerResult>({
          intentName: IntentNames.SaveState,
          scenarioId,
          [setupId]: retrieveSetupId(this)
        })
        const referenceId = retrieveReferenceId(this)
        const baseQuery = referenceId ? { referenceId } : {}
        const promise: Promise<TFetchBearerResult> = intent.apply(null, [
          { ...query, ...baseQuery },
          { ...init, method: 'PUT', body: JSON.stringify(body) }
        ])

        return new Promise((resolve, reject) => {
          // It does not make sense to use collection here.
          GetResourceIntent(promise)
            .then((payload: TEventPayload) => {
              if (this.el) {
                this.el.dispatchEvent(new CustomEvent(BearerStateSavedEvent, { detail: payload }))
              }
              resolve(payload)
              return payload
            })
            .catch(reject)
        })
      }
    }
    defineIntentProp(target, key, getter)
  }
}

// Usage
// @RetrieveStateIntent() propertyName: BearerFetch
// or
// @RetrieveStateIntent(IntentType.GetResource ) propertyName: BearerFetch
export function RetrieveStateIntent(type: IntentType = IntentType.GetCollection): IDecorator {
  return function(target: BearerComponent, key: string): void {
    const getter = (): BearerFetch => {
      return function(this: BearerComponent, params: { body?: any; [key: string]: any } = {}, init: Object = {}) {
        const scenarioId = this.SCENARIO_ID
        if (!scenarioId) {
          return missingScenarioId()
        }

        const referenceId = retrieveReferenceId(this)
        if (!referenceId) {
          return missingReferenceId()
        }

        const { body, ...query } = params
        const intent = intentRequest<TFetchBearerResult>({
          intentName: IntentNames.RetrieveState,
          scenarioId,
          [setupId]: retrieveSetupId(this)
        })
        const promise: Promise<TFetchBearerResult> = intent.apply(null, [
          { ...query, referenceId },
          { ...init, method: 'PUT', body: JSON.stringify(body) }
        ])
        return IntentMapper[type](promise)
      }
    }

    defineIntentProp(target, key, getter)
  }
}

// Common
// => reject if setupId/scenarioId/integrationId/query params is missing
//

// Intent
// => pass setupId/scenarioId/integrationId/query params
// => pass query params

// SaveStateIntent
// => pass setupId/scenarioId/integrationId/query params
// => trigger Save

// RetrieveStateIntent
// => do not perform if no referenceId
// => pass setupId/scenarioId/integrationId/query params

export function GetCollectionIntent(promise: Promise<TFetchBearerResult>): Promise<TCollectionData> {
  return new Promise((resolve, reject) => {
    promise
      .then((payload: TPayloadCollection) => {
        const { data, meta: { referenceId } = { referenceId: null } } = payload
        resolve({ items: data, referenceId })
      })
      .catch(e => {
        reject({ items: [], err: e })
      })
  })
}

export function GetResourceIntent(promise: Promise<TFetchBearerResult>): Promise<TResourceData> {
  return new Promise((resolve, reject) => {
    promise
      .then((payload: TPayloadResource) => {
        const { data, meta: { referenceId } = { referenceId: null } } = payload
        resolve({ object: data, referenceId })
      })
      .catch(e => {
        reject({ object: null, err: e })
      })
  })
}

/**
 * Helpers
 */

function missingScenarioId(): Promise<any> {
  console.warn('[BEARER]', 'Missing scenarioId, skipping api call')
  return Promise.reject(new BearerMissingScenarioId())
}

function missingReferenceId(): Promise<any> {
  console.warn('[BEARER]', 'Missing referenceId, skipping RetrieveState api call')
  return Promise.reject(new BearerMissingReferenceId())
}

function retrieveSetupId(target: BearerComponent) {
  return target.setupId || (target[BearerContext] && target[BearerContext][setupId])
}

function retrieveReferenceId(target: BearerComponent) {
  return target.referenceId
}

function defineIntentProp(target: BearerComponent, key: string, getter: any) {
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
  constructor(private readonly group: string = 'feature') {
    super()
  }
  message = `Attribute ${this.group}ReferenceId is missing. Cannot fetch data without any reference`
}

class BearerMissingScenarioId extends Error {
  message = 'Scenario ID is missing. Please add @RootComponent decorator above your class definition'
}
