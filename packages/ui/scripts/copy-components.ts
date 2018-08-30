// Whe highly rely on Ionic Framework for UI base components

/**
 * How to use this file
 * yarn tsc scripts/*.ts --lib es6
 * node scripts/copy-components.js
 */

import { execSync } from 'child_process'
import * as fs from 'fs-extra'
import * as path from 'path'

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

import compile from './compiler'

function exec(command: string, verbose: boolean = false) {
  console.log('[BEARER]', 'Executing command', command)
  const result = execSync(command)
  if (!result) {
    console.log('[BEARER]', 'Error')
    process.exit(1)
  } else if (verbose) {
    console.log(result.toString('utf8'))
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

// Copy global folder
exec(`cp -r ${baseCore}/src/global/ src/global/`, true)

// Copy required components
const components = [
  'popover',
  'animation-controller',
  'route',
  'router',
  'route-redirect',
  'badge',
  'grid',
  'button',
  'checkbox',
  'item',
  'item-sliding',
  'item-options',
  'item-option',
  'searchbar',
  'input',
  'list',
  'slides',
  'slide',
  'spinner',
  'toast',
  'card',
  'card-content',
  'card-header',
  'card-subtitle',
  'card-title',
  'modal',
  'loading',
  'loading-controller',
  'col',
  'row',
  'label'
]
components.map(component => {
  exec(`rsync -avz ${baseCore}/src/components/${component}/ src/components/${component}`)
})

// Copy required utils files
const includes = [
  'theme',
  'framework-delegate',
  'overlays',
  'overlays-interface',
  'input-interface',
  'gesture',
  'gesture-controller',
  'listener',
  'pointer-events',
  'recognizers',
  'config',
  'platform',
  'helpers',
  'media'
]
  .map(u => `--include="${u}.ts"`)
  .join(' ')

exec(`rsync -avz  --include="*/" ${includes} --exclude="*" ${baseCore}/src/utils/ src/utils`)

// Copy interface
exec(`cp ${baseCore}/src/interface.d.ts src/`)

// Comment non imported files
const patterns = [
  'action-sheet',
  'menu-interface',
  'alert-interface',
  'picker-interface',
  'nav-interface',
  'range-interface',
  'content-interface',
  'select-interface',
  'select-popover-interface',
  'tabbar-interface',
  'virtual-scroll-interface',
  'view-controller'
]

patterns.map(p => {
  exec(`sed -i '' -e '/^export.*${p}/s/^/\\/\\//g' src/interface.d.ts`)
})

// Copy theme files
exec(`rsync -avz  ${baseCore}/src/themes/ src/themes`)

// Rename component => ion-* to bearer-* , add styleUrl to fallback to md style by default
components.map(comp => {
  compile([path.join(`src/components/${comp}/${comp}.tsx`)])
})

// edge cases: replace manually created element
const filesWithIonPrefix = [
  'src/components/item/item.tsx',
  'src/components/loading-controller/loading-controller.tsx',
  'src/components/item-sliding/item-sliding.tsx',
  'src/utils/theme.ts'
]

filesWithIonPrefix.map(f => {
  exec(`sed -i '' -e 's/ion-/bearer-/g' ${f}`)
})

// replace HTMLIon by HTMLBearer
const files = ['src/components/router/utils/parser.ts', 'src/utils/overlays-interface.ts', 'src/utils/overlays.ts']

files.map(f => {
  exec(`sed -i '' -e 's/HTMLIon/HTMLBearer/g' ${f}`)
})

// add default to mode to list
exec(`sed -i '' -e "s/mode\\!/mode/g" src/components/list/list.tsx`)
exec(`sed -i '' -e "s/Mode\\;/Mode = 'md';/g" src/components/list/list.tsx`)

//Styles
// replace css variables, tag-name, .class-name
// exec(`find src/components/**/*.scss | xargs sed -i '' -e 's/ion-/bearer-/g'`)

//classnames
exec("find src/components/**/*.scss | xargs sed -i '' -e 's/\\.ion-/\\.bearer-/g'")
exec("find src/themes/*.scss | xargs sed -i '' -e 's/\\.ion-/\\.bearer-/g'")
// variables + functions
exec("find src/components/**/*.scss | xargs sed -i '' -e 's/\\([^a-z]\\)ion-/\\1bearer-/g'")
exec("find src/themes/*.scss | xargs sed -i '' -e 's/\\([^a-z]\\)ion-/\\1bearer-/g'")
exec("find src/themes/*.scss | xargs sed -i '' -e 's/--ion-/--bearer-/g'")
// tag names
exec("find src/components/**/*.scss | xargs sed -i '' -e 's/^ion-/bearer-/g'")

// adjust list styles
