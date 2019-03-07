import axios, { AxiosInstance } from 'axios'
class FetchDataError extends Error {}
class UpdateDataError extends Error {}
class CreateDataError extends Error {}

type TPersistedData = {
  Item: { referenceId: string; [key: string]: any }
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

  async getData(referenceId: string): Promise<TPersistedData | null> {
    if (!referenceId) {
      return Promise.resolve(null)
    }
    try {
      const data = await this.client.get(`api/v2/items/${referenceId}`, { params: { signature: this.signature } })
      return data.data
    } catch (error) {
      if (error.response && !(error.response.status === 404)) {
        throw new FetchDataError('Error while retrieving data')
      }
    }
    return Promise.resolve(null)
  }

  async updateData(referenceId, data): Promise<TPersistedData> {
    try {
      const response = await this.client.put(
        `api/v2/items/${referenceId}`,
        { ...data, ReadAllowed: true },
        { params: { signature: this.signature } }
      )
      return response.data
    } catch (error) {
      throw new UpdateDataError(`Error while updating data: ${error.toString()}`)
    }
  }

  async saveData(data): Promise<TPersistedData> {
    try {
      const response = await this.client.post(
        `api/v2/items`,
        { ...data, ReadAllowed: true },
        { params: { signature: this.signature } }
      )
      return response.data
    } catch (error) {
      throw new CreateDataError(`Error while creating data: ${error.toString()}`)
    }
  }
}

export default (baseURL: string, signature: string): DBClient => new DBClient(baseURL, signature)
