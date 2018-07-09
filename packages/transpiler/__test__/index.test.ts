import bearer from '../src/transformers/bearer'
test('bearer addImport updates the import clause', () => {})
// test('invoking transpiler', async () => {
//   let transpiler = new Transpiler(SRC_DIRECTORY)
//   await transpiler.run()
//   expect(
//     fs.existsSync(path.join(BUILD_DIRECTORY, 'exportObject.ts'))
//   ).toBeTruthy()
//   expect(
//     fs.existsSync(path.join(BUILD_DIRECTORY, 'classComponent.ts'))
//   ).toBeTruthy()
// })

// test('Adding BEARER_ID prop', async () => {
//   pending('circular calls')
//   let transpiler = new Transpiler(SRC_DIRECTORY)
//   await transpiler.run()
//   const builtFilePath = path.join(BUILD_DIRECTORY, 'classComponent.ts')
// })
