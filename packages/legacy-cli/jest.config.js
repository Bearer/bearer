module.exports = {
  transform: {
    '^.+\\.tsx?$': 'ts-jest'
  },
  testRegex: '(/__tests__/.*|(\\.|/)(test|spec))\\.(jsx?|tsx?)$',
  moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx', 'json', 'node'],
  automock: false,
  verbose: true,
  testEnvironment: 'node',
  testURL: 'http://localhost/',
  testPathIgnorePatterns: ['<rootDir>/dist/', '<rootDir>/node_modules']
}
