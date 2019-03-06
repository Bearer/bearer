const fs = require('fs')
const AWS = require('aws-sdk')

const s3 = new AWS.S3()

const { makeTmpDir, rmTmpDir } = require('../support/utils')

const { spawn, exec } = require('child_process')

const COMMAND = 'push',
  INTEGRATION_NAME = 'goatsAreFunTEST'

const params = {
  Bucket: 'vanilla-deployments',
  Key: INTEGRATION_NAME
}

let bin, tmpDir, integrationPath

describe.skip('pushCommand', () => {
  // body...

  beforeEach(() => {
    bin = `${__dirname}/../../bin/index.js`
    tmpDir = makeTmpDir()
    fs.writeFileSync(
      `${tmpDir}/index.js`,
      fs.readFileSync(`${__dirname}/../fixtures/integrations/getRepositories.js`).toString()
    )
    integrationPath = `${tmpDir}/${INTEGRATION_NAME}.zip`
  })

  afterEach(() => {
    rmTmpDir(tmpDir)
    s3.deleteObject(params, (err, data) => {
      if (err) console.log(err)
      else console.log(data)
    })
  })

  test('push command creates s3 object SLOW', done => {
    expect.assertions(2)

    const buildPackage = spawn(bin, ['package', INTEGRATION_NAME], {
      cwd: tmpDir
    })

    buildPackage.on('close', () => {
      const pushPackage = spawn(bin, [COMMAND, '--name', INTEGRATION_NAME, integrationPath])

      pushPackage.on('close', () => {
        s3.headObject(params, (err, data) => {
          if (err) console.log(err.toString())
          expect(err).toBe(null)
          expect(data).toBeDefined()
          done()
        })
      })
    })
  })

  test.skip('push command generates informative output SLOW', done => {
    expect.assertions(2)

    const buildPackage = spawn(bin, ['package', INTEGRATION_NAME], {
      cwd: tmpDir
    })

    buildPackage.on('close', () => {
      exec(
        [bin, COMMAND, '--name', INTEGRATION_NAME, integrationPath].join(' '),
        { cwd: tmpDir },
        (error, stdout, stderr) => {
          const re = new RegExp(`Pushing integration ${INTEGRATION_NAME}...`)
          expect(stdout).toMatch(re)
          expect(stdout).toMatch(/Integration has been uploaded/)
          done()
        }
      )
    })
  })
})
