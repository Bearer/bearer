import { cleanQuery, buildQuery, cleanOptions } from '../../lib/utils'

describe('cleanQuery', () => {
  it('filters empty params and returns a string', () => {
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

describe('cleanOptions', () => {
  it('removes keys with undefined valued', () => {
    cleanOptions
    const params = {
      aNullParams: null,
      undefinedParams: undefined,
      falseParams: false,
      aString: 'ok',
      aNumber: 1
    }

    expect(cleanOptions(params)).toEqual({ aNullParams: null, falseParams: false, aString: 'ok', aNumber: 1 })
  })
})

describe('buildQuery', () => {
  it('creates a valid query string', () => {
    cleanOptions
    const params = {
      aNullParams: null,
      undefinedParams: undefined,
      falseParams: false,
      aString: 'ok with space and accents ééà',
      aNumber: 1
    }

    expect(buildQuery(params)).toEqual(
      'aNullParams=null&undefinedParams=undefined&falseParams=false&aString=ok%20with%20space%20and%20accents%20%C3%A9%C3%A9%C3%A0&aNumber=1'
    )
  })

  it('returns and empty string', () => {
    const params = {}

    expect(buildQuery(params)).toEqual('')
  })
})
