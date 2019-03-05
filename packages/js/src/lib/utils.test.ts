import { formatQuery, cleanQuery } from './utils'

describe('formatQuery', () => {
  it('filters empty params and retunrs a string', () => {
    const params = {
      aNullParams: null,
      undefinedParams: undefined,
      falseParams: false,
      aString: 'ok',
      aNumber: 1
    }
    expect(formatQuery(params)).toEqual('aString=ok&aNumber=1')
  })

  it('returns and empty string', () => {
    const params = {}
    expect(formatQuery(params)).toEqual('')
  })
})

describe('cleanQuery', () => {
  it('filters empty params and retunrs a string', () => {
    const params = {
      aNullParams: null,
      undefinedParams: undefined,
      falseParams: false,
      aString: 'ok',
      aNumber: 1
    }
    expect(cleanQuery(params)).toEqual({ aString: 'ok', aNumber: 1 })
  })
})
