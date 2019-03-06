import Bearer from '@bearer/core'

import debug from '../logger'
const logger = debug('AuthenticationListener')

export class AuthenticationListener {
  protected INTEGRATION_ID!: string
  protected onAuthorized: () => void
  protected onRevoked: () => void
  protected onSessionInitialized!: () => void
  protected bearerContext: any
  protected authorizedListener: any
  protected revokedListener: any

  askAuthorization = (authRefId?: string) => {
    logger('authenticate %s %s', this.INTEGRATION_ID, this.bearerContext.setupId)
    Bearer.instance.askAuthorizations({
      authRefId,
      integrationId: this.INTEGRATION_ID,
      setupId: this.bearerContext.setupId
    })
  }

  revokeAuthorization = (authRefId?: string) => {
    Bearer.instance.revokeAuthorization(this.INTEGRATION_ID, authRefId)
  }

  componentDidLoad() {
    Bearer.instance.maybeInitialized
      .then(() => {
        if (this.onSessionInitialized) {
          this.onSessionInitialized()
        }

        this.authorizedListener = Bearer.onAuthorized(this.INTEGRATION_ID, this.onAuthorized)
        this.revokedListener = Bearer.onRevoked(this.INTEGRATION_ID, this.onRevoked)

        Bearer.instance
          .hasAuthorized(this.INTEGRATION_ID)
          .then(() => {
            logger('authorized')
            this.onAuthorized()
          })
          .catch(error => {
            logger('unauthorized %j', { error })
            this.onRevoked()
          })
      })
      .catch(error => {
        logger.extend('error')('Could not initialize session %j', { error })
      })
  }

  componentDidUnload = () => {
    if (this.authorizedListener) {
      this.authorizedListener.remove()
      this.authorizedListener = null
    }
    if (this.revokedListener) {
      this.revokedListener.remove()
      this.revokedListener = null
    }
  }
}
