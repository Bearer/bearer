import { Bearer } from './bearer'

import debug from './logger'
import { buildQuery, cleanOptions } from './utils'

export class IntegrationClient {
  private logger: debug.Debugger

  constructor(private readonly bearerInstance: Bearer, readonly integrationId: string) {
    this.logger = debug.extend('integration-client').extend(integrationId)
  }

  invoke = async <DataPayload = any>(
    functionName: string,
    { query = {}, ...params }: { query?: Record<string, string>; [key: string]: any } = {}
  ): Promise<TFetchBearerData<DataPayload>> => {
    const path = `/api/v4/functions/${this.integrationId}/${functionName}`
    try {
      const response = await this._jsonRequest(path, { query, params })
      const payload: {
        data: any
        referenceId?: string
        error: any
      } = await response.json()
      this.logger('successful request %j', payload)

      if (!payload.error) {
        return payload
      } else {
        throw { error: payload.error }
      }
    } catch (error) {
      this.logger('invoke failed %j', error, error.message)
      throw { error }
    }
  }

  private _jsonRequest = async (path: string, { query = {}, params = {} } = {}) => {
    const url = [this.bearerInstance.config.integrationHost, path].join('')
    const queryParams = {
      ...query,
      clientId: this.bearerInstance.clientId,
      secured: this.bearerInstance.config.secured
    }
    const queryString = buildQuery(cleanOptions(queryParams))

    this.logger('json request: path %s', path)

    return fetch(`${url}?${queryString}`, {
      method: 'POST',
      body: JSON.stringify(params || {}),
      headers: {
        'content-type': 'application/json'
      },
      credentials: 'include'
    }).then(async response => {
      if (response.status > 399) {
        this.logger('failing request %j', await response.clone().json())
      }
      return response
    })
  }
}

export type TFetchBearerData<T = any> = { data: T; referenceId?: string }

export default IntegrationClient
