import * as fs from 'fs-extra'
import * as path from 'path'

export function readFile(...args: string[]): string {
  try {
    return fs.readFileSync(path.join(...args), { encoding: 'utf8' })
  } catch (err) {
    return err.toString()
  }
}
