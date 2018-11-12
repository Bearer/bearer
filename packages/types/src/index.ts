export * from './authentications'
export * from './IntentTypes'

export interface BearerWindow {
  bearer?: { clientId: string; load(clientId: string): void; refreshScenarios?(): void }
}
