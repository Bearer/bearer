import debug, { BearerFunctionLogger } from '../src/logger'
import { IDebug } from 'debug'

describe('logger', () => {
  it('needs tests', () => {
    expect(debug).toBeInstanceOf(Function)
  })
})

describe('function logger', () => {
  it('returns an func logger', () => {
    const logger = new BearerFunctionLogger({})
    expect(logger).toBeInstanceOf(BearerFunctionLogger)
    expect(logger.log).toBeInstanceOf(Function)
  })
})
