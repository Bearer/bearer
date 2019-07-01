import axios, { AxiosRequestConfig, AxiosInstance } from 'axios'

export class BearerClient {
  protected readonly BEARER_API_KEY: string
  protected options: BearerClientOptions = { host: 'https://int.bearer.sh' }

  constructor(bearerApiKey: string | undefined, options?: BearerClientOptions) {
    if (!bearerApiKey) {
      throw new InvalidAPIKey(bearerApiKey)
    }

    this.options = { ...this.options, ...options }
    this.BEARER_API_KEY = bearerApiKey
  }

  public integration(integrationId: string) {
    return new BearerClientInstance(integrationId, this.options, this.BEARER_API_KEY)
  }

  /**
   * Deprecated. Please use integration(...).invoke(...) instead
   */

  public invoke(
    integrationId: string,
    functionName: string,
    { query, body }: { query?: any; body?: any } = { query: {}, body: {} }
  ) {
    const instance = new BearerClientInstance(integrationId, this.options, this.BEARER_API_KEY)
    return instance.invoke(functionName, { query, body })
  }
}

export class BearerClientInstance {
  protected readonly BEARER_API_KEY: string

  protected INTEGRATION_ID: string
  protected SETUP_ID?: string
  protected AUTH_ID?: string

  protected client: AxiosInstance
  protected clientOptions: BearerClientOptions

  constructor(integrationId: string, options: BearerClientOptions, BEARER_API_KEY: string) {
    this.INTEGRATION_ID = integrationId
    this.BEARER_API_KEY = BEARER_API_KEY
    this.clientOptions = options

    this.client = axios
  }

  public auth = ({ authId, setupId }: { authId?: string; setupId?: string }) => {
    this.AUTH_ID = authId || ''
    this.SETUP_ID = setupId || ''

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

  public options = (endpoint: string, parameters?: BearerRequestParameters, options?: any) => {
    return this.makeRequest('OPTIONS', endpoint, parameters, options)
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
    if (!this.BEARER_API_KEY) {
      throw new InvalidAPIKey(this.BEARER_API_KEY)
    }
    if (parameters && typeof parameters !== 'object') {
      throw new InvalidRequestOptions()
    }

    const preheaders: { [key: string]: string } = {
      Authorization: this.BEARER_API_KEY,
      'User-Agent': 'Bearer.sh (nodejs)',
      'X-Bearer-Proxy-User-Agent': 'Bearer.sh',
      'X-Bearer-Auth-Id': this.AUTH_ID!,
      'X-Bearer-Setup-Id': this.SETUP_ID!
    }

    if (parameters && parameters.headers) {
      for (let key in parameters.headers) {
        preheaders[`X-Bearer-Proxy-${key}`] = parameters.headers[key]
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
      baseURL: `${this.clientOptions.host}/api/v4/functions/backend/${this.INTEGRATION_ID}/bearer-proxy`,
      url: endpoint,
      headers: headers,
      params: parameters && parameters.query,
      data: parameters && parameters.body
    })
  }

  /**
   * Invoke custom functions
   */

  public invoke = (functionName: string, { query, body }: { query?: any; body?: any } = { query: {}, body: {} }) => {
    if (!this.BEARER_API_KEY) {
      throw new InvalidAPIKey(this.BEARER_API_KEY)
    }

    return this.client.request({
      baseURL: `${this.clientOptions.host}/api/v4/functions/backend/${this.INTEGRATION_ID}`,
      url: `/${functionName}`,
      headers: {
        Authorization: this.BEARER_API_KEY
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
