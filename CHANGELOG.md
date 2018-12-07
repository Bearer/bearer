# Change Log

All notable changes to this project will be documented in this file.
See [Conventional Commits](https://conventionalcommits.org) for commit guidelines.

<a name="0.74.4"></a>
## [0.74.4](https://github.com/Bearer/bearer/compare/v0.74.3...v0.74.4) (2018-12-07)


### Bug Fixes

* **cli:** use correct intent variable name reference ([#411](https://github.com/Bearer/bearer/issues/411)) ([fe579d8](https://github.com/Bearer/bearer/commit/fe579d8))





<a name="0.74.3"></a>
## [0.74.3](https://github.com/Bearer/bearer/compare/v0.74.2...v0.74.3) (2018-12-05)


### Bug Fixes

* **cli:** use intent filenames as intent names ([#408](https://github.com/Bearer/bearer/issues/408)) ([bafec2e](https://github.com/Bearer/bearer/commit/bafec2e))
* make local development work  ([#409](https://github.com/Bearer/bearer/issues/409)) ([1f930f1](https://github.com/Bearer/bearer/commit/1f930f1))





<a name="0.74.2"></a>
## [0.74.2](https://github.com/Bearer/bearer/compare/v0.74.1...v0.74.2) (2018-12-04)


### Bug Fixes

* **transpiler:** Core 150 event emitter issue ([#404](https://github.com/Bearer/bearer/issues/404)) ([5265573](https://github.com/Bearer/bearer/commit/5265573))





<a name="0.74.1"></a>
## [0.74.1](https://github.com/Bearer/bearer/compare/v0.74.0...v0.74.1) (2018-12-03)


### Bug Fixes

* **transpiler:** change state refId as to a prop ([#403](https://github.com/Bearer/bearer/issues/403)) ([249a176](https://github.com/Bearer/bearer/commit/249a176))





<a name="0.74.0"></a>
# [0.74.0](https://github.com/Bearer/bearer/compare/v0.73.2...v0.74.0) (2018-12-01)


### Bug Fixes

* export correct name for loader ([#401](https://github.com/Bearer/bearer/issues/401)) ([f129730](https://github.com/Bearer/bearer/commit/f129730))
* **transpiler:** [CORE-135] - output event typing ([#402](https://github.com/Bearer/bearer/issues/402)) ([01e3134](https://github.com/Bearer/bearer/commit/01e3134))


### Features

* **transpiler:** allow using existing property ([#400](https://github.com/Bearer/bearer/issues/400)) ([64f635d](https://github.com/Bearer/bearer/commit/64f635d))





<a name="0.73.2"></a>
## [0.73.2](https://github.com/Bearer/bearer/compare/v0.73.1...v0.73.2) (2018-11-23)


### Bug Fixes

* **legacy-cli:** prevent cli to fail when non existing ref is passed ([4aaaa08](https://github.com/Bearer/bearer/commit/4aaaa08))





<a name="0.73.1"></a>
## [0.73.1](https://github.com/Bearer/bearer/compare/v0.73.0...v0.73.1) (2018-11-23)

**Note:** Version bump only for package bearer-master





<a name="0.73.0"></a>
# [0.73.0](https://github.com/Bearer/bearer/compare/v0.72.1...v0.73.0) (2018-11-16)


### Features

* input/ouput options ([#393](https://github.com/Bearer/bearer/issues/393)) ([b231d86](https://github.com/Bearer/bearer/commit/b231d86))
* output decorator ([#391](https://github.com/Bearer/bearer/issues/391)) ([b388d9f](https://github.com/Bearer/bearer/commit/b388d9f))





<a name="0.72.1"></a>
## [0.72.1](https://github.com/Bearer/bearer/compare/v0.72.0...v0.72.1) (2018-11-15)


### Bug Fixes

* **core:** build issue ([#392](https://github.com/Bearer/bearer/issues/392)) ([66166e2](https://github.com/Bearer/bearer/commit/66166e2))





<a name="0.72.0"></a>
# [0.72.0](https://github.com/Bearer/bearer/compare/v0.71.4...v0.72.0) (2018-11-13)


### Bug Fixes

* **cli:** add soft option to generate open api ([dffe0ef](https://github.com/Bearer/bearer/commit/dffe0ef))
* **core:** resolve authorize per scenarioId ([e06834d](https://github.com/Bearer/bearer/commit/e06834d))


### Features

* **transpiler:** input decorator ([#383](https://github.com/Bearer/bearer/issues/383)) ([e1fac50](https://github.com/Bearer/bearer/commit/e1fac50))





<a name="0.71.4"></a>
## [0.71.4](https://github.com/Bearer/bearer/compare/v0.71.3...v0.71.4) (2018-11-09)


### Bug Fixes

* **core:** send clientId with every intent calls ([07cc014](https://github.com/Bearer/bearer/commit/07cc014))





<a name="0.71.3"></a>
## [0.71.3](https://github.com/Bearer/bearer/compare/v0.71.2...v0.71.3) (2018-10-30)


### Bug Fixes

* **slack-skin:** slack: add success and error state ([#386](https://github.com/Bearer/bearer/issues/386)) ([0f98c58](https://github.com/Bearer/bearer/commit/0f98c58))





<a name="0.71.2"></a>
## [0.71.2](https://github.com/Bearer/bearer/compare/v0.71.1...v0.71.2) (2018-10-29)


### Bug Fixes

* **ui:** Slack: add outline button prop, match style to LP ([65b4c3f](https://github.com/Bearer/bearer/commit/65b4c3f))
* **ui:** slack: set border color to gray ([1fbbd35](https://github.com/Bearer/bearer/commit/1fbbd35))





<a name="0.71.1"></a>
## [0.71.1](https://github.com/Bearer/bearer/compare/v0.71.0...v0.71.1) (2018-10-26)


### Bug Fixes

* **slack-skin:** create first btn-light styles ([f0870f3](https://github.com/Bearer/bearer/commit/f0870f3))
* **transpiler:** use global instead of no-group ([6ef6133](https://github.com/Bearer/bearer/commit/6ef6133))


### Performance Improvements

* **transpiler:** prevent processing if declaration file ([dbc801b](https://github.com/Bearer/bearer/commit/dbc801b))





<a name="0.71.0"></a>
# [0.71.0](https://github.com/Bearer/bearer/compare/v0.70.3...v0.71.0) (2018-10-24)


### Bug Fixes

* **ui:** add missin component definition ([6f61f1f](https://github.com/Bearer/bearer/commit/6f61f1f))


### Features

* **core:** allow log level customization ([fef36dc](https://github.com/Bearer/bearer/commit/fef36dc))
* **core:** remove log statements ([feba2f1](https://github.com/Bearer/bearer/commit/feba2f1))
* **ui:** change log to debug ([f3c3464](https://github.com/Bearer/bearer/commit/f3c3464))
* **ui:** remove log statement ([3540dff](https://github.com/Bearer/bearer/commit/3540dff))





<a name="0.70.3"></a>
## [0.70.3](https://github.com/Bearer/bearer/compare/v0.70.2...v0.70.3) (2018-10-22)


### Bug Fixes

* **core:** fix incorrect returned payload type ([8475503](https://github.com/Bearer/bearer/commit/8475503))





<a name="0.70.2"></a>
## [0.70.2](https://github.com/Bearer/bearer/compare/v0.70.1...v0.70.2) (2018-10-18)


### Bug Fixes

* **create-bearer:** wording ([612e1f6](https://github.com/Bearer/bearer/commit/612e1f6))





<a name="0.70.1"></a>
## [0.70.1](https://github.com/Bearer/bearer/compare/v0.70.0...v0.70.1) (2018-10-18)


### Bug Fixes

* **cli:** create empty openapi.json when starting ([9dd6f27](https://github.com/Bearer/bearer/commit/9dd6f27))
* **create-bearer:** add link to corporate site ([a8fab3c](https://github.com/Bearer/bearer/commit/a8fab3c))





<a name="0.70.0"></a>
# [0.70.0](https://github.com/Bearer/bearer/compare/v0.69.0...v0.70.0) (2018-10-18)


### Features

* **create-bearer:** init repo and make sure we can use npm init bearer ([70db666](https://github.com/Bearer/bearer/commit/70db666))





<a name="0.69.0"></a>
# [0.69.0](https://github.com/Bearer/bearer/compare/v0.68.1...v0.69.0) (2018-10-17)


### Features

* **core:** improve bearer fetch typing experience ([15a9373](https://github.com/Bearer/bearer/commit/15a9373))
* **core:** remove useless event triggered ([01ea31f](https://github.com/Bearer/bearer/commit/01ea31f))





<a name="0.68.1"></a>
## [0.68.1](https://github.com/Bearer/bearer/compare/v0.68.0...v0.68.1) (2018-10-16)


### Bug Fixes

* **cli:** tag name composition ([#373](https://github.com/Bearer/bearer/issues/373)) ([2585ab6](https://github.com/Bearer/bearer/commit/2585ab6))





<a name="0.68.0"></a>
# [0.68.0](https://github.com/Bearer/bearer/compare/v0.67.6...v0.68.0) (2018-10-16)


### Features

* **transpiler:** update component tag name parts ordering ([#372](https://github.com/Bearer/bearer/issues/372)) ([8d54be6](https://github.com/Bearer/bearer/commit/8d54be6))





<a name="0.67.6"></a>
## [0.67.6](https://github.com/Bearer/bearer/compare/v0.67.5...v0.67.6) (2018-10-15)

**Note:** Version bump only for package bearer-master





<a name="0.67.5"></a>
## [0.67.5](https://github.com/Bearer/bearer/compare/v0.67.4...v0.67.5) (2018-10-11)


### Bug Fixes

* don't generate the openapi.json on start ([3c8d92b](https://github.com/Bearer/bearer/commit/3c8d92b))
* **cli:** move prepareCOnfig to cli ([6486356](https://github.com/Bearer/bearer/commit/6486356))





<a name="0.67.4"></a>
## [0.67.4](https://github.com/Bearer/bearer/compare/v0.67.3...v0.67.4) (2018-10-11)


### Bug Fixes

* don't use vm use ts ([#367](https://github.com/Bearer/bearer/issues/367)) ([129f8f0](https://github.com/Bearer/bearer/commit/129f8f0))





<a name="0.67.3"></a>
## [0.67.3](https://github.com/Bearer/bearer/compare/v0.67.2...v0.67.3) (2018-10-10)


### Bug Fixes

* **cli:** fix refresh token regression ([#366](https://github.com/Bearer/bearer/issues/366)) ([d59effd](https://github.com/Bearer/bearer/commit/d59effd))





<a name="0.67.2"></a>
## [0.67.2](https://github.com/Bearer/bearer/compare/v0.67.1...v0.67.2) (2018-10-10)


### Bug Fixes

* **cli:** be flexible on alice code expectation ([#364](https://github.com/Bearer/bearer/issues/364)) ([9f52a52](https://github.com/Bearer/bearer/commit/9f52a52))





<a name="0.67.1"></a>
## [0.67.1](https://github.com/Bearer/bearer/compare/v0.67.0...v0.67.1) (2018-10-10)


### Bug Fixes

* **client:** put openapi.json to the correct folder ([d789cc1](https://github.com/Bearer/bearer/commit/d789cc1))





<a name="0.67.0"></a>
# [0.67.0](https://github.com/Bearer/bearer/compare/v0.66.1...v0.67.0) (2018-10-09)


### Bug Fixes

* **transpiler:** Add missing setup id ([#363](https://github.com/Bearer/bearer/issues/363)) ([d3293e8](https://github.com/Bearer/bearer/commit/d3293e8))


### Features

* add v2 storage endpoints for local dev ([d0a091a](https://github.com/Bearer/bearer/commit/d0a091a))
* build openapi.json doc ([690bd4b](https://github.com/Bearer/bearer/commit/690bd4b))





<a name="0.66.1"></a>
## [0.66.1](https://github.com/Bearer/bearer/compare/v0.66.0...v0.66.1) (2018-10-09)


### Bug Fixes

* **transpiler:** do not fail if no initiliazer is given ([#360](https://github.com/Bearer/bearer/issues/360)) ([dfb25d5](https://github.com/Bearer/bearer/commit/dfb25d5))





<a name="0.66.0"></a>
# [0.66.0](https://github.com/Bearer/bearer/compare/v0.65.1...v0.66.0) (2018-10-09)


### Bug Fixes

* **cli:** add trailing slash to the baseURL for axios/oauth2 ([#354](https://github.com/Bearer/bearer/issues/354)) ([c81c481](https://github.com/Bearer/bearer/commit/c81c481))
* **cli:** default value instead of default option ([#353](https://github.com/Bearer/bearer/issues/353)) ([680c4f5](https://github.com/Bearer/bearer/commit/680c4f5))
* **cli:** disable Stencil cache ([a7b7afa](https://github.com/Bearer/bearer/commit/a7b7afa))
* **cli:** file were not re-transpiled anymore ([585f9e9](https://github.com/Bearer/bearer/commit/585f9e9))
* **transpiler:** inject correct event name within output ([#359](https://github.com/Bearer/bearer/issues/359)) ([f370cde](https://github.com/Bearer/bearer/commit/f370cde))


### Features

* **cli:** inject Input to manifest file from props declaration ([#355](https://github.com/Bearer/bearer/issues/355)) ([b5d66c8](https://github.com/Bearer/bearer/commit/b5d66c8))
* **transpiler:** Event name refactoring ([#357](https://github.com/Bearer/bearer/issues/357)) ([3a0d011](https://github.com/Bearer/bearer/commit/3a0d011))
* **transpiler:** listen event refactoring ([#358](https://github.com/Bearer/bearer/issues/358)) ([b5294a8](https://github.com/Bearer/bearer/commit/b5294a8))
* **transpiler:** Output component api to manifest ([#356](https://github.com/Bearer/bearer/issues/356)) ([0363b5e](https://github.com/Bearer/bearer/commit/0363b5e))





<a name="0.65.1"></a>
## [0.65.1](https://github.com/Bearer/bearer/compare/v0.65.0...v0.65.1) (2018-10-03)


### Bug Fixes

* **ui:** wrong component defintition ([8b85992](https://github.com/Bearer/bearer/commit/8b85992))





<a name="0.65.0"></a>
# [0.65.0](https://github.com/Bearer/bearer/compare/v0.64.1...v0.65.0) (2018-10-02)


### Bug Fixes

* **cli:** use node module resolution ([0313f4b](https://github.com/Bearer/bearer/commit/0313f4b))
* **ui:** return to first screen when popover is closed ([dab9cef](https://github.com/Bearer/bearer/commit/dab9cef))


### Features

* add input and output eventually to preview root components ([331df53](https://github.com/Bearer/bearer/commit/331df53))




<a name="0.64.1"></a>
## [0.64.1](https://github.com/Bearer/bearer/compare/v0.64.0...v0.64.1) (2018-09-28)


### Bug Fixes

* **core:** fix incorrect destructuration ([5ebef63](https://github.com/Bearer/bearer/commit/5ebef63))
* **transpiler:** add base of work ([aeff483](https://github.com/Bearer/bearer/commit/aeff483))
* **transpiler:** add scenarioId prop to auth-screen & bearer-authorized ([7606356](https://github.com/Bearer/bearer/commit/7606356))
* **transpiler:** typo ([befff54](https://github.com/Bearer/bearer/commit/befff54))
* **ui:** handle newly injected prop ([5449468](https://github.com/Bearer/bearer/commit/5449468))





<a name="0.64.0"></a>
# [0.64.0](https://github.com/Bearer/bearer/compare/v0.63.0...v0.64.0) (2018-09-28)


### Bug Fixes

* **ui:** add missing components ([94680b2](https://github.com/Bearer/bearer/commit/94680b2))


### Features

* **cli:** add aliases for generate comment ([d4f0d1f](https://github.com/Bearer/bearer/commit/d4f0d1f))





<a name="0.63.0"></a>
# [0.63.0](https://github.com/Bearer/bearer/compare/v0.62.1...v0.63.0) (2018-09-26)


### Features

* use signature and new endpoints for fetching the data ([e305b24](https://github.com/Bearer/bearer/commit/e305b24))
* use signature passed in the context ([bc3d7de](https://github.com/Bearer/bearer/commit/bc3d7de))





<a name="0.62.1"></a>
## [0.62.1](https://github.com/Bearer/bearer/compare/v0.62.0...v0.62.1) (2018-09-25)


### Bug Fixes

* **pack:** pass production for better browser compatibility ([#341](https://github.com/Bearer/bearer/issues/341)) ([1bb4a9f](https://github.com/Bearer/bearer/commit/1bb4a9f))





<a name="0.62.0"></a>
# [0.62.0](https://github.com/Bearer/bearer/compare/v0.61.0...v0.62.0) (2018-09-25)


### Features

* backport loading user defined data ([f8f4313](https://github.com/Bearer/bearer/commit/f8f4313))
* hardcode props for root components ([d7814ba](https://github.com/Bearer/bearer/commit/d7814ba))





<a name="0.61.0"></a>
# [0.61.0](https://github.com/Bearer/bearer/compare/v0.60.13...v0.61.0) (2018-09-18)


### Bug Fixes

* add missing summary ([3d066cc](https://github.com/Bearer/bearer/commit/3d066cc))


### Features

* add static openapi.json file to the build ([4ffd55b](https://github.com/Bearer/bearer/commit/4ffd55b))





<a name="0.60.13"></a>
## [0.60.13](https://github.com/Bearer/bearer/compare/v0.60.12...v0.60.13) (2018-09-14)

**Note:** Version bump only for package bearer-master





<a name="0.60.12"></a>
## [0.60.12](https://github.com/Bearer/bearer/compare/v0.60.11...v0.60.12) (2018-09-13)




**Note:** Version bump only for package bearer-master

<a name="0.60.11"></a>
## [0.60.11](https://github.com/Bearer/bearer/compare/v0.60.10...v0.60.11) (2018-09-13)




**Note:** Version bump only for package bearer-master

<a name="0.60.10"></a>
## [0.60.10](https://github.com/Bearer/bearer/compare/v0.60.9...v0.60.10) (2018-09-13)




**Note:** Version bump only for package bearer-master

<a name="0.60.9"></a>
## [0.60.9](https://github.com/Bearer/bearer/compare/v0.60.8...v0.60.9) (2018-09-13)


### Bug Fixes

* **cli:** update typos in push ([2c3c905](https://github.com/Bearer/bearer/commit/2c3c905))




<a name="0.60.8"></a>
## [0.60.8](https://github.com/Bearer/bearer/compare/v0.60.7...v0.60.8) (2018-09-13)


### Bug Fixes

* **cli:** remove double // ([#337](https://github.com/Bearer/bearer/issues/337)) ([5269184](https://github.com/Bearer/bearer/commit/5269184))




<a name="0.60.7"></a>
## [0.60.7](https://github.com/Bearer/bearer/compare/v0.60.6...v0.60.7) (2018-09-13)


### Bug Fixes

* **cli:** display spec.ts update only for root component ([#336](https://github.com/Bearer/bearer/issues/336)) ([85e25a9](https://github.com/Bearer/bearer/commit/85e25a9))
* **cli:** keep the scenariorc ([509bb40](https://github.com/Bearer/bearer/commit/509bb40))




<a name="0.60.6"></a>
## [0.60.6](https://github.com/Bearer/bearer/compare/v0.60.5...v0.60.6) (2018-09-12)




**Note:** Version bump only for package bearer-master

<a name="0.60.5"></a>
## [0.60.5](https://github.com/Bearer/bearer/compare/v0.60.4...v0.60.5) (2018-09-12)




**Note:** Version bump only for package bearer-master

<a name="0.60.4"></a>
## [0.60.4](https://github.com/Bearer/bearer/compare/v0.60.3...v0.60.4) (2018-09-12)


### Bug Fixes

* **cli:** update push message ([#331](https://github.com/Bearer/bearer/issues/331)) ([f2183dc](https://github.com/Bearer/bearer/commit/f2183dc))





<a name="0.60.3"></a>
## [0.60.3](https://github.com/Bearer/bearer/compare/v0.60.2...v0.60.3) (2018-09-11)


### Bug Fixes

* **cli:** make writeFile async ([#330](https://github.com/Bearer/bearer/issues/330)) ([b34ec08](https://github.com/Bearer/bearer/commit/b34ec08))





<a name="0.60.2"></a>
## [0.60.2](https://github.com/Bearer/bearer/compare/v0.60.1...v0.60.2) (2018-09-07)

**Note:** Version bump only for package bearer-master





<a name="0.60.1"></a>
## [0.60.1](https://github.com/Bearer/bearer/compare/v0.60.0...v0.60.1) (2018-09-07)

**Note:** Version bump only for package bearer-master





<a name="0.60.0"></a>
# [0.60.0](https://github.com/Bearer/bearer/compare/v0.59.1...v0.60.0) (2018-09-07)


### Bug Fixes

* **cli:** update stencil config to address www ([684a460](https://github.com/Bearer/bearer/commit/684a460))


### Features

* **cli:** rename integrationId to setupid ([3380854](https://github.com/Bearer/bearer/commit/3380854))





<a name="0.59.1"></a>
## [0.59.1](https://github.com/Bearer/bearer/compare/v0.59.0...v0.59.1) (2018-09-04)


### Bug Fixes

* **cli:** add missing tsconfig.json fil ([ab1d63e](https://github.com/Bearer/bearer/commit/ab1d63e))





<a name="0.59.0"></a>
# [0.59.0](https://github.com/Bearer/bearer/compare/v0.58.8...v0.59.0) (2018-09-03)


### Bug Fixes

* **cli:** Pack intents command ([#318](https://github.com/Bearer/bearer/issues/318)) ([ebd6b93](https://github.com/Bearer/bearer/commit/ebd6b93))
* **cli:** remove useless setTimeour ^^ ([d832fec](https://github.com/Bearer/bearer/commit/d832fec))


### Features

* **cli:** add placeholder commands ([ae5255a](https://github.com/Bearer/bearer/commit/ae5255a))
* **cli:** Build intents command ([#315](https://github.com/Bearer/bearer/issues/315)) ([95995e1](https://github.com/Bearer/bearer/commit/95995e1))
* **cli:** Build views command ([#316](https://github.com/Bearer/bearer/issues/316)) ([a5803c7](https://github.com/Bearer/bearer/commit/a5803c7))
* **cli:** pack views command ([#319](https://github.com/Bearer/bearer/issues/319)) ([fa8d791](https://github.com/Bearer/bearer/commit/fa8d791))
* **cli:** Prepare command ([#312](https://github.com/Bearer/bearer/issues/312)) ([fde2dae](https://github.com/Bearer/bearer/commit/fde2dae))
* **cli:** Push command ([#317](https://github.com/Bearer/bearer/issues/317)) ([84906e4](https://github.com/Bearer/bearer/commit/84906e4))





<a name="0.58.8"></a>
## [0.58.8](https://github.com/Bearer/bearer/compare/v0.58.7...v0.58.8) (2018-08-28)


### Bug Fixes

* **transpiler:** forgotten path ([#311](https://github.com/Bearer/bearer/issues/311)) ([37bef59](https://github.com/Bearer/bearer/commit/37bef59))





<a name="0.58.7"></a>
## [0.58.7](https://github.com/Bearer/bearer/compare/v0.58.6...v0.58.7) (2018-08-28)


### Bug Fixes

* **ui:** add prepare command ([02bcea6](https://github.com/Bearer/bearer/commit/02bcea6))
* **ui:** build ui into dist folder ([0f613f3](https://github.com/Bearer/bearer/commit/0f613f3))





<a name="0.58.6"></a>
## [0.58.6](https://github.com/Bearer/bearer/compare/v0.58.5...v0.58.6) (2018-08-28)

**Note:** Version bump only for package bearer-master





<a name="0.58.5"></a>
## [0.58.5](https://github.com/Bearer/bearer/compare/v0.57.3...v0.58.5) (2018-08-27)


### Bug Fixes

* make it consistent - force higher version ([10ff5c4](https://github.com/Bearer/bearer/commit/10ff5c4))





<a name="0.57.3"></a>
## [0.57.3](https://github.com/Bearer/bearer/compare/v0.57.2...v0.57.3) (2018-08-27)


### Bug Fixes

* **ui:** fix incorrect package location ([#307](https://github.com/Bearer/bearer/issues/307)) ([a8b29be](https://github.com/Bearer/bearer/commit/a8b29be))





<a name="0.57.2"></a>
## [0.57.2](https://github.com/Bearer/bearer/compare/v0.57.1...v0.57.2) (2018-08-27)


### Bug Fixes

* .npmrc check ([f3773a2](https://github.com/Bearer/bearer/commit/f3773a2))





<a name="0.57.0"></a>
# [0.57.0](https://github.com/Bearer/bearer/compare/v0.56.3...v0.57.0) (2018-08-27)


### Bug Fixes

* **cli:** add .scenariorc to the gitignore ([#301](https://github.com/Bearer/bearer/issues/301)) ([3551f9c](https://github.com/Bearer/bearer/commit/3551f9c))
* **cli:** add missing suffix to imports ([0494032](https://github.com/Bearer/bearer/commit/0494032))
* **cli:** be consistent ([d821628](https://github.com/Bearer/bearer/commit/d821628))


### Features

* **cli:** add jenkins integration ([bd8185c](https://github.com/Bearer/bearer/commit/bd8185c))
* change RootComponent name to role ([61428aa](https://github.com/Bearer/bearer/commit/61428aa))





<a name="0.56.3"></a>
## [0.56.3](https://github.com/Bearer/bearer/compare/v0.56.2...v0.56.3) (2018-08-23)




**Note:** Version bump only for package bearer-master

<a name="0.56.2"></a>
## [0.56.2](https://github.com/Bearer/bearer/compare/v0.56.1...v0.56.2) (2018-08-23)




**Note:** Version bump only for package bearer-master

<a name="0.56.1"></a>
## [0.56.1](https://github.com/Bearer/bearer/compare/v0.56.0...v0.56.1) (2018-08-22)


### Bug Fixes

* **cli:** generate setup fil if it does not exist ([aff55c1](https://github.com/Bearer/bearer/commit/aff55c1))




<a name="0.56.0"></a>
# [0.56.0](https://github.com/Bearer/bearer/compare/v0.55.0...v0.56.0) (2018-08-22)


### Bug Fixes

* **cli:** allow env variable package version ([66bbdb3](https://github.com/Bearer/bearer/commit/66bbdb3))
* **cli:** fix tslint formatting which was preventing iframe to load ([808ea32](https://github.com/Bearer/bearer/commit/808ea32))
* **cli:** generate full manifest file ([6ac163e](https://github.com/Bearer/bearer/commit/6ac163e))
* **cli:** make prettier work in all editors ([07816c2](https://github.com/Bearer/bearer/commit/07816c2))
* **legacy-cli:** send scenario id to the deploy command ([#290](https://github.com/Bearer/bearer/issues/290)) ([7d04ee3](https://github.com/Bearer/bearer/commit/7d04ee3))
* **transpiler:** fix missing suffix issue ([e40cfe6](https://github.com/Bearer/bearer/commit/e40cfe6))


### Features

* **cli:** add silent flag ([0f248f5](https://github.com/Bearer/bearer/commit/0f248f5))
* **cli:** generate command rewrite ([#287](https://github.com/Bearer/bearer/issues/287)) ([2e91da7](https://github.com/Bearer/bearer/commit/2e91da7))
* **cli:** generate the spec.ts in bearer new cmd ([521da58](https://github.com/Bearer/bearer/commit/521da58))
* **cli:** link command rewrite ([#284](https://github.com/Bearer/bearer/issues/284)) ([40458ca](https://github.com/Bearer/bearer/commit/40458ca))
* **transpiler:** generate manfiest.js file ([ed761d9](https://github.com/Bearer/bearer/commit/ed761d9))


### Performance Improvements

* **cli:** improve cli startup time ([#288](https://github.com/Bearer/bearer/issues/288)) ([156d565](https://github.com/Bearer/bearer/commit/156d565))




<a name="0.55.0"></a>
# [0.55.0](https://github.com/Bearer/bearer/compare/v0.54.0...v0.55.0) (2018-08-17)


### Bug Fixes

* **cli:** prevent default scenarioUuid ([#281](https://github.com/Bearer/bearer/issues/281)) ([8e8fe5a](https://github.com/Bearer/bearer/commit/8e8fe5a))


### Features

* **cli:** rewrite login command ([#283](https://github.com/Bearer/bearer/issues/283)) ([0ae055f](https://github.com/Bearer/bearer/commit/0ae055f))
* **transpiler:** scope bearer components ([#279](https://github.com/Bearer/bearer/issues/279)) ([7b9453d](https://github.com/Bearer/bearer/commit/7b9453d))




<a name="0.54.0"></a>
# [0.54.0](https://github.com/Bearer/bearer/compare/v0.53.1...v0.54.0) (2018-08-16)


### Bug Fixes

* **cli:** Better deploy error message ([#275](https://github.com/Bearer/bearer/issues/275)) ([70b61af](https://github.com/Bearer/bearer/commit/70b61af))
* **webserver:** handle not existing intent gracefully ([#277](https://github.com/Bearer/bearer/issues/277)) ([6af021a](https://github.com/Bearer/bearer/commit/6af021a))


### Features

* **cli:** new commnad rewrite base ([#276](https://github.com/Bearer/bearer/issues/276)) ([f58d778](https://github.com/Bearer/bearer/commit/f58d778))




<a name="0.53.0"></a>
# [0.53.0](https://github.com/Bearer/bearer/compare/v0.52.0...v0.53.0) (2018-08-14)


### Bug Fixes

* **cli:** rollback removed features ([#273](https://github.com/Bearer/bearer/issues/273)) ([c81d929](https://github.com/Bearer/bearer/commit/c81d929))
* **core:** change LICENSE field to MIT ([c133805](https://github.com/Bearer/bearer/commit/c133805))
* **templates:** add navigator-auth-screen only when oauth2 ([#269](https://github.com/Bearer/bearer/issues/269)) ([785276c](https://github.com/Bearer/bearer/commit/785276c))
* **transpiler:** log less ([#267](https://github.com/Bearer/bearer/issues/267)) ([1faea1e](https://github.com/Bearer/bearer/commit/1faea1e))
* **ui:** change LICENSE to MIT ([d90917d](https://github.com/Bearer/bearer/commit/d90917d))


### Features

* **intents:** give access to dbclient ([#270](https://github.com/Bearer/bearer/issues/270)) ([fbf284e](https://github.com/Bearer/bearer/commit/fbf284e))




<a name="0.51.0"></a>
# [0.51.0](https://github.com/Bearer/bearer/compare/v0.50.0...v0.51.0) (2018-08-13)


### Bug Fixes

* **cli:** ask intenttype before name ([df52ae4](https://github.com/Bearer/bearer/commit/df52ae4))
* **cli:** ensure error message is returned only when exists ([f0dae67](https://github.com/Bearer/bearer/commit/f0dae67))
* **cli:** generate first screen ([b215a3f](https://github.com/Bearer/bearer/commit/b215a3f))
* **ui:** make setup more robust ([3c160fc](https://github.com/Bearer/bearer/commit/3c160fc))


### Features

* **cli:** move config.dev.js to example ([0d1d889](https://github.com/Bearer/bearer/commit/0d1d889))




<a name="0.50.0"></a>
# [0.50.0](https://github.com/Bearer/bearer/compare/v0.48.5...v0.50.0) (2018-08-13)


### Bug Fixes

* **bearer:** include typescript to fix missing typescript issue ([342b0b5](https://github.com/Bearer/bearer/commit/342b0b5))
* **cli:** fix serviceClient issue ([42b3c40](https://github.com/Bearer/bearer/commit/42b3c40))
* **cli:** generate component fixes ([ce8a6d4](https://github.com/Bearer/bearer/commit/ce8a6d4))
* **cli:** give button text a meaning ([1e5c78e](https://github.com/Bearer/bearer/commit/1e5c78e))
* **cli:** give more meaning to template variable names ([d25d0b6](https://github.com/Bearer/bearer/commit/d25d0b6))
* **cli:** intents variables generation ([1535f77](https://github.com/Bearer/bearer/commit/1535f77))
* **cli:** remove cache to ensure shrinkwrap file is properly generated ([e318925](https://github.com/Bearer/bearer/commit/e318925))
* **templates:** wrong callback call ([aa99089](https://github.com/Bearer/bearer/commit/aa99089))


### Features

* **cli:** proxy the notifications through package manager ([514e8aa](https://github.com/Bearer/bearer/commit/514e8aa))




<a name="0.49.0"></a>
# [0.49.0](https://github.com/Bearer/bearer/compare/v0.48.5...v0.49.0) (2018-08-12)


### Bug Fixes

* **bearer:** include typescript to fix missing typescript issue ([342b0b5](https://github.com/Bearer/bearer/commit/342b0b5))
* **cli:** fix serviceClient issue ([42b3c40](https://github.com/Bearer/bearer/commit/42b3c40))


### Features

* **cli:** proxy the notifications through package manager ([514e8aa](https://github.com/Bearer/bearer/commit/514e8aa))




<a name="0.48.5"></a>
## [0.48.5](https://github.com/Bearer/bearer/compare/v0.48.4...v0.48.5) (2018-08-12)




**Note:** Version bump only for package bearer-master

<a name="0.48.4"></a>
## [0.48.4](https://github.com/Bearer/bearer/compare/v0.48.3...v0.48.4) (2018-08-10)

**Note:** Version bump only for package bearer-master





<a name="0.48.3"></a>
## [0.48.3](https://github.com/Bearer/bearer/compare/v0.48.2...v0.48.3) (2018-08-10)

**Note:** Version bump only for package bearer-master





<a name="0.48.2"></a>
## [0.48.2](https://github.com/Bearer/bearer/compare/v0.48.1...v0.48.2) (2018-08-10)

**Note:** Version bump only for package bearer-master





<a name="0.48.1"></a>
## [0.48.1](https://github.com/Bearer/bearer/compare/v0.48.0...v0.48.1) (2018-08-10)




**Note:** Version bump only for package bearer-master

<a name="0.47.5"></a>
## [0.47.5](https://github.com/Bearer/bearer/compare/v0.47.4...v0.47.5) (2018-08-10)




**Note:** Version bump only for package bearer-master

<a name="0.46.1"></a>

## [0.46.1](https://github.com/Bearer/bearer/compare/v0.46.0...v0.46.1) (2018-08-10)

### Bug Fixes

- **cli:** disable plugins ([9013327](https://github.com/Bearer/bearer/commit/9013327))
- **cli:** remove autocomplete/not found/ not found plugins ([ce3977a](https://github.com/Bearer/bearer/commit/ce3977a))

<a name="0.45.2"></a>

## [0.45.2](https://github.com/Bearer/bearer/compare/v0.45.1...v0.45.2) (2018-08-10)

### Bug Fixes

- **cli:** do not rely on node_modules or bin anymore ([#256](https://github.com/Bearer/bearer/issues/256)) ([d839c42](https://github.com/Bearer/bearer/commit/d839c42))

<a name="0.45.1"></a>

## [0.45.1](https://github.com/Bearer/bearer/compare/v0.45.0...v0.45.1) (2018-08-10)

### Bug Fixes

- **cli:** rely on binary ([20b2c99](https://github.com/Bearer/bearer/commit/20b2c99))

<a name="0.43.1"></a>

## [0.43.1](https://github.com/Bearer/bearer/compare/v0.43.0...v0.43.1) (2018-08-10)

**Note:** Version bump only for package bearer

<a name="0.42.0"></a>

# [0.42.0](https://github.com/Bearer/bearer/compare/v0.41.9...v0.42.0) (2018-08-10)

### Features

- **cli:** add did you mean plugin ([68dfb25](https://github.com/Bearer/bearer/commit/68dfb25))
- **cli:** bring autocomplete to cli ([518252a](https://github.com/Bearer/bearer/commit/518252a))

<a name="0.41.9"></a>

## [0.41.9](https://github.com/Bearer/bearer/compare/v0.41.8...v0.41.9) (2018-08-09)

### Bug Fixes

- **cli:** empty build dir instead of src one ([32f1016](https://github.com/Bearer/bearer/commit/32f1016))

<a name="0.41.8"></a>

## [0.41.8](https://github.com/Bearer/bearer/compare/v0.41.7...v0.41.8) (2018-08-09)

### Bug Fixes

- **cli:** prevent cli to proceed if authentication failed ([b7f4f69](https://github.com/Bearer/bearer/commit/b7f4f69))
- **cli:** remove duplicate identifier ([b2b92ae](https://github.com/Bearer/bearer/commit/b2b92ae))
- **legacy-cli:** add setupId in template setup ([a0d52a5](https://github.com/Bearer/bearer/commit/a0d52a5))
- **ui:** test setupId in render ([7d96e14](https://github.com/Bearer/bearer/commit/7d96e14))

<a name="0.41.7"></a>

## [0.41.7](https://github.com/Bearer/bearer/compare/v0.41.6...v0.41.7) (2018-08-09)

### Bug Fixes

- **cli:** deploy with the correct namespace name ([a46b520](https://github.com/Bearer/bearer/commit/a46b520))
- **cli:** fix path invalidation ([2ec99d0](https://github.com/Bearer/bearer/commit/2ec99d0))

<a name="0.41.6"></a>

## [0.41.6](https://github.com/Bearer/bearer/compare/v0.41.5...v0.41.6) (2018-08-09)

### Bug Fixes

- **cli:** start was missing authorization host ([#248](https://github.com/Bearer/bearer/issues/248)) ([7b9edbf](https://github.com/Bearer/bearer/commit/7b9edbf))
- **cli:** use localhost as a tag name on start ([c0e2d95](https://github.com/Bearer/bearer/commit/c0e2d95))
- **ui:** add bearer style to sumbit button ([30d7c1b](https://github.com/Bearer/bearer/commit/30d7c1b))

<a name="0.41.5"></a>

## [0.41.5](https://github.com/Bearer/bearer/compare/v0.41.4...v0.41.5) (2018-08-09)

### Bug Fixes

- **cli:** ensure .bearer dir exists ([793e19e](https://github.com/Bearer/bearer/commit/793e19e))
- **cli:** pass tag name to stencil ([4e3d46f](https://github.com/Bearer/bearer/commit/4e3d46f))
- **ui:** stencil watch decorator does not work on state anymore ([113905b](https://github.com/Bearer/bearer/commit/113905b))

<a name="0.41.4"></a>

## [0.41.4](https://github.com/Bearer/bearer/compare/v0.41.3...v0.41.4) (2018-08-08)

### Bug Fixes

- **intents:** it looks like lambda body is already an object ([72487e9](https://github.com/Bearer/bearer/commit/72487e9))

<a name="0.41.3"></a>

## [0.41.3](https://github.com/Bearer/bearer/compare/v0.41.2...v0.41.3) (2018-08-08)

**Note:** Version bump only for package bearer

<a name="0.41.2"></a>

## [0.41.2](https://github.com/Bearer/bearer/compare/v0.41.1...v0.41.2) (2018-08-08)

### Bug Fixes

- **intents:** fallback body ([077dc7c](https://github.com/Bearer/bearer/commit/077dc7c))

<a name="0.41.1"></a>

## [0.41.1](https://github.com/Bearer/bearer/compare/v0.41.0...v0.41.1) (2018-08-08)

**Note:** Version bump only for package bearer

<a name="0.41.0"></a>

# [0.41.0](https://github.com/Bearer/bearer/compare/v0.40.0...v0.41.0) (2018-08-08)

### Bug Fixes

- **cli:** better logging when webpack is failing ([#239](https://github.com/Bearer/bearer/issues/239)) ([189296e](https://github.com/Bearer/bearer/commit/189296e))
- **intents:** mark data as optional for retrievestate intents ([154248a](https://github.com/Bearer/bearer/commit/154248a))
- **intents:** provide correct function signatures ([7125157](https://github.com/Bearer/bearer/commit/7125157))
- **intents:** typo ([4ee09c9](https://github.com/Bearer/bearer/commit/4ee09c9))
- **webserver:** unstringify ([5c7cb5f](https://github.com/Bearer/bearer/commit/5c7cb5f))

### Features

- **cli:** import client only when required ([#238](https://github.com/Bearer/bearer/issues/238)) ([2d95257](https://github.com/Bearer/bearer/commit/2d95257))

<a name="0.39.0"></a>

# [0.39.0](https://github.com/Bearer/bearer/compare/v0.38.1...v0.39.0) (2018-08-08)

### Bug Fixes

- **cli:** emit when setup is not required ([53269d1](https://github.com/Bearer/bearer/commit/53269d1))
- **cli:** let oclif deal with new versions ([#231](https://github.com/Bearer/bearer/issues/231)) ([4dd50f8](https://github.com/Bearer/bearer/commit/4dd50f8))
- **cli:** update the cdn host for production ([76348d2](https://github.com/Bearer/bearer/commit/76348d2))
- **core:** post robots now use correct url ([#235](https://github.com/Bearer/bearer/issues/235)) ([49abd5b](https://github.com/Bearer/bearer/commit/49abd5b))
- **core:** reject intent promise if error returned ([#233](https://github.com/Bearer/bearer/issues/233)) ([d1c36b9](https://github.com/Bearer/bearer/commit/d1c36b9))

### Features

- **cli:** remove auth navigator screen when not needed ([d259d0b](https://github.com/Bearer/bearer/commit/d259d0b))
- **core:** everything as post ([#234](https://github.com/Bearer/bearer/issues/234)) ([ae67130](https://github.com/Bearer/bearer/commit/ae67130))

<a name="0.37.2"></a>

## [0.37.2](https://github.com/Bearer/bearer/compare/v0.37.0...v0.37.2) (2018-08-08)

**Note:** Version bump only for package bearer

<a name="0.37.1"></a>

## [0.37.1](https://github.com/Bearer/bearer/compare/v0.37.0...v0.37.1) (2018-08-08)

**Note:** Version bump only for package bearer

<a name="0.37.0"></a>

# [0.37.0](https://github.com/Bearer/bearer/compare/v0.36.3...v0.37.0) (2018-08-08)

### Bug Fixes

- **cli:** remove bearer command reference ([d2e8f4b](https://github.com/Bearer/bearer/commit/d2e8f4b))
- **cli:** remove old bearer binary ([#223](https://github.com/Bearer/bearer/issues/223)) ([746553b](https://github.com/Bearer/bearer/commit/746553b))
- **transpiler:** fix tag starting with number ([682b651](https://github.com/Bearer/bearer/commit/682b651))
- **transpiler:** updated test suite ([5187fa5](https://github.com/Bearer/bearer/commit/5187fa5))

### Features

- **cli:** add new bearer cli OCLIF base ([#221](https://github.com/Bearer/bearer/issues/221)) ([d18abb1](https://github.com/Bearer/bearer/commit/d18abb1))

<a name="0.36.3"></a>

## [0.36.3](https://github.com/Bearer/bearer/compare/v0.36.2...v0.36.3) (2018-08-07)

**Note:** Version bump only for package undefined

<a name="0.36.2"></a>

## [0.36.2](https://github.com/Bearer/bearer/compare/v0.36.1...v0.36.2) (2018-08-07)

**Note:** Version bump only for package undefined

<a name="0.35.0"></a>

# [0.35.0](https://github.com/Bearer/bearer/compare/v0.34.2...v0.35.0) (2018-08-07)

### Bug Fixes

- **cli:** remove useless div ([e0c248a](https://github.com/Bearer/bearer/commit/e0c248a))
- **cli:** use correct intent type name ([2d5ffa4](https://github.com/Bearer/bearer/commit/2d5ffa4))
- **ui:** regenerate component definition ([3392a93](https://github.com/Bearer/bearer/commit/3392a93))
- **ui:** reset component style correctly ([56e3749](https://github.com/Bearer/bearer/commit/56e3749))

### Features

- **cli:** add default empty setup.css ([f6c4e20](https://github.com/Bearer/bearer/commit/f6c4e20))
- **cli:** clean up auth.config for oauth2 ([#207](https://github.com/Bearer/bearer/issues/207)) ([4b45e2e](https://github.com/Bearer/bearer/commit/4b45e2e))
- upgrade stencil required version ([d797a67](https://github.com/Bearer/bearer/commit/d797a67))

<a name="0.34.0"></a>

# [0.34.0](https://github.com/Bearer/bearer/compare/v0.33.0...v0.34.0) (2018-08-03)

### Bug Fixes

- preserve case for first letter ([ba13c81](https://github.com/Bearer/bearer/commit/ba13c81))
- **cli:** outdated paths ([d8b79c8](https://github.com/Bearer/bearer/commit/d8b79c8))
- **cli:** packagejson template ([5df1c5f](https://github.com/Bearer/bearer/commit/5df1c5f))
- **cli:** remove unnecessary dependencies ([aa483ff](https://github.com/Bearer/bearer/commit/aa483ff))
- **cli:** replace outdated path ([2550d84](https://github.com/Bearer/bearer/commit/2550d84))
- **transpiler:** no unused method ([c6dc79a](https://github.com/Bearer/bearer/commit/c6dc79a))
- **transpiler:** prevent double watcher transpilation ([a6f164c](https://github.com/Bearer/bearer/commit/a6f164c))
- **transpiler:** state injector update properly lifecycle methods ([d12a2c1](https://github.com/Bearer/bearer/commit/d12a2c1))
- **transpiler:** wrong filename ([1551a33](https://github.com/Bearer/bearer/commit/1551a33))

### Features

- **cli:** add support for npm when yarn is not available ([78ed3fb](https://github.com/Bearer/bearer/commit/78ed3fb))
- **core:** add support for shadow in RootComponent ([#205](https://github.com/Bearer/bearer/issues/205)) ([4529771](https://github.com/Bearer/bearer/commit/4529771))

<a name="0.33.0"></a>

# [0.33.0](https://github.com/Bearer/bearer/compare/v0.32.0...v0.33.0) (2018-08-03)

### Bug Fixes

- **cli:** remove node-jq dependency ([57bc84f](https://github.com/Bearer/bearer/commit/57bc84f))
- **core:** pass boolean to maybeInitializedSession promise ([76b1f67](https://github.com/Bearer/bearer/commit/76b1f67))
- **transpiler:** remove stupid doublon ([dcbde6d](https://github.com/Bearer/bearer/commit/dcbde6d))
- **ui:** prevent extra call when scenario is completed ([#196](https://github.com/Bearer/bearer/issues/196)) ([e1d5b8c](https://github.com/Bearer/bearer/commit/e1d5b8c))
- **ui:** update index ([ce7f0bc](https://github.com/Bearer/bearer/commit/ce7f0bc))
- jest upgrade ([25bdb02](https://github.com/Bearer/bearer/commit/25bdb02))

### Features

- **transpiler:** add metata parser ([6a066b2](https://github.com/Bearer/bearer/commit/6a066b2))
- **transpiler:** add metata parser ([423ff50](https://github.com/Bearer/bearer/commit/423ff50))
- **ui:** add disabled button prop ([#198](https://github.com/Bearer/bearer/issues/198)) ([765ef1d](https://github.com/Bearer/bearer/commit/765ef1d))
- **ui:** introduce bearer-authorized component ([69d5aac](https://github.com/Bearer/bearer/commit/69d5aac))
- **ui:** provide a way to authenticate people through a method ([5996faf](https://github.com/Bearer/bearer/commit/5996faf))

<a name="0.32.0"></a>

# [0.32.0](https://github.com/Bearer/bearer/compare/v0.31.0...v0.32.0) (2018-08-02)

### Bug Fixes

- **transpiler:** update test suite ([7da4def](https://github.com/Bearer/bearer/commit/7da4def))

### Features

- **core:** simplify decorator use ([#193](https://github.com/Bearer/bearer/issues/193)) ([fa3facc](https://github.com/Bearer/bearer/commit/fa3facc))

<a name="0.31.0"></a>

# [0.31.0](https://github.com/Bearer/bearer/compare/v0.30.0...v0.31.0) (2018-08-01)

### Bug Fixes

- **transpiler:** Fix root component ([#191](https://github.com/Bearer/bearer/issues/191)) ([ac68a8d](https://github.com/Bearer/bearer/commit/ac68a8d))

### Features

- **cli:** better local development index.html style ([35ee985](https://github.com/Bearer/bearer/commit/35ee985))
- **cli:** update link command to only accept one arg ([#188](https://github.com/Bearer/bearer/issues/188)) ([7b024de](https://github.com/Bearer/bearer/commit/7b024de))

<a name="0.30.0"></a>

# [0.30.0](https://github.com/Bearer/bearer/compare/v0.29.0...v0.30.0) (2018-08-01)

### Features

- **core:** trigger bearer:StateSaved event ([3d672c1](https://github.com/Bearer/bearer/commit/3d672c1))

<a name="0.29.0"></a>

# [0.29.0](https://github.com/Bearer/bearer/compare/v0.28.0...v0.29.0) (2018-08-01)

### Bug Fixes

- add .keep to pass tests ([e298811](https://github.com/Bearer/bearer/commit/e298811))
- add snapshot ([bb0f1e2](https://github.com/Bearer/bearer/commit/bb0f1e2))
- empty css is now using pascal case ([a996af9](https://github.com/Bearer/bearer/commit/a996af9))
- remove old comments ([21834eb](https://github.com/Bearer/bearer/commit/21834eb))

### Features

- add generate root group ([f365675](https://github.com/Bearer/bearer/commit/f365675))
- add root component transformer ([a710da2](https://github.com/Bearer/bearer/commit/a710da2))
- add root decorator ([0754737](https://github.com/Bearer/bearer/commit/0754737))
- add root group option ([f099802](https://github.com/Bearer/bearer/commit/f099802))
- add skel root component transformer ([a9fd496](https://github.com/Bearer/bearer/commit/a9fd496))
- **cli:** switch to dev portal ([#184](https://github.com/Bearer/bearer/issues/184)) ([506ac90](https://github.com/Bearer/bearer/commit/506ac90))
- **transpiler:** inject el into component ([7f00a76](https://github.com/Bearer/bearer/commit/7f00a76))

<a name="0.28.0"></a>

# [0.28.0](https://github.com/Bearer/bearer/compare/v0.27.4...v0.28.0) (2018-07-31)

### Bug Fixes

- **core:** add better typing ([bf2f477](https://github.com/Bearer/bearer/commit/bf2f477))
- **ui:** update typing form core update ([841318e](https://github.com/Bearer/bearer/commit/841318e))

### Features

- **cli:** add jq to invoke ([#175](https://github.com/Bearer/bearer/issues/175)) ([582180d](https://github.com/Bearer/bearer/commit/582180d))

<a name="0.27.3"></a>

## [0.27.3](https://github.com/Bearer/bearer/compare/v0.27.2...v0.27.3) (2018-07-30)

**Note:** Version bump only for package undefined

<a name="0.27.2"></a>

## [0.27.2](https://github.com/Bearer/bearer/compare/v0.27.1...v0.27.2) (2018-07-30)

### Bug Fixes

- **cli:** developer portal query ([#174](https://github.com/Bearer/bearer/issues/174)) ([759001e](https://github.com/Bearer/bearer/commit/759001e))

<a name="0.27.1"></a>

## [0.27.1](https://github.com/Bearer/bearer/compare/v0.27.0...v0.27.1) (2018-07-27)

**Note:** Version bump only for package undefined

<a name="0.27.0"></a>

# [0.27.0](https://github.com/Bearer/bearer/compare/v0.26.0...v0.27.0) (2018-07-27)

### Bug Fixes

- change terminology ([724a69d](https://github.com/Bearer/bearer/commit/724a69d))

### Features

- **cli:** add warning is missing or incorrect ([#167](https://github.com/Bearer/bearer/issues/167)) ([6af728c](https://github.com/Bearer/bearer/commit/6af728c))
- **cli:** rename invoke to run ([#172](https://github.com/Bearer/bearer/issues/172)) ([3a440a6](https://github.com/Bearer/bearer/commit/3a440a6))

<a name="0.26.0"></a>

# [0.26.0](https://github.com/Bearer/bearer/compare/v0.25.1...v0.26.0) (2018-07-26)

### Bug Fixes

- **intents:** remove async ([#169](https://github.com/Bearer/bearer/issues/169)) ([e1e10b3](https://github.com/Bearer/bearer/commit/e1e10b3))

### Features

- **cli:** add invoke command to run intent locally ([#165](https://github.com/Bearer/bearer/issues/165)) ([94a909d](https://github.com/Bearer/bearer/commit/94a909d))
- **cli:** prepare production ([#170](https://github.com/Bearer/bearer/issues/170)) ([2a6aa2e](https://github.com/Bearer/bearer/commit/2a6aa2e))
- **transpiler:** navigator screen rewrite transformer introduction ([e6e1d14](https://github.com/Bearer/bearer/commit/e6e1d14))
- **transpiler:** transformer slot to renderFunc ([7a1b017](https://github.com/Bearer/bearer/commit/7a1b017))

<a name="0.25.1"></a>

## [0.25.1](https://github.com/Bearer/bearer/compare/v0.25.0...v0.25.1) (2018-07-26)

### Bug Fixes

- **cli:** remove hardcoded wrong path ([3e08d3e](https://github.com/Bearer/bearer/commit/3e08d3e))

<a name="0.25.0"></a>

# [0.25.0](https://github.com/Bearer/bearer/compare/v0.24.1...v0.25.0) (2018-07-26)

### Bug Fixes

- remove lerna 3 option ([64312d3](https://github.com/Bearer/bearer/commit/64312d3))
- remove useless console.log ([a40fffd](https://github.com/Bearer/bearer/commit/a40fffd))
- remove useless enum item ([e5358be](https://github.com/Bearer/bearer/commit/e5358be))

### Features

- **ui:** complete prop on navigator ([8c39006](https://github.com/Bearer/bearer/commit/8c39006))
- add new options to generate command ([ed9ec03](https://github.com/Bearer/bearer/commit/ed9ec03))
- add optional name for generate command ([d4300e0](https://github.com/Bearer/bearer/commit/d4300e0))

<a name="0.24.1"></a>

## [0.24.1](https://github.com/Bearer/bearer/compare/v0.24.0...v0.24.1) (2018-07-25)

### Bug Fixes

- **ui:** prevent unnecessary screen rendering ([cb25e0b](https://github.com/Bearer/bearer/commit/cb25e0b))

<a name="0.24.0"></a>

# [0.24.0](https://github.com/Bearer/bearer/compare/v0.23.2...v0.24.0) (2018-07-25)

### Bug Fixes

- adapt template screens ([c4d29f0](https://github.com/Bearer/bearer/commit/c4d29f0))
- **cli:** buildartifact dist moved to .bearer ([#156](https://github.com/Bearer/bearer/issues/156)) ([fe16f2a](https://github.com/Bearer/bearer/commit/fe16f2a))
- **cli:** file renaming was missing ([d9faf0c](https://github.com/Bearer/bearer/commit/d9faf0c))
- **cli:** handle promise failures ([6d98d71](https://github.com/Bearer/bearer/commit/6d98d71))
- **cli:** No intent issue ([62e6f9a](https://github.com/Bearer/bearer/commit/62e6f9a))
- **cli:** rename screen folder to view ([9211568](https://github.com/Bearer/bearer/commit/9211568))
- add space between brackets ([9c3f671](https://github.com/Bearer/bearer/commit/9c3f671))
- clean build ([b638e71](https://github.com/Bearer/bearer/commit/b638e71))
- remove useless bits ([af89ec9](https://github.com/Bearer/bearer/commit/af89ec9))
- **ui:** get rid of conditional rendering within button popover ([b624c8b](https://github.com/Bearer/bearer/commit/b624c8b))

### Features

- **cli:** add bearer keep ([6eb8a2e](https://github.com/Bearer/bearer/commit/6eb8a2e))
- **cli:** add options shortcut ([36344ec](https://github.com/Bearer/bearer/commit/36344ec))
- **cli:** allow screens/intents deploy only ([98bdd6a](https://github.com/Bearer/bearer/commit/98bdd6a))
- **cli:** let deploy use views ([0404f10](https://github.com/Bearer/bearer/commit/0404f10))
- **cli:** make local deve working ([40c4705](https://github.com/Bearer/bearer/commit/40c4705))
- **cli:** rename component template name to feature ([63c0374](https://github.com/Bearer/bearer/commit/63c0374))
- **CLI:** Rename dev.config.js to config.dev.js ([ec04fca](https://github.com/Bearer/bearer/commit/ec04fca))
- **generator:** Generate setup and config only when required ([#129](https://github.com/Bearer/bearer/issues/129)) ([7091652](https://github.com/Bearer/bearer/commit/7091652)), closes [#148](https://github.com/Bearer/bearer/issues/148)
- **generator:** improve scenario name resilience ([#154](https://github.com/Bearer/bearer/issues/154)) ([c0099f6](https://github.com/Bearer/bearer/commit/c0099f6))
- **generator:** Separate TContext per AuthType ([#152](https://github.com/Bearer/bearer/issues/152)) ([b194c23](https://github.com/Bearer/bearer/commit/b194c23))
- **ui:** change default location ([9ea6cc2](https://github.com/Bearer/bearer/commit/9ea6cc2))
- add auth screen to oauth2 setup ([f93d5a9](https://github.com/Bearer/bearer/commit/f93d5a9))
- add callback to local auth ([9cf65b1](https://github.com/Bearer/bearer/commit/9cf65b1))
- rename screen to view ([57da6e5](https://github.com/Bearer/bearer/commit/57da6e5))

<a name="0.23.2"></a>

## [0.23.2](https://github.com/Bearer/bearer/compare/v0.23.1...v0.23.2) (2018-07-24)

### Bug Fixes

- **cli:** remove useless option ([#147](https://github.com/Bearer/bearer/issues/147)) ([c9acd2e](https://github.com/Bearer/bearer/commit/c9acd2e))

<a name="0.23.1"></a>

## [0.23.1](https://github.com/Bearer/bearer/compare/v0.23.0...v0.23.1) (2018-07-24)

### Bug Fixes

- **cli:** Prevents cli to hang on any command ([4ae6233](https://github.com/Bearer/bearer/commit/4ae6233))
- **cli:** remove localStorage for now ([#144](https://github.com/Bearer/bearer/issues/144)) ([d7a074c](https://github.com/Bearer/bearer/commit/d7a074c))

<a name="0.22.4"></a>

## [0.22.4](https://github.com/Bearer/bearer/compare/v0.22.2...v0.22.4) (2018-07-23)

### Bug Fixes

- **core:** handle non consistent payload ([dec01d5](https://github.com/Bearer/bearer/commit/dec01d5))

<a name="0.22.3"></a>

## [0.22.3](https://github.com/Bearer/bearer/compare/v0.22.3-0...v0.22.3) (2018-07-23)

### Bug Fixes

- **core:** handle non consistent payload ([2d409bb](https://github.com/Bearer/bearer/commit/2d409bb))

<a name="0.22.2"></a>

## [0.22.2](https://github.com/Bearer/bearer/compare/v0.22.1...v0.22.2) (2018-07-23)

### Bug Fixes

- **cli:** upgrade intents dependency ([#141](https://github.com/Bearer/bearer/issues/141)) ([0a575da](https://github.com/Bearer/bearer/commit/0a575da))

<a name="0.21.2"></a>

## [0.21.2](https://github.com/Bearer/bearer/compare/v0.21.1...v0.21.2) (2018-07-23)

### Bug Fixes

- **ui:** backport correct behaviour form previous implementation ([#137](https://github.com/Bearer/bearer/issues/137)) ([9f48b8b](https://github.com/Bearer/bearer/commit/9f48b8b))

<a name="0.21.1"></a>

## [0.21.1](https://github.com/Bearer/bearer/compare/v0.21.0...v0.21.1) (2018-07-23)

### Bug Fixes

- **cli:** prevent cli to fail if used outside of a scenario ([#136](https://github.com/Bearer/bearer/issues/136)) ([09dab53](https://github.com/Bearer/bearer/commit/09dab53))

<a name="0.21.0"></a>

# [0.21.0](https://github.com/Bearer/bearer/compare/v0.20.2...v0.21.0) (2018-07-23)

### Bug Fixes

- **cli:** Fix package manager returning to early ([#123](https://github.com/Bearer/bearer/issues/123)) ([7e68dca](https://github.com/Bearer/bearer/commit/7e68dca))
- **cli:** Update {{intentName}}.ts ([#132](https://github.com/Bearer/bearer/issues/132)) ([05effbb](https://github.com/Bearer/bearer/commit/05effbb))
- **ui:** Remove conditional default ([a4504da](https://github.com/Bearer/bearer/commit/a4504da))

### Features

- **cli:** Add basic auth in the init command ([#128](https://github.com/Bearer/bearer/issues/128)) ([95a2adc](https://github.com/Bearer/bearer/commit/95a2adc))

<a name="0.20.2"></a>

## [0.20.2](https://github.com/Bearer/bearer/compare/v0.20.1...v0.20.2) (2018-07-20)

**Note:** Version bump only for package undefined

<a name="0.20.1"></a>

## [0.20.1](https://github.com/Bearer/bearer/compare/v0.20.0...v0.20.1) (2018-07-20)

**Note:** Version bump only for package undefined

<a name="0.20.0"></a>

# [0.20.0](https://github.com/Bearer/bearer/compare/v0.19.1...v0.20.0) (2018-07-20)

### Features

- **templates:** Add BasicAuth and update Client ([#120](https://github.com/Bearer/bearer/issues/120)) ([afc0192](https://github.com/Bearer/bearer/commit/afc0192))

<a name="0.19.0"></a>

# [0.19.0](https://github.com/Bearer/bearer/compare/v0.18.0...v0.19.0) (2018-07-19)

### Features

- **generator:** Improve Intent generator ([#113](https://github.com/Bearer/bearer/issues/113)) ([8f86f46](https://github.com/Bearer/bearer/commit/8f86f46))

<a name="0.18.0"></a>

# [0.18.0](https://github.com/Bearer/bearer/compare/v0.17.5...v0.18.0) (2018-07-19)

### Features

- **generator:** improves client.ts ([#107](https://github.com/Bearer/bearer/issues/107)) ([5c5ffcc](https://github.com/Bearer/bearer/commit/5c5ffcc))
- **templates:** renames devIntents.config to dev.config ([#108](https://github.com/Bearer/bearer/issues/108)) ([f5d617f](https://github.com/Bearer/bearer/commit/f5d617f))

<a name="0.17.5"></a>

## [0.17.5](https://github.com/Bearer/bearer/compare/v0.17.4...v0.17.5) (2018-07-18)

**Note:** Version bump only for package undefined

<a name="0.17.4"></a>

## [0.17.4](https://github.com/Bearer/bearer/compare/v0.17.3...v0.17.4) (2018-07-18)

### Bug Fixes

- **local-dev-server:** allow override ([#112](https://github.com/Bearer/bearer/issues/112)) ([f1dd000](https://github.com/Bearer/bearer/commit/f1dd000))
- **transpiler:** typo ([ea1ed75](https://github.com/Bearer/bearer/commit/ea1ed75))
