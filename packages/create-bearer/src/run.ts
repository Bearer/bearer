import { run } from '@bearer/cli/lib/index'

export default async (...args: any[]) => {
  let isTemplate = false
  const newArgs = args.reduce((acc, current) => {
    if (isTemplate) {
      if (current.startsWith('-')) {
        acc.push('https://github.com/Bearer/templates')
      }
      isTemplate = false
    } else {
      isTemplate = current === '-t' || current === '--template'
    }
    acc.push(current)
    return acc
  }, [])
  // template args is the latest
  if (isTemplate) {
    newArgs.push('https://github.com/Bearer/templates')
  }
  await run(['new', ...newArgs])
}
