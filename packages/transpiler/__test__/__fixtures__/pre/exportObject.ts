class SpongeBobHelper {
  private name: string
  constructor(name: string) {
    this.name = name
  }
}

export class PatrickHelper {}

export default {
  spongeHelper: SpongeBobHelper,
  patrickHelper: PatrickHelper
}
