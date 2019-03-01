import * as React from 'react'

export interface IConnectProps {
  integration: string
  setupId: string
  onSuccess: (detail: { authId: string }) => void
  render: (props: { loading: boolean; connect: () => void; error: any }) => JSX.Element
}

class Connect extends React.Component<IConnectProps, {}> {
  constructor(props) {
    super(props)
    this.connect = this.connect.bind(this)
  }

  connect() {
    setTimeout(() => {
      this.props.onSuccess({ authId: 'sponge bob' })
    }, 2000)
  }
  render() {
    return this.props.render({ loading: false, connect: this.connect, error: null })
  }
}

export default Connect
