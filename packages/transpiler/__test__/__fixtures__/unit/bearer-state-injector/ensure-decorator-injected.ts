import { BearerState } from '@bearer/core'

class AddDecoratorOnTranspile {
  @BearerState() pullRequests: Array<any> = []
}
