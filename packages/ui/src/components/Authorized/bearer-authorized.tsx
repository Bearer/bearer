import { Component, Method, Prop, State } from '@bearer/core'

import { AuthenticationListener } from '../../utils/with-authentication'

export type FWithAuthenticate = (params: { authenticate(): Promise<boolean> }) => any
export type FWithRevoke = (params: { revoke(authRefId: string): Promise<boolean> }) => any

import debug from '../../logger'
const logger = debug('bearer-authorized')
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
    logger('onAuthorized %s', !!this.pendingAuthorizationResolve)
    this.isAuthorized = true
    if (this.pendingAuthorizationResolve) {
      this.pendingAuthorizationResolve(true)
    }
    this.resetPendingPromises()
  }

  onRevoked = () => {
    this.isAuthorized = false
    logger('onRevoked %s', !!this.pendingAuthorizationReject)
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
        logger('authenticate success %j', data)
      })
      .catch(error => {
        logger('authenticate error %j', error)
      })
  }

  @Method()
  revoke(authRefId: string) {
    this.revokePromise(authRefId)
  }

  authenticatePromise = (authRefId?: string): Promise<boolean> => {
    const promise = new Promise<boolean>((resolve, reject) => {
      this.pendingAuthorizationResolve = resolve
      this.pendingAuthorizationReject = reject
    })
    this.askAuthorization(authRefId)
    return promise
  }

  revokePromise = (authRefId: string): Promise<boolean> => {
    const promise = new Promise<boolean>((resolve, reject) => {
      this.pendingAuthorizationResolve = resolve
      this.pendingAuthorizationReject = reject
    })
    this.revokeAuthorization(authRefId)
    return promise
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
