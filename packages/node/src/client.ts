import axios, { AxiosRequestConfig, AxiosInstance } from 'axios'

export class BearerClient {
  protected readonly bearerApiKey: string
  protected options: BearerClientOptions = { host: 'https://int.bearer.sh' }

  constructor(bearerApiKey: string, options?: BearerClientOptions) {
    this.options = { ...this.options, ...options }
    this.bearerApiKey = bearerApiKey
  }

  public integration(integrationId: string) {
    return new BearerClientInstance(integrationId, this.options, this.bearerApiKey)
  }

  /**
   * Deprecated. Please use integration(...).invoke(...) instead
   */

  public invoke(
    integrationId: string,
    functionName: string,
    { query, body }: { query?: any; body?: any } = { query: {}, body: {} }
  ) {
    const instance = new BearerClientInstance(integrationId, this.options, this.bearerApiKey)
    return instance.invoke(functionName, { query, body })
  }
}

export class BearerClientInstance {
  protected readonly bearerApiKey: string

  protected integrationId: string
  protected setupId?: string
  protected authId?: string

  protected client: AxiosInstance
  protected clientOptions: BearerClientOptions

  constructor(integrationId: string, options: BearerClientOptions, bearerApiKey: string) {
    this.integrationId = integrationId
    this.bearerApiKey = bearerApiKey
    this.clientOptions = options

    this.client = axios
  }

  public auth = ({ authId, setupId }: { authId?: string; setupId?: string }) => {
    this.setupId = setupId || ''
    this.authId = authId || ''

    return this
  }

  public authenticate = this.auth // Alias

  /**
   * HTTP methods
   */

  public get = (endpoint: string, parameters?: BearerRequestParameters, options?: any) => {
    return this.makeRequest('GET', endpoint, parameters, options)
  }

  public head = (endpoint: string, parameters?: BearerRequestParameters, options?: any) => {
    return this.makeRequest('HEAD', endpoint, parameters, options)
  }

  public post = (endpoint: string, parameters?: BearerRequestParameters, options?: any) => {
    return this.makeRequest('POST', endpoint, parameters, options)
  }

  public put = (endpoint: string, parameters?: BearerRequestParameters, options?: any) => {
    return this.makeRequest('PUT', endpoint, parameters, options)
  }

  public delete = (endpoint: string, parameters?: BearerRequestParameters, options?: any) => {
    return this.makeRequest('DELETE', endpoint, parameters, options)
  }

  public patch = (endpoint: string, parameters?: BearerRequestParameters, options?: any) => {
    return this.makeRequest('PATCH', endpoint, parameters, options)
  }

  protected makeRequest = (
    method: BearerRequestMethod,
    endpoint: string,
    parameters?: BearerRequestParameters,
    options?: BearerRequestOptions
  ) => {
    if (parameters && typeof parameters !== 'object') {
      throw new InvalidRequestOptions()
    }

    const preheaders: { [key: string]: string } = {
      Authorization: this.bearerApiKey,
      'User-Agent': 'Bearer.sh',
      'Bearer-Auth-Id': this.authId!,
      'Bearer-Setup-Id': this.setupId!
    }

    if (parameters && parameters.headers) {
      for (const key in parameters.headers) {
        preheaders[`Bearer-Proxy-${key}`] = parameters.headers[key]
      }
    }

    const headers = Object.keys(preheaders).reduce(
      (acc, key) => {
        const header = preheaders[key]

        if (header !== undefined && header !== null) {
          acc[key] = preheaders[key]
        }

        return acc
      },
      {} as any
    )

    return this.client.request({
      method,
      headers,
      baseURL: `${this.clientOptions.host}/api/v4/functions/backend/${this.integrationId}/bearer-proxy`,
      url: endpoint,
      params: parameters && parameters.query,
      data: parameters && parameters.body
    })
  }

  /**
   * Invoke custom functions
   */

  public invoke = (functionName: string, { query, body }: { query?: any; body?: any } = { query: {}, body: {} }) => {
    return this.client.request({
      baseURL: `${this.clientOptions.host}/api/v4/functions/backend/${this.integrationId}`,
      url: `/${functionName}`,
      headers: {
        Authorization: this.bearerApiKey
      },
      method: 'post',
      data: body,
      params: query
    })
  }
}

/**
 * Export
 */

export default (apiKey: string | undefined): BearerClient => {
  if (!apiKey) {
    throw new InvalidAPIKey(apiKey)
  }

  return new BearerClient(apiKey)
}

/**
 * Types
 */

interface BearerRequestParameters {
  headers?: { [key: string]: string }
  query?: string | { [key: string]: string }
  body?: any
}

type BearerRequestMethod = AxiosRequestConfig['method']
type BearerRequestOptions = any

type BearerClientOptions = { [key: string]: string }

/**
 * Errors handling
 */

class InvalidAPIKey extends Error {
  constructor(token: any) {
    super(`Invalid Bearer API key provided.  Value: ${token} \
You'll find you API key at this location: https://app.bearer.sh/keys`)
  }
}

class InvalidRequestOptions extends Error {
  constructor() {
    super(`Unable to trigger API request. Request parameters should be an object \
in the form "{ headers: { "Foo": "bar" }, body: "My body" }"`)
  }
}
