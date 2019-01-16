import { DBClient as CLIENT } from './db-client'

export * from './declaration'
export { SaveState } from './intents/save-state'
export { RetrieveState } from './intents/retrieve-state'
export { FetchData } from './intents/fetch'

// tslint:disable-next-line:variable-name
export const DBClient = CLIENT.instance
