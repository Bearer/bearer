// Whe highly rely on Ionic Framework for UI base components

/**
 * How to use this file
 * yarn tsc scripts/*.ts --lib es6
 * node scripts/copy-components.js
 */

import { execSync } from 'child_process'
import * as fs from 'fs-extra'
/**
 * * get ionic archive
 * * extract archive
 * *¨copy interesting files
 * * rename component tag names
 * * copy tests
 */
const VERSION = 'v4.0.0-beta.3'
const ARCHIVE = `https://github.com/ionic-team/ionic/archive/${VERSION}.zip`
const TMP_ZIP = '/tmp/ionic.zip'
const TMP_DIR = '/tmp/ionic-extract'

function exec(command: string) {
  console.log('[BEARER]', 'Executing command', command)
  if (!execSync(command)) {
    console.log('[BEARER]', 'Error')
    process.exit(1)
  }
}

function rm(file: string) {
  if (fs.existsSync(file)) {
    console.log('[BEARER]', 'Deleting file', file)
    fs.removeSync(file)
  }
}

function rmDir(dir: string) {
  if (fs.existsSync(dir)) {
    console.log('[BEARER]', 'Deleting dir', dir)
    fs.emptyDirSync(dir)
  }
}

rm(TMP_ZIP)

if (!fs.existsSync(TMP_ZIP)) {
  console.log('[BEARER]', 'Dowloading archive')
  exec(`wget -q -O ${TMP_ZIP}  ${ARCHIVE}`)
}

rmDir(TMP_DIR)

console.log('[BEARER]', 'Extracting archive')
exec(`unzip -d ${TMP_DIR} -q ${TMP_ZIP} `)

const baseCore = `${TMP_DIR}/ionic-${VERSION.replace('v', '')}/core`
exec(`cp -R ${baseCore}/src/utils src/utils`)
