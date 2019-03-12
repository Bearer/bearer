import axios, { AxiosInstance } from 'axios'

import BaseCommand from '../base-command'

export class IntegrationClient {
  private client: AxiosInstance

  constructor(baseURL: string, authorization?: string, version?: string) {
    const headers = {
      Authorization: authorization,
      ['BEARER-CLI-VERSION']: version
    }
    this.client = axios.create({
      headers,
      baseURL
    })
  }

  async getIntegrationArchiveUploadUrl(orgId: string, integrationId: string): Promise<string> {
    try {
      const { data } = await this.client.post('integration-urls', { orgId, integrationId })
      return data.url
    } catch (e) {
      throw e
    }
  }
}

export default (command: BaseCommand): IntegrationClient => {
  return new IntegrationClient(
    command.bearerConfig.DeploymentUrl,
    command.bearerConfig.bearerConfig.authorization.AuthenticationResult!.IdToken,
    command.config.version
  )
}
