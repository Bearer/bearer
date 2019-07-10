import axios, { AxiosRequestConfig, AxiosInstance } from 'axios'

class Bearer {
  protected readonly bearerApiKey: string
  protected options: BearerClientOptions = { host: 'https://int.bearer.sh' }

  constructor(bearerApiKey: string, options?: BearerClientOptions) {
    this.options = { ...this.options, ...options }
    this.bearerApiKey = bearerApiKey
  }

  public integration(integrationId: string) {
    return new BearerClient(integrationId, this.options, this.bearerApiKey)
  }

  /**
   * Deprecated. Please use integration(...).invoke(...) instead
   */

  public invoke(
    integrationId: string,
    functionName: string,
    { query, body }: { query?: any; body?: any } = { query: {}, body: {} }
  ) {
    const instance = new BearerClient(integrationId, this.options, this.bearerApiKey)
    return instance.invoke(functionName, { query, body })
  }
}

class BearerClient {
  protected readonly client: AxiosInstance = axios

  constructor(
    readonly integrationId: string,
    readonly options: BearerClientOptions,
    readonly bearerApiKey: string,
    readonly setupId?: string,
    readonly authId?: string
  ) {}

  public auth = (authId: string) => {
    return new BearerClient(this.integrationId, this.options, this.bearerApiKey, this.setupId, authId)
  }

  public setup = (setupId: string) => {
    return new BearerClient(this.integrationId, this.options, this.bearerApiKey, setupId, this.authId)
  }

  public authenticate = this.auth // Alias

  /**
   * HTTP methods
   */

  public get = (endpoint: string, parameters?: BearerRequestParameters, options?: any) => {
    return this.request('GET', endpoint, parameters, options)
  }

  public head = (endpoint: string, parameters?: BearerRequestParameters, options?: any) => {
    return this.request('HEAD', endpoint, parameters, options)
  }

  public post = (endpoint: string, parameters?: BearerRequestParameters, options?: any) => {
    return this.request('POST', endpoint, parameters, options)
  }

  public put = (endpoint: string, parameters?: BearerRequestParameters, options?: any) => {
    return this.request('PUT', endpoint, parameters, options)
  }

  public delete = (endpoint: string, parameters?: BearerRequestParameters, options?: any) => {
    return this.request('DELETE', endpoint, parameters, options)
  }

  public patch = (endpoint: string, parameters?: BearerRequestParameters, options?: any) => {
    return this.request('PATCH', endpoint, parameters, options)
  }

  public request = (
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
      baseURL: `${this.options.host}/api/v4/functions/backend/${this.integrationId}/bearer-proxy`,
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
      baseURL: `${this.options.host}/api/v4/functions/backend/${this.integrationId}`,
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
 * Types
 */

type BearerRequestMethod = AxiosRequestConfig['method']
interface BearerRequestParameters {
  headers?: { [key: string]: string }
  query?: string | { [key: string]: string }
  body?: any
}

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

/**
 * Exports
 */

export default (apiKey: string | undefined, options?: BearerClientOptions): Bearer => {
  if (!apiKey) {
    throw new InvalidAPIKey(apiKey)
  }

  return new Bearer(apiKey, options)
}

export { Bearer }
