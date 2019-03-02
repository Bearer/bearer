import * as BearerReact from './index'

import BearerConnect from './Connect'
import BearerProvider from './bearer-provider'
import bearerFactory from './factory'

describe('main file', () => {
  describe('exports', () => {
    it('exposes Context Provider', () => {
      expect(BearerReact.Bearer).toBe(BearerProvider)
    })

    it('exposes Connect', () => {
      expect(BearerReact.Connect).toBe(BearerConnect)
    })

    it('exposes Factory', () => {
      expect(BearerReact.factory).toBe(bearerFactory)
    })
  })
})
