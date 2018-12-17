import { capitalize } from './string'

export function retrieveIntentName(name: string): string {
  return `retrieve${capitalize(name)}`
}

export function retrieveFetcherName(name: string): string {
  return `fetcherRetrieve${capitalize(name)}`
}
