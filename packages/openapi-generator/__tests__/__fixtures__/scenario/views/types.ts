export type PullRequest = {
  id: string
  title: string
  comments: Comments
}

export type Repo = {
  id: string
  name: string
  totalCount: number
}

type Comments = {
  totalCount: number
  visible: boolean
}
