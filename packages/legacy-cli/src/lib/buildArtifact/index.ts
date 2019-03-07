import * as globby from 'globby'
import * as pathJs from 'path'
import * as webpack from 'webpack'

export function transpileFunctions(entriesPath: string, distPath: string): Promise<boolean | Array<any>> {
  return new Promise((resolve, reject) => {
    // Note: it works because we have client.ts present
    globby([`${entriesPath}/*.ts`])
      .then(files => {
        if (files.length) {
          const entries = files.reduceRight(
            (entriesAcc, file) => ({
              ...entriesAcc,
              [pathJs.basename(file).split('.')[0]]: file
            }),
            {}
          )
          webpack(
            {
              mode: 'production',
              // optimization: {
              //   minimize: false
              // },
              entry: entries,
              module: {
                rules: [
                  {
                    test: /\.tsx?$/,
                    loader: 'ts-loader',
                    exclude: /node_modules/,
                    options: {
                      onlyCompileBundledFiles: true,
                      compilerOptions: {
                        allowUnreachableCode: false,
                        declaration: false,
                        lib: ['es2017'],
                        noUnusedLocals: false,
                        noUnusedParameters: false,
                        allowSyntheticDefaultImports: true,
                        experimentalDecorators: true,
                        moduleResolution: 'node',
                        module: 'es6',
                        target: 'es2017'
                      }
                    }
                  }
                ]
              },
              resolve: {
                extensions: ['.tsx', '.ts', '.js']
              },
              target: 'node',
              output: {
                libraryTarget: 'commonjs2',
                filename: '[name].js',
                path: distPath
              }
              // TODO: check if it is necessary
              // context: pathJs.resolve(path)
            },
            (err, stats) => {
              if (err || stats.hasErrors()) {
                reject(stats)
              } else {
                resolve(true)
              }
            }
          )
        } else {
          reject([{ error: 'No functions to process' }])
        }
      })
      .catch(error => {
        reject([{ error }])
      })
  })
}
