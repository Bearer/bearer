import axios, { AxiosInstance } from 'axios'

type TClientOptions = {
  hostUrl: string
}

export class BearerClient {
  private static defaultOptions = { hostUrl: 'https://int.bearer.sh/backend/api/v1/' }
  private options: TClientOptions
  private client: AxiosInstance

  constructor(private readonly token: string, clientOptions: Partial<TClientOptions> = {}) {
    this.options = { ...BearerClient.defaultOptions, ...clientOptions }

    this.client = axios.create({
      baseURL: this.options.hostUrl,
      headers: {
        Authorization: token
      }
    })
  }

  public call(
    scenarioName: string,
    intentName: string,
    { query, body }: { query: any; body: any } = { query: {}, body: {} }
  ) {
    return this.client.post(`${scenarioName}/${intentName}`, body, {
      params: query
    })
  }
}

export default (token: string): BearerClient => {
  return new BearerClient(token)
}
