#!/usr/bin/env node
import * as fs from 'fs'
import * as path from 'path'
import run from './run'

const bearer = fs.readFileSync(path.join(__dirname, '../static/message.txt'), { encoding: 'utf8' })
console.log(bearer)

run(...process.argv.slice(2))
  .then()
  .catch(() => {
    process.exit(0)
  })
