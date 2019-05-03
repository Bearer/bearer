import runCreate from './run'
import * as cli from '@bearer/cli/lib/index'

// @ts-ignore
cli.run = jest.fn(() => Promise.resolve())

describe('run', () => {
  describe('when -t options is passed', () => {
    beforeEach(() => {
      // @ts-ignore
      cli.run.mockReset()
    })

    describe('when empty', () => {
      it('fallbacks on bearer template url', async () => {
        await runCreate('aName', '-t', '--authType', 'ok')
        expect(cli.run).toHaveBeenCalledWith([
          'new',
          'aName',
          '-t',
          'https://github.com/Bearer/templates',
          '--authType',
          'ok'
        ])
      })
    })

    describe('when empty and latest', () => {
      it('fallbacks on bearer template url', async () => {
        await runCreate('aName', '-t')
        expect(cli.run).toHaveBeenCalledWith(['new', 'aName', '-t', 'https://github.com/Bearer/templates'])
      })
    })

    describe('when present', () => {
      it('does not change', async () => {
        await runCreate('aName', '-t', 'anotherUrl', '--authType', 'ok')
        expect(cli.run).toHaveBeenCalledWith(['new', 'aName', '-t', 'anotherUrl', '--authType', 'ok'])
      })
    })
  })
})
