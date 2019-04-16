import * as inquirer from 'inquirer'
type Omit<T, K> = Pick<T, Exclude<keyof T, K>>

export type Options = Partial<Omit<inquirer.Question, 'message' | 'name'>>

export async function askForString(phrase: string, options: Options = {}): Promise<string> {
  const { response } = await inquirer.prompt<{ response: string }>([
    {
      message: `${phrase}:`,
      name: 'response',
      ...options
    }
  ])
  return response
}

export async function askForPassword(phrase: string, options: Options = {}): Promise<string> {
  return askForString(phrase, { ...options, type: 'password' })
}
