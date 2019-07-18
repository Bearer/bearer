import exportFromIndex from './index'
import { middleware } from './middleware'

describe('index', () => {
  it('exports middleware', () => {
    expect(exportFromIndex).toEqual(middleware)
  })
})
