/**
 * cleanQuery
 * @param params {object} remove all falsy values
 */
export function cleanQuery(params: Record<string, any>) {
  return Object.keys(params).reduce(
    (acc, key) => {
      if (params[key]) {
        acc[key] = params[key]
      }
      return acc
    },
    {} as Record<string, string>
  )
}

/**
 * cleanOptions remove all undefined keys
 * @param obj {object}
 */
export function cleanOptions(obj: Record<string, any>) {
  return Object.keys(obj).reduce(
    (acc, key: string) => {
      if (obj[key] !== undefined) {
        acc[key] = obj[key]
      }
      return acc
    },
    {} as Record<string, any>
  )
}

/**
 * buildQuery: transform an object to a valid query string
 * @param params {object}
 */

export function buildQuery(params: Record<string, any>) {
  function encode(k: string) {
    return encodeURIComponent(k) + '=' + encodeURIComponent(params[k])
  }

  return Object.keys(params)
    .map(encode)
    .join('&')
}
