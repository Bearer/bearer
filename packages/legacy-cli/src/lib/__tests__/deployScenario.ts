import * as fs from 'fs-extra'

import * as deployScenario from '../deployScenario'
import { Config } from '../types'
test('deployIntents is defined', () => {
  expect(deployScenario.deployIntents).toBeDefined()
})

let emit = jest.fn((...args) => console.log(...args))

beforeEach(() => {
  fs.ensureDirSync('./tmp/views')
  fs.ensureDirSync('./tmp/intents')
})

afterEach(() => {
  fs.removeSync('./tmp')
})

describe('deployViews', () => {
  test.skip(
    'Is not hanging in the end',
    async () => {
      expect.assertions(1)
      const locator = {} as any
      const config = {
        scenarioConfig: { scenarioTitle: 'test', orgId: '4l1c3' },
        rootPathRc: './tmp/.test'
      } as Config
      await expect(deployScenario.deployViews({ emit }, config, locator)).resolves.toEqual({})
    },
    1000
  )
})
