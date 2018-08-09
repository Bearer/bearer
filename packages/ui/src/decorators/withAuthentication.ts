import Bearer from '@bearer/core'

export default function AuthenticationListener() {
  return function(target) {
    const oldDidLoad = target.prototype.componentDidLoad
    target.prototype.componentDidLoad = function(this: IAuthenticatedLike) {
      console.log('[BEARER]', 'componentDidLoad authentication', this)
      Bearer.instance.maybeInitialized
        .then(() => {
          this.authorizedListener = Bearer.onAuthorized(this.SCENARIO_ID, this.onAuthorized)
          this.revokedListener = Bearer.onRevoked(this.SCENARIO_ID, this.onRevoked)

          if (this.onSessionInitialized) {
            this.onSessionInitialized()
          }

          Bearer.instance
            .hasAuthorized(this.SCENARIO_ID)
            .then(() => {
              console.log('[BEARER]', 'authorized')
              this.onAuthorized()
            })
            .catch(error => {
              console.log('[BEARER]', 'unauthorized', { error })
              this.onRevoked()
            })
        })
        .catch(error => {
          console.error('[BEARER]', 'Could not initialize session', { error })
        })
      if (oldDidLoad) {
        oldDidLoad()
      }
    }

    const componentDidUnload = target.prototype.componentDidUnload
    target.prototype.componentDidUnload = function(this: IAuthenticatedLike) {
      if (this.authorizedListener) {
        this.authorizedListener.remove()
        this.authorizedListener = null
      }
      if (this.revokedListener) {
        this.revokedListener.remove()
        this.revokedListener = null
      }
      componentDidUnload()
    }

    target.prototype.revoke = function() {
      Bearer.instance.revokeAuthorization(this.SCENARIO_ID)
    }

    target.prototype.SCENARIO_ID = 'BEARER_SCENARIO_ID'

    target.prototype.authenticate = function() {
      console.log('[BEARER]', 'this.bearerContext', this, this.bearerContext.setupId)
      Bearer.instance.askAuthorizations({
        scenarioId: this.SCENARIO_ID,
        setupId: this.bearerContext.setupId
      })
    }
  }
}

export interface IAuthenticated {
  onSessionInitialized?(): void
  onAuthorized(): void
  onRevoked(): void
  // revoke(): void
}

export interface IAuthenticatedLike extends IAuthenticated {
  [key: string]: any
}

export class WithAuthenticationMethods {
  // revoke = () => {}
  // authenticate = () => {}
}
