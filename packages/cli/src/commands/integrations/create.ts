import { flags } from '@oclif/command'
// @ts-ignore
import * as suggest from 'inquirer-prompt-suggest'

import BaseCommand from '../../base-command'
import { linkIntegration } from '../../utils/commands'
import { ensureFreshToken } from '../../utils/decorators'
import { askForString } from '../../utils/prompts'
import * as inquirer from 'inquirer'

export default class IntegrationsCreate extends BaseCommand {
  static description = 'create a new Integration'

  static flags = {
    ...BaseCommand.flags,
    description: flags.string({ char: 'd' }),
    name: flags.string({ char: 'n' }),
    skipLink: flags.boolean({ char: 'l' })
  }

  @ensureFreshToken()
  async run() {
    const { flags } = this.parse(IntegrationsCreate)
    const name = flags.name || (await this.askForName())
    const description = flags.description || (await askForString('Description (optional)'))

    try {
      const { data } = await this.devPortalClient.request<CreateIntegration>({
        query: MUTATION,
        variables: { name, description }
      })
      if (data.data) {
        const { integration } = data.data.createIntegration
        this.success('Integration successfully created')
        this.log(
          '      name: %s\n      uuid: %s\nidentifier: %s\n       Url:',
          integration.name,
          integration.uuid,
          `${this.constants.DeveloperPortalUrl}integrations/${integration.uuid}`
        )
        if (this.isIntegrationLocation) {
          // tslint:disable-next-line no-boolean-literal-compare
          if (!flags.skipLink) {
            await linkIntegration.bind(this)(integration.uuid)
          }
        }
      } else {
        this.debug(data)
        this.error('Unable to create this integration, please retry')
      }
    } catch (e) {
      this.debug('%j', e.response)
      this.error(e)
    }
  }

  async askForName(): Promise<string> {
    inquirer.registerPrompt('suggest', suggest)
    const suggestions = Array.from(Array(30).keys()).map(randomName)
    this.debug('%j', suggestions)
    return askForString('Integration name', {
      suggestions,
      validate: (input: string) => {
        return input.length > 0
      },
      type: 'suggest'
    } as any)
  }
}

type Integration = {
  uuid: string
  name: string
  latestActivity?: {
    state: string
  }
}

function randomName() {
  return [bearSample(), nameSample()].join(' ')
}

function bearSample() {
  return bears[Math.floor(Math.random() * bears.length)]
}

function nameSample() {
  return names[Math.floor(Math.random() * names.length)]
}

const bears = [
  'Cinnamon',
  'Florida',
  'Glacier',
  'Haida Gwaii',
  'Kermode',
  'Spirit',
  'Louisiana',
  'Newfoundland',
  'Baluchistan',
  'Formosan',
  'Himalayan',
  'Ussuri',
  'Alaska Peninsula',
  'Atlas',
  'Bergman',
  'Cantabrian',
  'Gobi',
  'Grizzly',
  'Kamchatka',
  'Kodiak',
  'Marsican',
  'Sitka',
  'Stickeen',
  'Ussuri',
  'Giant',
  'Qinling',
  'Sloth',
  'Sun',
  'Polar',
  'Ursid hybrid',
  'Spectacled'
]
const names = [
  'Alpha',
  'Bravo',
  'Charlie',
  'Delta',
  'Echo',
  'Foxtrot',
  'Golf',
  'Hotel',
  'India',
  'Juliet',
  'Kilo',
  'Lima',
  'Mike',
  'November',
  'Oscar',
  'Papa',
  'Quebec',
  'Romeo',
  'Sierra',
  'Tango',
  'Uniform',
  'Victor',
  'Whiskey',
  'X-ray',
  'Yankee',
  'Zulu'
]

type CreateIntegration = { createIntegration: { integration: Integration } }

const MUTATION = `
mutation CLICreateIntegration($name: String!, $description: String!) {
  createIntegration(name: $name, description: $description) {
    integration {
      uuid: uuidv2
      name
    }
  }
} 
`
