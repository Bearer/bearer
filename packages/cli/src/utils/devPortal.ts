import axios, { AxiosInstance } from 'axios'
import BaseCommand from '../base-command'
import { withFreshToken } from './decorators'

export class Client {
  private instance: AxiosInstance
  constructor(private readonly command: BaseCommand) {
    this.instance = axios.create({
      baseURL: command.constants.DeveloperPortalAPIUrl
    })
    this.instance.interceptors.response.use(
      r => r,
      error => {
        if (error.response && error.response.status) {
          command.error('Unauthorized action, please run bearer login first')
        }
        return Promise.reject(error)
      }
    )
  }

  async request<DataReturned>(data: { query: string; variables?: any }) {
    await withFreshToken(this.command)

    const token = await this.command.bearerConfig.getToken()

    return await this.instance.post<{ data?: DataReturned; errors?: any }>('', data, {
      headers: {
        Authorization: `Bearer ${token!.access_token}`
      }
    })
  }
}

export const devPortalClient = (command: BaseCommand) => {
  return new Client(command)
}
