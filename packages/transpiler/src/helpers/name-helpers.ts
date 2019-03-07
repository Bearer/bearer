import { capitalize } from './string'

export function retrieveFunctionName(name: string): string {
  return `retrieve${capitalize(name)}`
}

export function retrieveFetcherName(name: string): string {
  return `fetcherRetrieve${capitalize(name)}`
}

export function loadName(name: string): string {
  return `_load${capitalize(name)}`
}

export function initialName(name: string): string {
  return `${name}Initial`
}

export function idName(name: string): string {
  return `${name}Id`
}
