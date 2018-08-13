import * as intents from '@bearer/intents'

import { getComponentVars, getIntentVars } from '../commands/generateCommand'
const TestingValues = ['spongebob', 'SpongeBob', 'spongeBob', 'sponge_bob', 'sponge-bob']
const AuthTypes = ['oauth2', 'apiKey', 'noAuth', 'basicAuth']

describe('Generate command', () => {
  describe('get components variables', () => {
    TestingValues.forEach(value => {
      AuthTypes.forEach(authType => {
        describe(authType, () => {
          describe(value, () => {
            it('formats variables correctly', () => {
              expect(getComponentVars(value, { authType })).toMatchSnapshot()
            })
          })
        })
      })
    })
  })

  describe('get intents variables', () => {
    Object.keys(intents).forEach(intentType => {
      AuthTypes.forEach(authType => {
        describe(authType, () => {
          describe(intentType, () => {
            TestingValues.forEach(value => {
              describe(value, () => {
                it('formats variables correctly', () => {
                  expect(getIntentVars(value, intentType, { authType })).toMatchSnapshot()
                })
              })
            })
          })
        })
      })
    })
  })
})
