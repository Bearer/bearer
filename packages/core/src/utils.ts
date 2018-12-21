export function cleanQuery(params: Record<string, any>) {
  return Object.keys(params).reduce((acc, key) => {
    if (params[key]) {
      acc[key] = params[key]
    }
    return acc
  }, {})
}

export function formatQuery(params: Record<string, any>) {
  return Object.keys(cleanQuery(params))
    .reduce((acc, key) => {
      return [...acc, [key, params[key]].join('=')]
    }, [])
    .join('&')
}
