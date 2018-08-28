const transpiler = require('@bearer/transpiler/lib/bin/bearer-tst.js').default
transpiler(process.argv.slice(2))
process.send({ event: 'transpiler:initialized' })
