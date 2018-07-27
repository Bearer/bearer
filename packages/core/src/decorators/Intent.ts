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

const MISSING_SCENARIO_ID = 'Scenario ID is missing. Please add @Component decorator above your class definition'
// Usage
// @Intent('intentName') propertyName: BearerFetch
// or
// @Intent('intentNameResource',IntentType.GetResource ) propertyName: BearerFetch
export function Intent(intentName: string, type: IntentType = IntentType.GetCollection): IDecorator {
  return function(target: any, key: string): void {
    const getter = (): BearerFetch => {
      return function(params = {}, init = {}) {
        if (!target['SCENARIO_ID']) {
          console.warn(MISSING_SCENARIO_ID)
        }
        const scenarioId = target['SCENARIO_ID']
        if (!scenarioId) {
          return Promise.reject(new Error(MISSING_SCENARIO_ID))
        } else {
          const intent = intentRequest({
            intentName,
            scenarioId,
            setupId: retrieveSetupId(target)
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
// @SaveStateIntent() propertyName: BearerFetch
// or
// @SaveStateIntent(IntentType.GetResource ) propertyName: BearerFetch
export function SaveStateIntent(type: IntentType = IntentType.GetCollection): IDecorator {
  return function(target: any, key: string): void {
    const getter = (): BearerFetch => {
      return function(params: { body?: any; [key: string]: any } = {}, init: Object = {}) {
        const scenarioId = target['SCENARIO_ID']

        if (!scenarioId) {
          return Promise.reject(new Error(MISSING_SCENARIO_ID))
        } else {
          const { body, ...query } = params
          const intent = intentRequest({
            intentName: 'SaveState',
            scenarioId,
            setupId: retrieveSetupId(this)
          })
          const referenceId = retrieveReferenceId(this)
          const baseQuery = referenceId ? { referenceId } : {}

          return IntentMapper[type](
            intent.apply(null, [{ ...query, ...baseQuery }, { ...init, method: 'PUT', body: JSON.stringify(body) }])
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

function retrieveSetupId(target: any) {
  const setupId = target.setupId || (target['bearerContext'] && target['bearerContext']['setupId'])
  return setupId
}

function retrieveReferenceId(target: any) {
  return target.referenceId
}

// Usage
// @RetrieveStateIntent() propertyName: BearerFetch
// or
// @RetrieveStateIntent(IntentType.GetResource ) propertyName: BearerFetch
export function RetrieveStateIntent(type: IntentType = IntentType.GetCollection): IDecorator {
  return Intent('RetrieveState', type)
}

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

export interface BearerFetch {
  (...args: any[]): Promise<any>
}

export declare const BearerFetch: BearerFetch
