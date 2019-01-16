import { DBClient } from './db-client'
import { RetrieveState, SaveState, FetchData, DBClient as CLIENT } from './index'

describe('index', () => {
  it('export SaveState', () => {
    expect(SaveState).toBeTruthy()
  })

  it('export RetrieveState', () => {
    expect(RetrieveState).toBeTruthy()
  })

  it('export FetchData', () => {
    expect(FetchData).toBeTruthy()
  })

  it('export a dbclient instance', () => {
    expect(CLIENT).toEqual(DBClient.instance)
  })
})
