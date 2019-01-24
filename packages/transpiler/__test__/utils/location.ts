import * as path from 'path'

export const FIXTURES = path.join(__dirname, '..', '__fixtures__')
export const BUILD = path.join(__dirname, '..', '.build/src')

export function unitFixtureDirectory(dir: string) {
  return path.join(FIXTURES, 'unit', dir)
}

export function buildunitFixtureDirectory(dir: string) {
  return path.join(BUILD, 'unit', dir)
}
