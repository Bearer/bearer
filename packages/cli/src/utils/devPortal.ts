import cliUx from 'cli-ux'
import axios, { AxiosInstance } from 'axios'

import BaseCommand from '../base-command'
import { promptToLogin } from '../actions/login'
import { LOGIN_CLIENT_ID } from './constants'

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

export async function withFreshToken(command: BaseCommand) {
  const { expires_at, refresh_token } = (await command.bearerConfig.getToken()) || {
    expires_at: null,
    refresh_token: null
  }

  if (expires_at && refresh_token) {
    try {
      if (expires_at < Date.now()) {
        cliUx.action.start('Refreshing token')
        await refreshMyToken(command, refresh_token)
        cliUx.action.stop()
      }
    } catch (error) {
      cliUx.action.stop(`Failed`)
      command.error(error.message)
    }
  } else {
    await promptToLogin(command)
  }
}

// tslint:disable-next-line variable-name
async function refreshMyToken(command: BaseCommand, refresh_token: string): Promise<boolean | Error> {
  // TODO: rework refresh mechanism
  const response = await axios.post(`${command.constants.LoginDomain}/oauth/token`, {
    refresh_token,
    grant_type: 'refresh_token',
    client_id: LOGIN_CLIENT_ID
  })
  await command.bearerConfig.storeToken({ ...response.data, refresh_token })
  return true
}
