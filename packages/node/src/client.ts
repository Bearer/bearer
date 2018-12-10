import axios, { AxiosInstance } from 'axios'

type TClientOptions = {
  hostUrl: string
}

type TIntentParams = { query?: any; body?: any }
const defaultIntentParams = { query: {}, body: {} }

export class BearerClient {
  protected static defaultOptions = { hostUrl: 'https://int.bearer.sh/backend/api/v1/' }
  protected options: TClientOptions
  protected client: AxiosInstance

  constructor(protected readonly token: string, clientOptions: Partial<TClientOptions> = {}) {
    this.options = { ...BearerClient.defaultOptions, ...clientOptions }

    this.client = axios.create({
      baseURL: this.options.hostUrl,
      headers: {
        Authorization: token
      }
    })
  }

  public call(scenarioName: string, intentName: string, { query, body }: TIntentParams = defaultIntentParams) {
    return this.client.post(`${scenarioName}/${intentName}`, body, {
      params: query
    })
  }
}

export class ScenarioClient {
  private bearerClient: BearerClient

  constructor(token: string, clientOptions: Partial<TClientOptions> = {}, private readonly scenarioName: string) {
    this.bearerClient = new BearerClient(token, clientOptions)
  }

  public call(intentName: string, intentParams: TIntentParams = defaultIntentParams) {
    return this.bearerClient.call(this.scenarioName, intentName, intentParams)
  }
}

export default (token: string): BearerClient => {
  return new BearerClient(token)
}
