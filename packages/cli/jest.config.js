module.exports = {
  automock: false,
  verbose: true,
  coveragePathIgnorePatterns: ['<rootDir>/node_modules', '<rootDir>/test'],
  moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx', 'json', 'node'],
  setupFiles: ['<rootDir>/test/setup.js'],
  globalSetup: '<rootDir>/test/globalSetup.ts',
  testEnvironment: 'node',
  modulePathIgnorePatterns: ['<rootDir>/templates'],
  testPathIgnorePatterns: [
    '/node_modules/',
    '<rootDir>/templates',
    '<rootDir>/test/.artifacts/',
    '<rootDir>/src/commands/generate/spec.ts',
    '<rootDir>/lib'
  ],
  testRegex: '(/__tests__/.*|(\\.|/)(test|spec))\\.(jsx?|tsx?)$',
  transform: {
    '^.+\\.tsx?$': 'ts-jest'
  },
  collectCoverageFrom: ['<rootDir>/src/**/*.{ts,tsx}', '!<rootDir>/src/**/*.{spec,test}.{ts,tsx}']
}
