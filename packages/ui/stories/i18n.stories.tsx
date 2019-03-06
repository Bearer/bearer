import * as React from 'react'
import { storiesOf } from '@storybook/react'

storiesOf('I18n translate', module)
  .addWithJSX('Missing Key: default Value', () => [
    <bearer-i18n _key="missing.key" var='{"ok": "Sponge bob"}' default="Default value : {{ok}}" />
  ])
  .addWithJSX('Existing Key', () => [
    <bearer-i18n _key="existing.key" var='{"ok": "🐻"}' default="Not displayed : {{ok}}" />
  ])

storiesOf('I18n pluralize', module)
  .addWithJSX('Missing Key: default Value', () => [
    <integration-placeholder>
      <bearer-i18n _key="missing.key" count="0" default="A default text : {{count}}" />
    </integration-placeholder>
  ])
  .addWithJSX('Existing Key: 0', () => [
    <bearer-i18n _key="withCount" count="0" default="A default text : {{count}}" />
  ])
  .addWithJSX('Existing Key: 1', () => [
    <bearer-i18n _key="withCount" count="1" default="A default text : {{count}}" />
  ])
  .addWithJSX('Existing Key: 2', () => [
    <bearer-i18n _key="withCount" count="2" default="A default text : {{count}}" />
  ])
  .addWithJSX('Existing Key: 3', () => [
    <bearer-i18n _key="withCount" count="3" default="A default text : {{count}}" />
  ])
