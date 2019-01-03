import { CLI } from './index'
import * as program from 'commander'

const command = {
  useWith: program => {
    program.command('goats', 'Show me a great animals')
  }
}

test('using a command', () => {
  const cli = new CLI(program, null, {
    scenarioConfig: { config: '/tmp/scenariorc' }
  } as any)
  cli.use(command)

  expect(cli.program.commands[0]._name).toBe('goats')
})
