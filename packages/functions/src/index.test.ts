import { DBClient } from './db-client'
import { FetchData, DBClient as CLIENT } from './index'

describe('index', () => {
  it('export FetchData', () => {
    expect(FetchData).toBeTruthy()
  })

  it('export a dbclient instance', () => {
    expect(CLIENT).toEqual(DBClient.instance)
  })
})
