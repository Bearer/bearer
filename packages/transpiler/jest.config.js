module.exports = {
  automock: false,
  moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx', 'json', 'node'],
  setupFiles: ['<rootDir>/__test__/utils/setup.ts'],
  testPathIgnorePatterns: ['/node_modules/', '/__fixtures__/', '/utils/', '/spec.ts', 'lib/'],
  testRegex: '(/__tests__/.*|(\\.|/)(test|spec))\\.(jsx?|tsx?)$',
  transform: { '^.+\\.tsx?$': 'ts-jest' }
}
