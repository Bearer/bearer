import axios, { AxiosInstance } from 'axios'
class FetchDataError extends Error {}
class UpdateDataError extends Error {}
class CreateDataError extends Error {}

type TPersistedData = {
  Item: { referenceId: string; [key: string]: any }
}
class UserDataClient {
  private client: AxiosInstance
  constructor(private readonly baseURL: string) {
    console.log('[BEARER]', 'baseURL', baseURL)
    this.client = axios.create({
      baseURL,
      timeout: 3000,
      headers: {
        Accept: 'application/json',
        'User-Agent': 'Bearer'
      }
    })
  }

  async retrieveState(referenceId: string): Promise<TPersistedData> {
    if (!referenceId) {
      return Promise.resolve(null)
    }
    try {
      const data = await this.client.get(`api/v1/items/${referenceId}`)
      return data.data
    } catch (error) {
      if (error.response && !(error.response.status === 404)) {
        throw new FetchDataError('Error while retrieving data')
      }
    }
    return Promise.resolve(null)
  }

  async updateState(referenceId, state): Promise<TPersistedData> {
    try {
      const response = await this.client.put(`api/v1/items/${referenceId}`, { ...state, ReadAllowed: true })
      return response.data
    } catch (error) {
      throw new UpdateDataError('Error while updating data')
    }
  }

  async saveState(state): Promise<TPersistedData> {
    try {
      const response = await this.client.post(`api/v1/items`, { ...state, ReadAllowed: true })
      return response.data
    } catch (error) {
      throw new CreateDataError(`Error while creating data: ${error.toString()}`)
    }
  }
}

export default (baseURL: string): UserDataClient => new UserDataClient(baseURL)
