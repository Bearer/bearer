const CLI = require('../src/lib/cli').CLI

const deployCmd = require('../src/lib/commands/deployCommand')
const initCmd = require('../src/lib/commands/initCommand')
const generateCmd = require('../src/lib/commands/generateCommand')

const program = require('commander')

const cli = new CLI(program, null, {
  HandlerBase: 'index.js',
  scenarioConfig: { config: '/tmp/scenariorc' }
})

cli.use(deployCmd)
cli.use(generateCmd)
cli.use(initCmd)

describe('deploy command', () => {
  test('program have `deploy` command regirstered', () => {
    expect(program.commands.map(cmd => cmd._name)).toContain('deploy')
  })
})

describe('new command', () => {
  test('program have `new` command regirstered', () => {
    expect(program.commands.map(cmd => cmd._name)).toContain('new')
  })
})

describe('generate command', () => {
  test('program have `generate` command regirstered', () => {
    const command = program.commands.find(
      command => command._name === 'generate'
    )

    expect(command.options).toEqual(expect.arrayContaining([]))
    expect(program.commands.map(cmd => cmd._name)).toContain('generate')
  })
})
