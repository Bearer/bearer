const CLI = require('../src/lib/cli').CLI

const deployCmd = require('../src/lib/commands/deployCommand')

const program = require('commander')

const cli = new CLI(program, null, {
  scenarioConfig: { config: '/tmp/scenariorc' }
})

cli.use(deployCmd)

describe('deploy command', () => {
  test('program have `deploy` command regirstered', () => {
    expect(program.commands.map(cmd => cmd._name)).toContain('deploy')
  })
})
