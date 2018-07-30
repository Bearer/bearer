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

interface IDecorator {
  (target: any, key: string): void
}

export interface BearerFetch {
  (...args: any[]): Promise<any>
}

export declare const BearerFetch: BearerFetch

export const BearerContext = 'bearerContext'
export const setupId = 'setupId'

/**
 * Intents
 */

// Usage
// @Intent('intentName') propertyName: BearerFetch
// or
// @Intent('intentNameResource',IntentType.GetResource ) propertyName: BearerFetch
export function Intent(intentName: string, type: IntentType = IntentType.GetCollection): IDecorator {
  return function(target: any, key: string): void {
    const getter = (): BearerFetch => {
      return function(params = {}, init = {}) {
        const scenarioId = target['SCENARIO_ID']
        if (!scenarioId) {
          return missingScenarioId()
        }
        const intent = intentRequest({
          intentName,
          scenarioId,
          [setupId]: retrieveSetupId(target)
        })
        const referenceId = retrieveReferenceId(this)
        const baseQuery = referenceId ? { referenceId } : {}

        return IntentMapper[type](
          intent.apply(null, [
            {
              ...params,
              ...baseQuery
            },
            init
          ])
        )
      }
    }

    defineIntentProp(target, key, getter)
  }
}

// Usage
// @SaveStateIntent() propertyName: BearerFetch
// or
// @SaveStateIntent(IntentType.GetResource ) propertyName: BearerFetch
export function SaveStateIntent(type: IntentType = IntentType.GetCollection): IDecorator {
  return function(target: any, key: string): void {
    const getter = (): BearerFetch => {
      return function(params: { body?: any; [key: string]: any } = {}, init: Object = {}) {
        const scenarioId = target['SCENARIO_ID']
        if (!scenarioId) {
          return missingScenarioId()
        }

        const { body, ...query } = params
        const intent = intentRequest({
          intentName: IntentNames.SaveState,
          scenarioId,
          [setupId]: retrieveSetupId(this)
        })
        const referenceId = retrieveReferenceId(this)
        const baseQuery = referenceId ? { referenceId } : {}

        return IntentMapper[type](
          intent.apply(null, [{ ...query, ...baseQuery }, { ...init, method: 'PUT', body: JSON.stringify(body) }])
        )
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
  return function(target: any, key: string): void {
    const getter = (): BearerFetch => {
      return function(params: { body?: any; [key: string]: any } = {}, init: Object = {}) {
        const scenarioId = target['SCENARIO_ID']
        if (!scenarioId) {
          return missingScenarioId()
        }

        const referenceId = retrieveReferenceId(this)
        if (!referenceId) {
          return missingReferenceId()
        }

        const { body, ...query } = params
        const intent = intentRequest({
          intentName: IntentNames.RetrieveState,
          scenarioId,
          [setupId]: retrieveSetupId(this)
        })

        return IntentMapper[type](
          intent.apply(null, [{ ...query, referenceId }, { ...init, method: 'PUT', body: JSON.stringify(body) }])
        )
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

export function GetCollectionIntent(promise): Promise<any> {
  return new Promise((resolve, reject) => {
    promise
      .then(({ data, collection }: { data: Array<any>; collection: Array<any> }) => {
        resolve({ items: data || collection })
      })
      .catch(e => {
        reject({ items: [], err: e })
      })
  })
}

export function GetResourceIntent(promise): Promise<any> {
  return new Promise((resolve, reject) => {
    promise
      .then(({ data, referenceId }: { data: any; referenceId?: string }) => {
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

function retrieveSetupId(target: any) {
  return target.setupId || (target[BearerContext] && target[BearerContext][setupId])
}

function retrieveReferenceId(target: any) {
  return target.referenceId
}

function defineIntentProp(target: any, key: string, getter: any) {
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
