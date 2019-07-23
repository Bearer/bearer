import axios, { AxiosRequestConfig } from 'axios'

import { Bearer } from './bearer'

import debug from './logger'
import { cleanOptions } from './utils'

export class IntegrationClient {
  private logger: any

  constructor(
    private readonly bearerInstance: Bearer,
    readonly integrationId: string,
    readonly setupId?: string,
    readonly authId?: string
  ) {
    this.logger = debug.extend('integration-client').extend(integrationId)
    this.logger('Integration initialized', this)
  }

  /**
   * `auth` set authId so that API calls are performed with the given identity
   */
  auth = (authId: string) => new IntegrationClient(this.bearerInstance, this.integrationId, this.setupId, authId)

  /**
   * `setup` specify which setupId to use when calling Integration service
   */
  setup = (setupId: string) => new IntegrationClient(this.bearerInstance, this.integrationId, setupId, this.authId)

  /**
   * `invoke` invoke a custom function for the given integration
   * @param functionName {string} the name of the function you would like to invoke
   * @param params {object} extra information required when invoking the function
   */
  invoke = async <DataPayload = any>(
    functionName: string,
    { query = {}, ...params }: { query?: Record<string, string>; [key: string]: any } = {}
  ): Promise<TFetchBearerData<DataPayload>> => {
    const queryParams = {
      authId: this.authId,
      setupId: this.setupId,
      ...query,
      clientId: this.bearerInstance.clientId,
      secured: this.bearerInstance.config.secured
    }
    this.logger('json request: path %s', functionName)

    try {
      const { data: payload } = await axios.request<{
        data: any
        referenceId?: string
        error: any
      }>({
        method: 'POST',
        baseURL: `${this.bearerInstance.config.integrationHost}/api/v4/functions/${this.integrationId}`,
        url: functionName,
        params: cleanOptions(queryParams),
        data: params || {}
      })
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

  /**
   * `get` perform get request to integration service
   */

  get = <DataReturned = any>(endpoint: string, parameters?: BearerRequestParameters) => {
    return this.request<DataReturned>('GET', endpoint, parameters)
  }

  /**
   * `head` perform head request to integration service
   */

  head = <DataReturned = any>(endpoint: string, parameters?: BearerRequestParameters) => {
    return this.request<DataReturned>('HEAD', endpoint, parameters)
  }

  /**
   * `post` perform post request to integration service
   */

  post = <DataReturned = any>(endpoint: string, parameters?: BearerRequestParameters) => {
    return this.request<DataReturned>('POST', endpoint, parameters)
  }

  /**
   * `put` perform put request to integration service
   */

  put = <DataReturned = any>(endpoint: string, parameters?: BearerRequestParameters) => {
    return this.request<DataReturned>('PUT', endpoint, parameters)
  }

  /**
   * `delete` perform delete request to integration service
   */

  delete = <DataReturned = any>(endpoint: string, parameters?: BearerRequestParameters) => {
    return this.request<DataReturned>('DELETE', endpoint, parameters)
  }

  /**
   * `patch` perform patch request to integration service
   */

  patch = <DataReturned = any>(endpoint: string, parameters?: BearerRequestParameters) => {
    return this.request<DataReturned>('PATCH', endpoint, parameters)
  }

  private request = <TData = any>(
    method: BearerRequestMethod,
    endpoint: string,
    parameters: BearerRequestParameters = {}
  ) => {
    if (parameters && typeof parameters !== 'object') {
      throw new InvalidRequestOptions()
    }

    const headers: BearerHeaders = {
      'Bearer-Auth-Id': this.authId!,
      'Bearer-Setup-Id': this.setupId!
    }

    if (parameters && parameters.headers) {
      for (const key in parameters.headers) {
        headers[`Bearer-Proxy-${key}`] = parameters.headers[key]
      }
    }

    return axios.request<TData>({
      method,
      headers: cleanOptions(headers),
      baseURL: `${this.bearerInstance.config.integrationHost}/api/v4/functions/${this.integrationId}/bearer-proxy`,
      url: endpoint,
      params: { ...parameters.query, clientId: this.bearerInstance.clientId },
      data: parameters && parameters.body
    })
  }
}

export type TFetchBearerData<T = any> = { data: T; referenceId?: string }

type BearerHeaders = Record<string, string | number | undefined>
type BearerRequestMethod = AxiosRequestConfig['method']

interface BearerRequestParameters {
  headers?: BearerHeaders
  query?: Record<string, string | number>
  body?: any
}
export default IntegrationClient

/**
 * Errors handling
 */

class InvalidRequestOptions extends Error {
  constructor() {
    super(`Unable to trigger API request. Request parameters should be an object \
in the form "{ headers: { "Foo": "bar" }, body: "My body" }"`)
  }
}
