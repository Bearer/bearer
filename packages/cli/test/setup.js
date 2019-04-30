// ease output comparison
// listr renders differently if isTTY is not true
process.stdout.isTTY = undefined
jest.setTimeout(10000)
