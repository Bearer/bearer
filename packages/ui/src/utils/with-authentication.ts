import Bearer from '@bearer/core'

export class AuthenticationListener {
  protected SCENARIO_ID!: string
  protected onAuthorized: () => void
  protected onRevoked: () => void
  protected onSessionInitialized!: () => void
  protected bearerContext: any
  protected authorizedListener: any
  protected revokedListener: any

  askAuthorization = (authRefId?: string) => {
    console.debug('[BEARER]', 'authenticate', this.SCENARIO_ID, this.bearerContext.setupId)
    Bearer.instance.askAuthorizations({
      scenarioId: this.SCENARIO_ID,
      setupId: this.bearerContext.setupId,
      authRefId
    })
  }

  revokeAuthorization = () => {
    Bearer.instance.revokeAuthorization(this.SCENARIO_ID)
  }

  componentDidLoad() {
    Bearer.instance.maybeInitialized
      .then(() => {
        if (this.onSessionInitialized) {
          this.onSessionInitialized()
        }

        this.authorizedListener = Bearer.onAuthorized(this.SCENARIO_ID, this.onAuthorized)
        this.revokedListener = Bearer.onRevoked(this.SCENARIO_ID, this.onRevoked)

        Bearer.instance
          .hasAuthorized(this.SCENARIO_ID)
          .then(() => {
            console.debug('[BEARER]', 'authorized')
            this.onAuthorized()
          })
          .catch(error => {
            console.debug('[BEARER]', 'unauthorized', { error })
            this.onRevoked()
          })
      })
      .catch(error => {
        console.error('[BEARER]', 'Could not initialize session', { error })
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
