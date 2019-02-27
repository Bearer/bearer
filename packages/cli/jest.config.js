module.exports = {
  automock: false,
  coveragePathIgnorePatterns: ['<rootDir>/node_modules', '<rootDir>/test'],
  moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx', 'json', 'node'],
  setupFiles: ['<rootDir>/test/setup.js'],
  testEnvironment: 'node',
  modulePathIgnorePatterns: ['<rootDir>/templates'],
  testPathIgnorePatterns: [
    '/node_modules/',
    '<rootDir>/templates',
    '/.bearer/',
    '<rootDir>/src/commands/generate/spec.ts',
    '<rootDir>/lib'
  ],
  testRegex: '(/__tests__/.*|(\\.|/)(test|spec))\\.(jsx?|tsx?)$',
  transform: {
    '^.+\\.tsx?$': 'ts-jest'
  }
}
