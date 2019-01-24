#!/usr/bin/env node
import * as fs from 'fs'
import * as path from 'path'
const { run } = require('@bearer/cli/lib/index')

const bearer = fs.readFileSync(path.join(__dirname, '../static/message.txt'), { encoding: 'utf8' })
console.log(bearer)

run(['new', ...process.argv.slice(2)])
  .then()
  .catch(() => {
    process.exit(0)
  })
