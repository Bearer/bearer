import axios, { AxiosInstance } from 'axios'

import BaseCommand from '../base-command'

export class IntegrationClient {
  private client: AxiosInstance

  constructor(baseUrl: string, authorization?: string, version?: string) {
    const headers = {
      Authorization: authorization,
      ['BEARER-CLI-VERSION']: version
    }
    this.client = axios.create({
      headers,
      baseURL: baseUrl
    })
  }

  async getIntegrationArchiveUploadUrl(orgId: string, integrationId: string): Promise<string> {
    try {
      const response = await this.client.post('scenario-archive-url', { orgId, scenarioId: integrationId, type: 'src' })
      return response.data.url
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
