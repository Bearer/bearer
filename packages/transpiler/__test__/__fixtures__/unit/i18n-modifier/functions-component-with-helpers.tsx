import { t, p } from '@bearer/core'

export default () => (
  <span>
    {t('key.target', 'default value')}
    {p('key.target', 0, 'default value')}
  </span>
)
