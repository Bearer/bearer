const transpiler = require('@bearer/transpiler/dist/bin/bearer-tst.js').default
transpiler(process.argv.slice(2))
process.send({ event: 'transpiler:initialized' })
