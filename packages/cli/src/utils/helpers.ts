import * as copy from 'copy-template-dir'
import * as path from 'path'

import BaseCommand from '../BaseCommand'

export function copyFiles(
  command: BaseCommand,
  sourceDirectory: string,
  outPutDirectory: string,
  vars: any
): Promise<boolean> {
  function printFiles(files: Array<string>) {
    files.forEach(file => {
      command.log(command.colors.gray(`    create: `) + command.colors.white(file.replace(outPutDirectory + '/', '')))
    })
    command.log('\n')
  }

  return new Promise((resolve, reject) => {
    copy(path.join(__dirname, '..', '..', sourceDirectory), outPutDirectory, vars, (err, createdFiles) => {
      if (err) {
        reject(false)
      } else {
        printFiles(createdFiles)
        resolve(true)
      }
    })
  })
}
