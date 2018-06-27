import { intentRequest } from '../requests'

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

const MISSING_SCENARIO_ID =
  'Scenario ID is missing. Please add @Component decorator above your class definition'

// Usage
// @Intent('intentName') propertyName: BearerFetch
// or
// @Intent('intentNameResource',IntentType.GetResource ) propertyName: BearerFetch
export function Intent(
  intentName: string,
  type: IntentType = IntentType.GetCollection
): IDecorator {
  return function(target: any, key: string): void {
    const getter = (): BearerFetch => {
      if (!target['SCENARIO_ID']) {
        console.warn(MISSING_SCENARIO_ID)
      }

      return function(...args) {
        const scenarioId = target['SCENARIO_ID']

        if (!scenarioId) {
          return Promise.reject(new Error(MISSING_SCENARIO_ID))
        } else {
          const intent = intentRequest({ intentName, scenarioId })
          return IntentMapper[type](intent.apply(null, [...args]))
        }
      }
    }

    const setter = () => {}

    if (delete target[key]) {
      Object.defineProperty(target, key, {
        get: getter,
        set: setter
      })
    }
  }
}

// Usage
// @SaveStateItent() propertyName: BearerFetch
// or
// @SaveStateItent(IntentType.GetResource ) propertyName: BearerFetch
export function SaveStateItent(
  type: IntentType = IntentType.GetCollection
): IDecorator {
  return function(target: any, key: string): void {
    const getter = (): BearerFetch => {
      if (!target['SCENARIO_ID']) {
        console.warn(MISSING_SCENARIO_ID)
      }

      return function(
        params: { body?: any; [key: string]: any } = {},
        init: Object = {}
      ) {
        const scenarioId = target['SCENARIO_ID']

        if (!scenarioId) {
          return Promise.reject(new Error(MISSING_SCENARIO_ID))
        } else {
          const { body, ...query } = params
          const intent = intentRequest({ intentName: 'saveState', scenarioId })
          return IntentMapper[type](
            intent.apply(null, [
              {},
              { ...init, method: 'PUT', body: JSON.stringify(body) }
            ])
          )
        }
      }
    }

    const setter = () => {}

    if (delete target[key]) {
      Object.defineProperty(target, key, {
        get: getter,
        set: setter
      })
    }
  }
}

// Usage
// @RetrieveStateItent() propertyName: BearerFetch
// or
// @RetrieveStateItent(IntentType.GetResource ) propertyName: BearerFetch
export function RetrieveStateItent(
  type: IntentType = IntentType.GetCollection
): IDecorator {
  return Intent('retrieveState', type)
}

export function GetCollectionIntent(promise): Promise<any> {
  return new Promise((resolve, reject) => {
    promise
      .then(collection => {
        resolve({ items: collection })
      })
      .catch(e => {
        reject({ items: [], err: e })
      })
  })
}

export function GetResourceIntent(promise): Promise<any> {
  return new Promise((resolve, reject) => {
    promise
      .then(resource => {
        resolve({ object: resource })
      })
      .catch(e => {
        reject({ object: null, err: e })
      })
  })
}

export interface BearerFetch {
  (...args: any[]): Promise<any>
}

export declare const BearerFetch: BearerFetch
