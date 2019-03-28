import axios, { AxiosInstance } from 'axios'
class FetchDataError extends Error {}
class SaveDataError extends Error {}

type TPersistedData<T extends {} = { [key: string]: any }> = {
  Item: { referenceId: string; data?: T; ReadAllowed: boolean }
}

export class DBClient {
  static instance(signature) {
    return new DBClient(process.env.bearerBaseURL!, signature)
  }

  private client: AxiosInstance

  constructor(private readonly baseURL: string, private readonly signature: string) {
    this.client = axios.create({
      baseURL: this.baseURL,
      timeout: 3000,
      headers: {
        Accept: 'application/json',
        'User-Agent': 'Bearer'
      }
    })
  }

  async getData<SavedData = {}>(referenceId: string): Promise<TPersistedData<SavedData>> {
    if (!referenceId) {
      return Promise.resolve({ Item: { referenceId, ReadAllowed: false } })
    }
    try {
      const data = await this.client.get<TPersistedData<SavedData>>(`api/v2/items/${referenceId}`, {
        params: { signature: this.signature }
      })
      return data.data
    } catch (error) {
      if (error.response && !(error.response.status === 404)) {
        throw new FetchDataError('Error while retrieving data')
      }
    }
    return Promise.resolve({ Item: { referenceId, ReadAllowed: false } })
  }

  async upsertData<InputData = {}>(referenceId, data): Promise<TPersistedData<InputData>> {
    if (!referenceId) {
      throw new SaveDataError(`Error while updating data: no reference given`)
    }
    try {
      await this.client.put<TPersistedData<InputData>>(
        `api/v2/items/${referenceId}`,
        { data, referenceId, ReadAllowed: true },
        { params: { signature: this.signature } }
      )
      return { Item: { data, referenceId, ReadAllowed: true } }
    } catch (error) {
      throw new SaveDataError(`Error while updating data: ${error.toString()}`)
    }
  }
}

export default (baseURL: string, signature: string): DBClient => new DBClient(baseURL, signature)
