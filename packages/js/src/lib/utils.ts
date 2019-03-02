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

export function formatQuery(params: Record<string, any>) {
  return Object.keys(cleanQuery(params))
    .reduce(
      (acc, key) => {
        return [...acc, [key, params[key]].join('=')]
      },
      [] as string[]
    )
    .join('&')
}
