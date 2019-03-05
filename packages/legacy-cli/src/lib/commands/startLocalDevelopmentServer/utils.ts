import { getRows } from './storage'
import debug from '../../logger'

const logger = debug.extend('utils')

export async function loadUserDefinedData({ query }): Promise<any> {
  const userDataIds = Object.keys(query).filter(key => key.endsWith('Id') && query[key])

  if (userDataIds.length > 0) {
    return (await Promise.all(
      userDataIds.map(async id => {
        try {
          return await getRows(query[id])
        } catch (e) {
          logger('Error while loading %s associated data (value: %s)', id, query[id] || 'No Value')
          return null
        }
      })
    )).reduce((acc, datum, index) => {
      if (datum) {
        const parsed = JSON.parse(datum)
        const key = userDataIds[index]
        logger('%s loaded: \n%j', itemName(key), datum)
        acc[itemName(key)] = parsed.ReadAllowed ? parsed : null
      } else {
      }
      return acc
    }, {})
  }
  return {}
}

function itemName(key: string): string {
  return key.replace(/Id$/, '')
}
