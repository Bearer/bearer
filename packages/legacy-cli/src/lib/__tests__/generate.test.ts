import * as intents from '@bearer/intents'
import Authentications from '@bearer/types/lib/Authentications'

import { getComponentVars, getIntentVars } from '../commands/generateCommand'
const TestingValues = ['spongebob', 'SpongeBob', 'spongeBob', 'sponge_bob', 'sponge-bob']

describe('Generate command', () => {
  describe('get components variables', () => {
    TestingValues.forEach(value => {
      Object.keys(Authentications).forEach(authType => {
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
    Object.keys(intents)
      .filter(i => i !== 'DBClient')
      .forEach(intentType => {
        Object.keys(Authentications).forEach(authType => {
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
