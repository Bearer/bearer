import { BearerEnv, BaseConfig } from '../types'

export const LOGIN_CLIENT_ID = process.env.BEARER_LOGIN_CLIENT_ID || 'Wgll39KqWnJWud473wq7hZhiXxeNjEU7'
export const BEARER_ENV = process.env.BEARER_ENV || 'production'
export const BEARER_AUTH_PORT = 45677
export const BEARER_LOGIN_PORT = 56789

export const CONFIGS: Record<BearerEnv, BaseConfig> = {
  dev: {
    IntegrationServiceHost: 'https://int.dev.bearer.sh/',
    IntegrationServiceUrl: 'https://int.dev.bearer.sh/api/v1/',
    DeveloperPortalAPIUrl: 'https://api.staging.bearer.sh/graphql',
    DeveloperPortalUrl: 'https://app.staging.bearer.sh/',
    CdnHost: 'https://static.dev.bearer.sh',
    BearerEnv: 'dev',
    LoginDomain: 'https://login.bearer.sh'
  },
  staging: {
    IntegrationServiceHost: 'https://int.staging.bearer.sh/',
    IntegrationServiceUrl: 'https://int.staging.bearer.sh/api/v1/',
    DeveloperPortalAPIUrl: 'https://api.staging.bearer.sh/graphql',
    DeveloperPortalUrl: 'https://app.staging.bearer.sh/',
    CdnHost: 'https://static.staging.bearer.sh',
    BearerEnv: 'staging',
    LoginDomain: 'https://login.bearer.sh'
  },
  production: {
    IntegrationServiceHost: 'https://int.bearer.sh/',
    IntegrationServiceUrl: 'https://int.bearer.sh/api/v1/',
    DeveloperPortalAPIUrl: 'https://api.bearer.sh/graphql',
    DeveloperPortalUrl: 'https://app.bearer.sh/',
    CdnHost: 'https://static.bearer.sh',
    BearerEnv: 'production',
    LoginDomain: 'https://login.bearer.sh'
  }
}

// tslint:disable max-line-length
export const SUCCESS_LOGIN_PAGE = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <title>Authentication callback</title>
    <meta charset="utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>
      * {
        -webkit-font-smoothing: antialiased;
        -moz-osx-font-smoothing: grayscale;
      }
      html,
      body {
        background-color: #f5f7fb;
        font-family: 'Proxima Nova', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif,
          'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol';
        text-align: center;
        font-size: 16px;
        line-height: 16px;
      }
      h1 {
        color: #00c682;
        font-size: 2rem;
        font-weight: 600;
        letter-spacing: 0.99px;
        line-height: 29px;
      }
      p {
        color: #343c5d;
        letter-spacing: 0.56px;
      }
      a {
        border-radius: 4px;
        display: inline-block;
        margin-top: 32px;
        padding: 12px 18px;
        background-color: #030d36;
        color: #ffffff;
        font-weight: 600;
        letter-spacing: 0.2px;
        text-decoration: none;
      }
      .outer {
        display: table;
        position: absolute;
        top: 0;
        left: 0;
        height: 100%;
        width: 100%;
      }
      .middle {
        display: table-cell;
        vertical-align: middle;
      }
      .inner {
        margin-left: auto;
        margin-right: auto;
        max-width: 700px;
      }
      .hint {
        font-size: 0.9rem;
        padding: 1.5rem;
        border: 1px solid #c2c9ea;
        border-radius: 4px;
        background-color: #ffffff;
        position: relative;
        top: 60px;
        margin-bottom: 60px;
      }
    </style>
  </head>
  <body>
    <div class="outer">
      <div class="middle">
        <div class="inner">
          <svg width="37" height="40" xmlns="http://www.w3.org/2000/svg">
            <path
              d="M.8 19.2h25.269l-6.635-6.634a.8.8 0 0 1 1.132-1.132l8 8a.8.8 0 0 1 0 1.132l-8 8a.8.8 0 0 1-1.132-1.132L26.07 20.8H.8a.8.8 0 0 1 0-1.6zm14.4 11.2a.8.8 0 0 1 .8.8v6.4a.8.8 0 0 0 .8.8h17.6a.8.8 0 0 0 .8-.8V2.4a.8.8 0 0 0-.8-.8H16.8a.8.8 0 0 0-.8.8v6.4a.8.8 0 0 1-1.6 0V2.4A2.4 2.4 0 0 1 16.8 0h17.6a2.4 2.4 0 0 1 2.4 2.4v35.2a2.4 2.4 0 0 1-2.4 2.4H16.8a2.4 2.4 0 0 1-2.4-2.4v-6.4a.8.8 0 0 1 .8-.8z"
              fill="#00C682"
              fill-rule="nonzero"
            />
          </svg>
          <h1>Successfully authenticated</h1>
          <p>You can close this window</p>
          <br />
        </div>
      </div>
    </div>
  </body>
</html>
`
