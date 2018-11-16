export * from './authentications'
export * from './intent-types'

export interface BearerWindow {
  bearer?: { clientId: string; load(clientId: string): void; refreshScenarios?(): void }
}
