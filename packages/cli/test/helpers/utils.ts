import * as fs from 'fs-extra'
import * as path from 'path'

export function readFile(...args: string[]): string {
  return fs.readFileSync(path.join(...args), { encoding: 'utf8' })
}
