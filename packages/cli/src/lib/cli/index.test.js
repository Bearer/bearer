const { CLI } = require('./index')
const program = require('commander')
let command = {
  useWith: program => {
    program.command('goats', 'Show me a great animals')
  }
}

test('using a command', () => {
  const cli = new CLI(program)
  cli.use(command)
  expect(cli.program.commands[0]._name).toBe('goats')
})
