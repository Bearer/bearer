module.exports = {
  testEnvironment: 'node',
  transform: {
    '^.+\\.tsx?$': 'ts-jest'
  },
  testRegex: '(/__tests__/.*|(\\.|/)(test|spec))\\.(jsx?|tsx?)$',
  moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx', 'json', 'node'],
  automock: false,
  coveragePathIgnorePatterns: ['<rootDir>/node_modules', '<rootDir>/test'],
  testPathIgnorePatterns: ['/node_modules/', '<rootDir>/templates', '/.bearer/', '<rootDir>/src', '<rootDir>/lib'],
  setupFiles: ['<rootDir>/test/setup.js']
}
