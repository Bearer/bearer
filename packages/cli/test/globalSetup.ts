import * as fs from 'fs-extra'

export default async () => {
  if (!fs.existsSync('test/.artifacts')) {
    fs.mkdirpSync('test/.artifacts')
  }
  fs.emptyDirSync('test/.artifacts')
}
