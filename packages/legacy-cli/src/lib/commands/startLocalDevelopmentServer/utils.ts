import { getRows } from './storage'

export async function loadUserDefinedData({ query }): Promise<any> {
  const userDataIds = Object.keys(query).filter(key => key.endsWith('Id') && query[key])

  if (userDataIds.length > 0) {
    return (await Promise.all(
      userDataIds.map(async id => {
        try {
          return await getRows(query[id])
        } catch (e) {
          console.log('[BEARER]', `Error while loading ${id} associated data (value: ${query[id] || 'No Value'})`)
          return null
        }
      })
    )).reduce((acc, datum, index) => {
      if (datum) {
        const parsed = JSON.parse(datum)
        const key = userDataIds[index]
        console.log('[BEARER]', `${itemName(key)} loaded :\n`, JSON.stringify(datum, null, 2))
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
