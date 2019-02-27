import debug from 'debug'
export default (scope: string) => debug(`bearer:${scope}`)

export class BearerIntentLogger {
  context: any

  constructor(context: any) {
    this.context = context
  }

  log = (message: any) => {
    debug('bearer:intents')('%j', {
      data: message,
      scenarioId: this.context.integrationUuid,
      intentName: this.context.intentName
    })
  }
}
