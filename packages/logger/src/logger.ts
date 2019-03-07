import debug from 'debug'
export default (scope: string) => debug(`bearer:${scope}`)

export class BearerFunctionLogger {
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
    debug('bearer:functions')('%j', {
      data,
      integrationId: this.context.integrationUuid,
      functionName: this.context.functionName
    })
  }
}
