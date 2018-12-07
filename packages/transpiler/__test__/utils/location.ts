import * as path from 'path'

export const FIXTURES = path.join(__dirname, '..', '__fixtures__')
export const BUILD = path.join(__dirname, '..', '.build/src')

export function UnitFixtureDirectory(dir: string) {
  return path.join(FIXTURES, 'unit', dir)
}

export function BuildUnitFixtureDirectory(dir: string) {
  return path.join(BUILD, 'unit', dir)
}
