import Bearer, { Component, Prop, State } from '@bearer/core'

@Component({
  tag: 'bearer-setup-display',
  styleUrl: 'setup-display.scss',
  shadow: true
})
export class BearerSetupDisplay {
  @Prop()
  integrationId = ''
  @State()
  isSetup: boolean = false
  @Prop({ mutable: true })
  setupId = ''

  componentDidLoad() {
    Bearer.emitter.addListener(`setup_success:${this.integrationId}`, data => {
      this.setupId = data.referenceId
      this.isSetup = true
    })
  }

  render() {
    const label = <strong>Setup-id</strong>
    if (this.isSetup || this.setupId) {
      return (
        <div>
          {label}:&nbsp; {this.setupId}
        </div>
      )
    }
    return (
      <div>
        {label}:&nbsp;{' '}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          xmlnsXlink="http://www.w3.org/1999/xlink"
          version="1.1"
          x="0px"
          y="0px"
          viewBox="0 0 100 100"
          enable-background="new 0 0 100 100"
          xmlSpace="preserve"
          class="icon-svg"
        >
          <g>
            <g>
              <g>
                <path
                  fill="#E4611F"
                  // tslint:disable-next-line:max-line-length
                  d="M87.999,92.65H12.097c-4.451,0-8.034-1.783-9.83-4.896c-1.795-3.108-1.549-7.104,0.677-10.957     l37.95-65.733C43.122,7.21,46.458,5,50.049,5c3.592,0,6.928,2.211,9.153,6.065L97.15,76.799     c2.227,3.854,2.474,7.848,0.679,10.957C96.033,90.867,92.45,92.65,87.999,92.65z M50.049,12.978     c-0.521,0-1.417,0.642-2.245,2.077L9.853,80.787c-0.829,1.434-0.937,2.531-0.677,2.982c0.26,0.448,1.264,0.903,2.921,0.903     h75.902c1.657,0,2.66-0.455,2.921-0.905c0.26-0.449,0.15-1.545-0.678-2.98l-37.95-65.733     C51.465,13.62,50.568,12.978,50.049,12.978z"
                />
              </g>
            </g>
            <path
              fill="#E4611F"
              // tslint:disable-next-line:max-line-length
              d="M45.313,57.395c0.21,2.179,0.563,3.781,1.081,4.9c0.65,1.408,1.856,2.152,3.49,2.152   c1.599,0,2.813-0.754,3.513-2.18c0.569-1.164,0.928-2.744,1.092-4.818l1.406-16.12c0.155-1.509,0.234-3.025,0.234-4.504   c0-2.643-0.347-4.642-1.064-6.113c-0.576-1.181-1.875-2.587-4.787-2.587c-1.872,0-3.413,0.633-4.58,1.881   c-1.149,1.229-1.732,2.916-1.732,5.02c0,1.356,0.1,3.593,0.297,6.655L45.313,57.395z"
            />
            <path
              fill="#E4611F"
              // tslint:disable-next-line:max-line-length
              d="M50.016,67.895c-1.683,0-3.126,0.595-4.292,1.758c-1.167,1.166-1.758,2.601-1.758,4.259   c0,1.889,0.633,3.396,1.882,4.479c1.201,1.045,2.625,1.571,4.233,1.571c1.59,0,3.001-0.535,4.194-1.596   c1.229-1.095,1.854-2.591,1.854-4.455c0-1.664-0.604-3.101-1.8-4.269C53.143,68.484,51.691,67.895,50.016,67.895z"
            />
          </g>
        </svg>{' '}
        Not set
      </div>
    )
  }
}
