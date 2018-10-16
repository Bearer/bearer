import { Component, Method, Prop, State } from '@bearer/core'

import { AuthenticationListener } from '../../utils/withAuthentication'

export type FWithAuthenticate = (params: { authenticate(): Promise<boolean> }) => any
export type FWithRevoke = (params: { revoke(): void }) => any

// TODO: scope  authenticatePromise per scenario/setup
// @WithAuthentication()
@Component({
  tag: 'bearer-authorized'
})
export class BearerAuthorized extends AuthenticationListener {
  @State()
  isAuthorized: boolean | null = null

  @State()
  sessionInitialized: boolean = false

  @Prop()
  renderUnauthorized: FWithAuthenticate

  @Prop()
  renderAuthorized: FWithRevoke

  @Prop({ context: 'bearer' })
  bearerContext: any

  @Prop()
  scenarioId: string

  get SCENARIO_ID(): string {
    return this.scenarioId
  }

  private pendingAuthorizationResolve: (authorize: boolean) => void
  private pendingAuthorizationReject: (authorize: boolean) => void

  onAuthorized = () => {
    console.log('[BEARER]', 'onAuthorized', !!this.pendingAuthorizationResolve)
    this.isAuthorized = true
    if (this.pendingAuthorizationResolve) {
      this.pendingAuthorizationResolve(true)
    }
  }

  onRevoked = () => {
    this.isAuthorized = false
    console.log('[BEARER]', 'onRevoked', !!this.pendingAuthorizationReject)
    if (this.pendingAuthorizationReject) {
      this.pendingAuthorizationReject(false)
    }
  }

  onSessionInitialized = () => {
    this.sessionInitialized = true
  }

  @Method()
  authenticate() {
    console.log('[BEARER]', 'bearer-authorized', 'authenticate')
    this.authenticatePromise()
      .then(data => {
        console.log('[BEARER]', 'bearer-authorized', 'authenticated', data)
      })
      .catch(error => {
        console.log('[BEARER]', 'bearer-authenticated', 'error', error)
      })
  }

  @Method()
  revoke() {
    console.log('[BEARER]', 'bearer-authorized', 'revoke')
    this.revokePromise()
  }

  authenticatePromise = (): Promise<boolean> => {
    const promise = new Promise<boolean>((resolve, reject) => {
      this.pendingAuthorizationResolve = resolve
      this.pendingAuthorizationReject = reject
    })
    this.askAuthorization()
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
}
