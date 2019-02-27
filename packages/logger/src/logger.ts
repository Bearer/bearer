import debug from 'debug'
export default (scope: string) => debug(`bearer:${scope}`)

export class BearerIntentLogger {
  context: any
  /**
   * @param  {any} context
   */
  constructor(context: any) {
    this.context = context
  }
  /**
   * @param  {any} data
   */
  log = (data: any) => {
    debug('bearer:intents')('%j', {
      data,
      scenarioId: this.context.integrationUuid,
      intentName: this.context.intentName
    })
  }
}
