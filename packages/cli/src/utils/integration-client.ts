import axios, { AxiosInstance } from 'axios'

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
