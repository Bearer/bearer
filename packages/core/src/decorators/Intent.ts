import { intentRequest } from '../requests'

export enum IntentType {
  GetCollection = 'GetCollection',
  GetResource = 'GetResource'
}

const IntentMapper = {
  [IntentType.GetCollection]: GetCollectionIntent,
  [IntentType.GetResource]: GetResourceIntent
}

const MISSING_SCENARIO_ID =
  'Scenario ID is missing. Add BearerComponent above @Component({...}) decorator'

// Property Decorator
// How to use it
// @Intent('intentName') propertyName: BearerFetch
export function Intent(
  intentName: string,
  type: IntentType = IntentType.GetCollection
): (target: any, key: string) => void {
  return function(target: any, key: string) {
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

export function SaveStateItent(
  type: IntentType = IntentType.GetCollection
): (target: any, key: string) => void {
  return Intent('saveState', type)
}

export function RetrieveStateItent(
  type: IntentType = IntentType.GetCollection
): (target: any, key: string) => void {
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
