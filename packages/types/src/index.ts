export * from './authentications'
export * from './function-types'

export interface BearerWindow {
  bearer?: { clientId: string; load(clientId: string): void; refreshIntegrations?(): void }
}
