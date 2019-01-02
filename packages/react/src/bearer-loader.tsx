import * as React from 'react'
import { Helmet } from 'react-helmet'

const bearerScriptTag = (clientId: string, bearerHost: string) =>
  // tslint:disable-next-line max-line-length
  `!function(){var t=window.bearer=window.bearer||[];t.initialize||(t.invoked?window.console&&console.error&&console.error("Bearer snippet included twice."):(t.invoked=!0,t.load=function(e,o){var r=document.createElement("script");r.type="text/javascript",r.async=!0,r.src="${bearerHost}"+e+"/bearer.min.js";var n=document.getElementsByTagName("script")[0];n.parentNode.insertBefore(r,n),t._loadOptions=o,t.clientId=e},t.SNIPPET_VERSION="1.0.0",t.load("${clientId}")))}();`

interface IBearerLoaderProps {
  clientId: string
  intHost?: string
}

export default (props: IBearerLoaderProps) => (
  <Helmet>
    <script>{bearerScriptTag(props.clientId, props.intHost || 'https://int.bearer.sh/v1/')}</script>
  </Helmet>
)
