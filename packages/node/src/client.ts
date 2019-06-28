import axios, { AxiosInstance } from 'axios'

type TClientOptions = {
  hostUrl: string
}

type TFunctionParams = { query?: any; body?: any }
const defaultFunctionParams = { query: {}, body: {} }

export class BearerClient<T = string> {
  protected static defaultOptions = { hostUrl: 'https://int.bearer.sh/api/v4/functions/backend' }
  protected options: TClientOptions
  protected client: AxiosInstance

  constructor(protected readonly token: string, clientOptions: Partial<TClientOptions> = {}) {
    if (!token) {
      throw new InvalidAPIKey(token)
    }
    this.options = { ...BearerClient.defaultOptions, ...clientOptions }

    this.client = axios.create({
      baseURL: this.options.hostUrl,
      headers: {
        Authorization: token
      }
    })
  }

  public invoke(integrationName: string, functionName: T, { query, body }: TFunctionParams = defaultFunctionParams) {
    return this.client.post(`${integrationName}/${functionName}`, body, {
      params: query
    })
  }
}

export class IntegrationClient<T = string> {
  private bearerClient: BearerClient<T>

  constructor(token: string, clientOptions: Partial<TClientOptions> = {}, private readonly integrationName: string) {
    this.bearerClient = new BearerClient<T>(token, clientOptions)
  }

  public invoke(functionName: T, functionParams: TFunctionParams = defaultFunctionParams) {
    return this.bearerClient.invoke(this.integrationName, functionName, functionParams)
  }
}

export default (token: string): BearerClient => {
  return new BearerClient(token)
}

class InvalidAPIKey extends Error {
  constructor(token: any) {
    super(`Invalid Bearer API key provided.  Value: ${token}
You'll find you API key at this location: https://app.bearer.sh/keys`)
  }
}
