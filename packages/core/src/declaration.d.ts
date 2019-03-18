import { Bearer } from '@bearer/js'

// Since it is a peerDependency, we defined it to have it present in the TS wolrd
declare global {
  const bearer: Bearer

  interface Window {
    bearer: Bearer
  }
}
