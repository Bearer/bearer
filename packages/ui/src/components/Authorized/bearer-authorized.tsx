import { Component, Method, Prop, State } from '@bearer/core'

import { AuthenticationListener } from '../../utils/with-authentication'

export type FWithAuthenticate = (params: { authenticate(): Promise<boolean> }) => any
export type FWithRevoke = (params: { revoke(): Promise<boolean> }) => any

// TODO: scope  authenticatePromise per scenario/setup
// @WithAuthentication()
@Component({
  tag: 'bearer-authorized'
})
export class BearerAuthorized extends AuthenticationListener {
  get SCENARIO_ID(): string {
    return this.scenarioId
  }
  @State()
  isAuthorized: boolean | null = null

  @State()
  sessionInitialized = false

  @Prop()
  renderUnauthorized: FWithAuthenticate

  @Prop()
  renderAuthorized: FWithRevoke

  @Prop({ context: 'bearer' })
  bearerContext: any

  @Prop()
  scenarioId: string

  private pendingAuthorizationResolve: (authorize: boolean) => void
  private pendingAuthorizationReject: (authorize: boolean) => void

  onAuthorized = () => {
    console.debug('[BEARER]', 'onAuthorized', !!this.pendingAuthorizationResolve)
    this.isAuthorized = true
    if (this.pendingAuthorizationResolve) {
      this.pendingAuthorizationResolve(true)
    }
    this.resetPendingPromises()
  }

  onRevoked = () => {
    this.isAuthorized = false
    console.debug('[BEARER]', 'onRevoked', !!this.pendingAuthorizationReject)
    if (this.pendingAuthorizationReject) {
      this.pendingAuthorizationReject(false)
    }
    this.resetPendingPromises()
  }

  onSessionInitialized = () => {
    this.sessionInitialized = true
  }

  @Method()
  authenticate(authRefId?: string) {
    this.authenticatePromise(authRefId)
      .then(data => {
        console.debug('[BEARER]', 'bearer-authorized', 'authenticated', data)
      })
      .catch(error => {
        console.debug('[BEARER]', 'bearer-authenticated', 'error', error)
      })
  }

  @Method()
  revoke() {
    this.revokePromise()
  }

  authenticatePromise = (authRefId?: string): Promise<boolean> => {
    const promise = new Promise<boolean>((resolve, reject) => {
      this.pendingAuthorizationResolve = resolve
      this.pendingAuthorizationReject = reject
    })
    this.askAuthorization(authRefId)
    return promise
  }

  revokePromise = (): Promise<boolean> => {
    this.revokeAuthorization()
    return Promise.resolve(true)
  }

  render() {
    if (!this.sessionInitialized || this.isAuthorized === null) {
      return null
    }
    if (!this.isAuthorized) {
      return this.renderUnauthorized
        ? this.renderUnauthorized({ authenticate: this.authenticatePromise })
        : 'Unauthorized'
    }
    return this.renderAuthorized ? this.renderAuthorized({ revoke: this.revokePromise }) : <slot />
  }

  private resetPendingPromises = () => {
    this.pendingAuthorizationResolve = null
    this.pendingAuthorizationReject = null
  }
}
