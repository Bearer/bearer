import debug, { BearerIntentLogger } from '../src/logger'
import { IDebug } from 'debug'

describe('logger', () => {
  it('needs tests', () => {
    expect(debug).toBeInstanceOf(Function)
  })
})

describe('intent logger', () => {
  it('returns an intent logger', () => {
    const logger = new BearerIntentLogger({})
    expect(logger).toBeInstanceOf(BearerIntentLogger)
    expect(logger.log).toBeInstanceOf(Function)
  })
})
