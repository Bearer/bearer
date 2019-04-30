import * as fs from 'fs-extra'
import * as path from 'path'

export function readFile(...args: string[]): string {
  try {
    return fs.readFileSync(path.join(...args), { encoding: 'utf8' })
  } catch (err) {
    if (err.code === 'ENOENT') {
      return `Not found: ${path.join(...args)}`
    }
    return err.toString()
  }
}

export const ARTIFACT_FOLDER = path.join(__dirname, '..', '.artifacts')
export const FIXTURE_FOLDER = path.join(__dirname, '..', '__FIXTURES__')

export function fixturesPath(...args: string[]) {
  return path.join(FIXTURE_FOLDER, ...args)
}

export function artifactPath(...args: string[]) {
  return path.join(ARTIFACT_FOLDER, ...args)
}
