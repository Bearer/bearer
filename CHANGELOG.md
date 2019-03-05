# Change Log

All notable changes to this project will be documented in this file.
See [Conventional Commits](https://conventionalcommits.org) for commit guidelines.

# [0.100.0](https://github.com/Bearer/bearer/compare/v0.99.2...v0.100.0) (2019-03-04)


### Features

* **react:** Add react connect component ([#554](https://github.com/Bearer/bearer/issues/554)) ([9507947](https://github.com/Bearer/bearer/commit/9507947))





## [0.99.2](https://github.com/Bearer/bearer/compare/v0.99.1...v0.99.2) (2019-03-02)

**Note:** Version bump only for package bearer-master





## [0.99.1](https://github.com/Bearer/bearer/compare/v0.99.0...v0.99.1) (2019-03-01)


### Bug Fixes

* add missing type field to packages ([#551](https://github.com/Bearer/bearer/issues/551)) ([3d96f83](https://github.com/Bearer/bearer/commit/3d96f83))
* **js:** fix incorrect lookup ([#552](https://github.com/Bearer/bearer/issues/552)) ([a9bc997](https://github.com/Bearer/bearer/commit/a9bc997))





# [0.99.0](https://github.com/Bearer/bearer/compare/v0.98.0...v0.99.0) (2019-02-28)


### Bug Fixes

* **logger:** issue with multiple logger export ([#550](https://github.com/Bearer/bearer/issues/550)) ([5f5281d](https://github.com/Bearer/bearer/commit/5f5281d))


### Features

* **intents:** add Oauth1 typing ([#525](https://github.com/Bearer/bearer/issues/525)) ([62ce12a](https://github.com/Bearer/bearer/commit/62ce12a))





# [0.98.0](https://github.com/Bearer/bearer/compare/v0.97.4...v0.98.0) (2019-02-28)


### Features

* **js:** allow bearer to receive options ([#549](https://github.com/Bearer/bearer/issues/549)) ([96b4977](https://github.com/Bearer/bearer/commit/96b4977))





## [0.97.4](https://github.com/Bearer/bearer/compare/v0.97.3...v0.97.4) (2019-02-27)


### Bug Fixes

* **cli:** don't rely on weird tags for openapi-generator ([#548](https://github.com/Bearer/bearer/issues/548)) ([89b544d](https://github.com/Bearer/bearer/commit/89b544d))
* **transpiler:** virtually scope events ([42d054c](https://github.com/Bearer/bearer/commit/42d054c))





## [0.97.3](https://github.com/Bearer/bearer/compare/v0.97.2...v0.97.3) (2019-02-27)


### Bug Fixes

* **openapi-generator:** let lerna be aware about the package ([#547](https://github.com/Bearer/bearer/issues/547)) ([1a5b4b0](https://github.com/Bearer/bearer/commit/1a5b4b0))





## [0.97.2](https://github.com/Bearer/bearer/compare/v0.97.1...v0.97.2) (2019-02-27)


### Bug Fixes

* **cli:** introduce openapi generator ([#538](https://github.com/Bearer/bearer/issues/538)) ([0f8c78b](https://github.com/Bearer/bearer/commit/0f8c78b))
* make the logger truly work in the browser ([cc05f9d](https://github.com/Bearer/bearer/commit/cc05f9d))
* **openapi-generator:** fix the intent order ([#545](https://github.com/Bearer/bearer/issues/545)) ([14c1087](https://github.com/Bearer/bearer/commit/14c1087))
* **ui:** CSS customization on buttons ([#536](https://github.com/Bearer/bearer/issues/536)) ([4cda4c2](https://github.com/Bearer/bearer/commit/4cda4c2))





## [0.97.1](https://github.com/Bearer/bearer/compare/v0.97.0...v0.97.1) (2019-02-26)


### Bug Fixes

* **core:** update event name ([09c0f4c](https://github.com/Bearer/bearer/commit/09c0f4c))





# [0.97.0](https://github.com/Bearer/bearer/compare/v0.96.1...v0.97.0) (2019-02-26)


### Bug Fixes

* **core:** pass authId to revoke request ([9db5854](https://github.com/Bearer/bearer/commit/9db5854))
* **js:** remove logger use for now ([af3a772](https://github.com/Bearer/bearer/commit/af3a772))
* **ui:** forward authId to bearer revoke request ([d374415](https://github.com/Bearer/bearer/commit/d374415))


### Features

* **js:** backport functionalities ([acd9fba](https://github.com/Bearer/bearer/commit/acd9fba))
* **openapi-generator:** introduce openapi-generator package ([#537](https://github.com/Bearer/bearer/issues/537)) ([e154413](https://github.com/Bearer/bearer/commit/e154413))





## [0.96.1](https://github.com/Bearer/bearer/compare/v0.96.0...v0.96.1) (2019-02-22)


### Bug Fixes

* **logger:** fix browser distribution ([d5ecfde](https://github.com/Bearer/bearer/commit/d5ecfde))





# [0.96.0](https://github.com/Bearer/bearer/compare/v0.95.1...v0.96.0) (2019-02-21)


### Bug Fixes

* **js:** add packave version ([8a3049f](https://github.com/Bearer/bearer/commit/8a3049f))


### Features

* **intents:** allow Alice to set error response ([#526](https://github.com/Bearer/bearer/issues/526)) ([ca47652](https://github.com/Bearer/bearer/commit/ca47652))
* **js:** create base of [@bearer](https://github.com/bearer)/js package ([f440483](https://github.com/Bearer/bearer/commit/f440483))





## [0.95.1](https://github.com/Bearer/bearer/compare/v0.95.0...v0.95.1) (2019-02-20)

**Note:** Version bump only for package bearer-master





# [0.95.0](https://github.com/Bearer/bearer/compare/v0.94.1...v0.95.0) (2019-02-15)


### Bug Fixes

* **cli:** add placeholder to test translations ([0bb4e35](https://github.com/Bearer/bearer/commit/0bb4e35))
* **core:** export fixes ([1aa5e83](https://github.com/Bearer/bearer/commit/1aa5e83))
* **transpiler:** remove aliases ([c130c18](https://github.com/Bearer/bearer/commit/c130c18))
* **ui:** add missing definition ([a760c7a](https://github.com/Bearer/bearer/commit/a760c7a))
* **ui:** add react hack ([e32c619](https://github.com/Bearer/bearer/commit/e32c619))
* **ui:** fix hot reload ([1319506](https://github.com/Bearer/bearer/commit/1319506))
* **ui:** make i18n component accept a scope ([31b7170](https://github.com/Bearer/bearer/commit/31b7170))


### Features

* **core:** update i18n signatures before transpiler rewrite ([91fe00a](https://github.com/Bearer/bearer/commit/91fe00a))
* **transpiler:** add scope attribute to i18n elements ([65ad24e](https://github.com/Bearer/bearer/commit/65ad24e))
* **transpiler:** alias non aliased imports ([7936133](https://github.com/Bearer/bearer/commit/7936133))
* **transpiler:** enable i18n modifier actions ([574de9c](https://github.com/Bearer/bearer/commit/574de9c))
* **transpiler:** inject scoped variables ([b237b30](https://github.com/Bearer/bearer/commit/b237b30))
* **transpiler:** rename t/p to aliased version ([6335599](https://github.com/Bearer/bearer/commit/6335599))
* **transpiler:** use case not handled for the moment ([24c7567](https://github.com/Bearer/bearer/commit/24c7567))





## [0.94.1](https://github.com/Bearer/bearer/compare/v0.94.0...v0.94.1) (2019-02-12)


### Bug Fixes

* **cli:** fix maximum version of [@types](https://github.com/types) node ([9841329](https://github.com/Bearer/bearer/commit/9841329))





# [0.94.0](https://github.com/Bearer/bearer/compare/v0.93.4...v0.94.0) (2019-02-12)


### Features

* update intent endpoint ([#513](https://github.com/Bearer/bearer/issues/513)) ([9363f8d](https://github.com/Bearer/bearer/commit/9363f8d))





## [0.93.4](https://github.com/Bearer/bearer/compare/v0.93.3...v0.93.4) (2019-02-12)

**Note:** Version bump only for package bearer-master





## [0.93.3](https://github.com/Bearer/bearer/compare/v0.93.2...v0.93.3) (2019-02-05)


### Bug Fixes

* **cli:** remove jest for now ([4252d66](https://github.com/Bearer/bearer/commit/4252d66))





## [0.93.2](https://github.com/Bearer/bearer/compare/v0.93.1...v0.93.2) (2019-02-05)


### Bug Fixes

* **cli:** allow email to be passed as env variable ([613dc96](https://github.com/Bearer/bearer/commit/613dc96))





## [0.93.1](https://github.com/Bearer/bearer/compare/v0.93.0...v0.93.1) (2019-02-05)


### Bug Fixes

* **cli:** update hard coded beta version for now ([7ac2b14](https://github.com/Bearer/bearer/commit/7ac2b14))





# [0.93.0](https://github.com/Bearer/bearer/compare/v0.92.1...v0.93.0) (2019-02-05)


### Bug Fixes

* **bearer-cli:** add missing routes to support v3 ([#511](https://github.com/Bearer/bearer/issues/511)) ([b76c3dc](https://github.com/Bearer/bearer/commit/b76c3dc))


### Features

* bump all versions ([4fb7615](https://github.com/Bearer/bearer/commit/4fb7615))
* need to bump version for everything ([f413953](https://github.com/Bearer/bearer/commit/f413953))





## [0.92.1](https://github.com/Bearer/bearer/compare/v0.92.0...v0.92.1) (2019-02-01)


### Bug Fixes

* **core:** typo ([ed5851d](https://github.com/Bearer/bearer/commit/ed5851d))





# [0.92.0](https://github.com/Bearer/bearer/compare/v0.91.8...v0.92.0) (2019-01-31)


### Features

* **ui:** minor enhancements to button-popover ([#505](https://github.com/Bearer/bearer/issues/505)) ([d527ac4](https://github.com/Bearer/bearer/commit/d527ac4))





## [0.91.8](https://github.com/Bearer/bearer/compare/v0.91.7...v0.91.8) (2019-01-30)


### Bug Fixes

* **transpiler:** allow using a different prop as referenceId value ([b6da586](https://github.com/Bearer/bearer/commit/b6da586))





## [0.91.7](https://github.com/Bearer/bearer/compare/v0.91.6...v0.91.7) (2019-01-30)


### Bug Fixes

* **legacy-cli:** failed intents should return 200 ([#503](https://github.com/Bearer/bearer/issues/503)) ([4763fa7](https://github.com/Bearer/bearer/commit/4763fa7))





## [0.91.6](https://github.com/Bearer/bearer/compare/v0.91.5...v0.91.6) (2019-01-30)

**Note:** Version bump only for package bearer-master





## [0.91.5](https://github.com/Bearer/bearer/compare/v0.91.4...v0.91.5) (2019-01-29)


### Bug Fixes

* **transpiler:** trigger publish ([7b561fe](https://github.com/Bearer/bearer/commit/7b561fe))





## [0.91.4](https://github.com/Bearer/bearer/compare/v0.91.3...v0.91.4) (2019-01-29)


### Bug Fixes

* **ui:** Fix button shadow ([#500](https://github.com/Bearer/bearer/issues/500)) ([4bc562d](https://github.com/Bearer/bearer/commit/4bc562d))





## [0.91.3](https://github.com/Bearer/bearer/compare/v0.91.2...v0.91.3) (2019-01-29)


### Bug Fixes

* **tranpiler:** run tranpile phase twice ([678e9cc](https://github.com/Bearer/bearer/commit/678e9cc))





## [0.91.2](https://github.com/Bearer/bearer/compare/v0.91.1...v0.91.2) (2019-01-25)


### Bug Fixes

* **cli:** [CORE-191] - fix script packages ([#495](https://github.com/Bearer/bearer/issues/495)) ([183d22a](https://github.com/Bearer/bearer/commit/183d22a))





## [0.91.1](https://github.com/Bearer/bearer/compare/v0.91.0...v0.91.1) (2019-01-24)


### Bug Fixes

* **create-bearer:** prefix command with new + doc ([#494](https://github.com/Bearer/bearer/issues/494)) ([089a7ce](https://github.com/Bearer/bearer/commit/089a7ce))





# [0.91.0](https://github.com/Bearer/bearer/compare/v0.90.6...v0.91.0) (2019-01-24)


### Features

* **cli:** allow password to be passed as env variable ([#491](https://github.com/Bearer/bearer/issues/491)) ([bfdf4d3](https://github.com/Bearer/bearer/commit/bfdf4d3))
* **cli:** bearer cli setup action components templates ([#492](https://github.com/Bearer/bearer/issues/492)) ([06b5d5c](https://github.com/Bearer/bearer/commit/06b5d5c))
* **create-bearer:** allow npm/yarn init npx command runners ([#493](https://github.com/Bearer/bearer/issues/493)) ([56145fa](https://github.com/Bearer/bearer/commit/56145fa))





## [0.90.6](https://github.com/Bearer/bearer/compare/v0.90.5...v0.90.6) (2019-01-23)


### Bug Fixes

* **cli:** update intent filter to handle new intent format ([#489](https://github.com/Bearer/bearer/issues/489)) ([4dab3aa](https://github.com/Bearer/bearer/commit/4dab3aa))





## [0.90.5](https://github.com/Bearer/bearer/compare/v0.90.4...v0.90.5) (2019-01-23)


### Bug Fixes

* **intents:** remove unnecessary typing ([#490](https://github.com/Bearer/bearer/issues/490)) ([e2f5eed](https://github.com/Bearer/bearer/commit/e2f5eed))





## [0.90.4](https://github.com/Bearer/bearer/compare/v0.90.3...v0.90.4) (2019-01-23)


### Bug Fixes

* **intents:** make intent writing less restrictive ([#488](https://github.com/Bearer/bearer/issues/488)) ([75d8145](https://github.com/Bearer/bearer/commit/75d8145))





## [0.90.3](https://github.com/Bearer/bearer/compare/v0.90.2...v0.90.3) (2019-01-23)


### Bug Fixes

* update tag names when using scenario name scoping ([#483](https://github.com/Bearer/bearer/issues/483)) ([c7c48db](https://github.com/Bearer/bearer/commit/c7c48db))
* **legacy-cli:** make CLI works as intended ([#487](https://github.com/Bearer/bearer/issues/487)) ([0f6ace7](https://github.com/Bearer/bearer/commit/0f6ace7))





## [0.90.2](https://github.com/Bearer/bearer/compare/v0.90.1...v0.90.2) (2019-01-22)

**Note:** Version bump only for package bearer-master





## [0.90.1](https://github.com/Bearer/bearer/compare/v0.90.0...v0.90.1) (2019-01-22)

**Note:** Version bump only for package bearer-master





<a name="0.90.0"></a>
# [0.90.0](https://github.com/Bearer/bearer/compare/v0.89.2...v0.90.0) (2019-01-22)


### Features

* **cli:** Core 216/intents as instance ([#482](https://github.com/Bearer/bearer/issues/482)) ([6ed680f](https://github.com/Bearer/bearer/commit/6ed680f))
* [CORE-216] - intents abstract classes ([#484](https://github.com/Bearer/bearer/issues/484)) ([035684c](https://github.com/Bearer/bearer/commit/035684c))





<a name="0.89.2"></a>
## [0.89.2](https://github.com/Bearer/bearer/compare/v0.89.1...v0.89.2) (2019-01-17)


### Bug Fixes

* **cli:** change the way we generate files ([cc36d56](https://github.com/Bearer/bearer/commit/cc36d56))
* **cli:** disable open api spec generation for now ([67e8dde](https://github.com/Bearer/bearer/commit/67e8dde))
* **cli:** stop using templates repository ([4a33bf1](https://github.com/Bearer/bearer/commit/4a33bf1))





<a name="0.89.1"></a>
## [0.89.1](https://github.com/Bearer/bearer/compare/v0.89.0...v0.89.1) (2019-01-17)


### Bug Fixes

* **ui:** make bearer form great again ([#478](https://github.com/Bearer/bearer/issues/478)) ([660fa7f](https://github.com/Bearer/bearer/commit/660fa7f))





<a name="0.89.0"></a>
# [0.89.0](https://github.com/Bearer/bearer/compare/v0.88.0...v0.89.0) (2019-01-17)


### Bug Fixes

* migrate to async for good ([#477](https://github.com/Bearer/bearer/issues/477)) ([da95eb5](https://github.com/Bearer/bearer/commit/da95eb5))


### Features

* **ui:** dropdown style as variable ([#472](https://github.com/Bearer/bearer/issues/472)) ([3b48117](https://github.com/Bearer/bearer/commit/3b48117))





<a name="0.88.0"></a>
# [0.88.0](https://github.com/Bearer/bearer/compare/v0.87.3...v0.88.0) (2019-01-17)


### Bug Fixes

* **transpiler:** input output transpiler intent arguments ([90cd61a](https://github.com/Bearer/bearer/commit/90cd61a))


### Features

* **intents:** add first fetsh promise version ([cf98b8c](https://github.com/Bearer/bearer/commit/cf98b8c))
* **intents:** first version of promisified fetch intent ([d9fbbdc](https://github.com/Bearer/bearer/commit/d9fbbdc))
* **intents:** implement save state intent ([45284b9](https://github.com/Bearer/bearer/commit/45284b9))





<a name="0.87.3"></a>
## [0.87.3](https://github.com/Bearer/bearer/compare/v0.87.2...v0.87.3) (2019-01-16)


### Bug Fixes

* **legacy-cli:** allow fetching data using setup-id ([#474](https://github.com/Bearer/bearer/issues/474)) ([8f11d52](https://github.com/Bearer/bearer/commit/8f11d52))





<a name="0.87.2"></a>
## [0.87.2](https://github.com/Bearer/bearer/compare/v0.87.1...v0.87.2) (2019-01-16)


### Bug Fixes

* **core:** allow setupId to be passed as parameter of the intent call ([d2a18f0](https://github.com/Bearer/bearer/commit/d2a18f0))





<a name="0.87.1"></a>
## [0.87.1](https://github.com/Bearer/bearer/compare/v0.87.0...v0.87.1) (2019-01-15)


### Bug Fixes

* **intents:** allow not passing ref-id to SaveState intent ([769ef42](https://github.com/Bearer/bearer/commit/769ef42))





<a name="0.87.0"></a>
# [0.87.0](https://github.com/Bearer/bearer/compare/v0.86.2...v0.87.0) (2019-01-15)


### Bug Fixes

* don't set property when data is undefined ([#466](https://github.com/Bearer/bearer/issues/466)) ([558c6d5](https://github.com/Bearer/bearer/commit/558c6d5))
* open to closed ([0e236a4](https://github.com/Bearer/bearer/commit/0e236a4))
* **tsconfig:** add null checks ([9b0c945](https://github.com/Bearer/bearer/commit/9b0c945))
* **ui:** allow dropdown-button to be customized ([#469](https://github.com/Bearer/bearer/issues/469)) ([54c0b8e](https://github.com/Bearer/bearer/commit/54c0b8e))


### Features

* **types,transpiler:** allow passing arguments to input/output decorators ([68cb454](https://github.com/Bearer/bearer/commit/68cb454))
* **ui:** add action button ([c14879a](https://github.com/Bearer/bearer/commit/c14879a))
* **ui:** bearer button customization ([f77fd1e](https://github.com/Bearer/bearer/commit/f77fd1e))





<a name="0.86.2"></a>
## [0.86.2](https://github.com/Bearer/bearer/compare/v0.86.1...v0.86.2) (2019-01-07)


### Bug Fixes

* **cli:** change default Authorization header value ([07ce7c8](https://github.com/Bearer/bearer/commit/07ce7c8))





<a name="0.86.1"></a>
## [0.86.1](https://github.com/Bearer/bearer/compare/v0.86.0...v0.86.1) (2019-01-03)


### Bug Fixes

* **intents:** rewrite auth context types ([f276847](https://github.com/Bearer/bearer/commit/f276847))





<a name="0.86.0"></a>
# [0.86.0](https://github.com/Bearer/bearer/compare/v0.85.4...v0.86.0) (2019-01-03)


### Bug Fixes

* **security:** remove deprecation warning ([a30afc6](https://github.com/Bearer/bearer/commit/a30afc6))


### Features

* **cli:** introduce cli encrypt/decrypt helpers ([76dbb33](https://github.com/Bearer/bearer/commit/76dbb33))





<a name="0.85.4"></a>
## [0.85.4](https://github.com/Bearer/bearer/compare/v0.85.3...v0.85.4) (2019-01-02)


### Bug Fixes

* **react:** allow different int host for react applications ([9024547](https://github.com/Bearer/bearer/commit/9024547))
* **react:** make it discoverable on npm ([30923f5](https://github.com/Bearer/bearer/commit/30923f5))





<a name="0.85.3"></a>
## [0.85.3](https://github.com/Bearer/bearer/compare/v0.85.2...v0.85.3) (2018-12-31)


### Bug Fixes

* **node:** add better error ([ff191e1](https://github.com/Bearer/bearer/commit/ff191e1))





<a name="0.85.2"></a>
## [0.85.2](https://github.com/Bearer/bearer/compare/v0.85.1...v0.85.2) (2018-12-31)

**Note:** Version bump only for package bearer-master





<a name="0.85.1"></a>
## [0.85.1](https://github.com/Bearer/bearer/compare/v0.85.0...v0.85.1) (2018-12-28)


### Bug Fixes

* **core:** use v1 items endpoints ([3d3a161](https://github.com/Bearer/bearer/commit/3d3a161))





<a name="0.85.0"></a>
# [0.85.0](https://github.com/Bearer/bearer/compare/v0.84.2...v0.85.0) (2018-12-27)


### Features

* **security:** simplify usage for now ([100154c](https://github.com/Bearer/bearer/commit/100154c))





<a name="0.84.2"></a>
## [0.84.2](https://github.com/Bearer/bearer/compare/v0.84.1...v0.84.2) (2018-12-27)


### Bug Fixes

* **core:** delegate path to request builders ([41094f5](https://github.com/Bearer/bearer/commit/41094f5))





<a name="0.84.1"></a>
## [0.84.1](https://github.com/Bearer/bearer/compare/v0.84.0...v0.84.1) (2018-12-27)


### Bug Fixes

* **core:** filter falsy params ([21d7760](https://github.com/Bearer/bearer/commit/21d7760))





<a name="0.84.0"></a>
# [0.84.0](https://github.com/Bearer/bearer/compare/v0.83.5...v0.84.0) (2018-12-21)


### Bug Fixes

* **security:** fix typing ([35ed19e](https://github.com/Bearer/bearer/commit/35ed19e))


### Features

* **tsconfig:** create tsconfig package ([0be2ac5](https://github.com/Bearer/bearer/commit/0be2ac5))





<a name="0.83.5"></a>
## [0.83.5](https://github.com/Bearer/bearer/compare/v0.83.4...v0.83.5) (2018-12-21)


### Bug Fixes

* **tslint-config:** remove comment from json ([5896e43](https://github.com/Bearer/bearer/commit/5896e43))





<a name="0.83.4"></a>
## [0.83.4](https://github.com/Bearer/bearer/compare/v0.83.3...v0.83.4) (2018-12-21)

**Note:** Version bump only for package bearer-master





<a name="0.83.3"></a>
## [0.83.3](https://github.com/Bearer/bearer/compare/v0.83.2...v0.83.3) (2018-12-21)


### Bug Fixes

* **security:** fix package build ([#447](https://github.com/Bearer/bearer/issues/447)) ([4c250b3](https://github.com/Bearer/bearer/commit/4c250b3))





<a name="0.83.2"></a>
## [0.83.2](https://github.com/Bearer/bearer/compare/v0.83.1...v0.83.2) (2018-12-20)


### Bug Fixes

* **cli:** handle new intent routes ([068cb23](https://github.com/Bearer/bearer/commit/068cb23))





<a name="0.83.1"></a>
## [0.83.1](https://github.com/Bearer/bearer/compare/v0.83.0...v0.83.1) (2018-12-20)


### Bug Fixes

* **core:** use new intent endpoints ([b97fc28](https://github.com/Bearer/bearer/commit/b97fc28))





<a name="0.83.0"></a>
# [0.83.0](https://github.com/Bearer/bearer/compare/v0.82.0...v0.83.0) (2018-12-20)


### Bug Fixes

* **core:** forward client id to auth request ([d661714](https://github.com/Bearer/bearer/commit/d661714))


### Features

* **core:** use integration-service v2 endpoint ([43156f9](https://github.com/Bearer/bearer/commit/43156f9))





<a name="0.82.0"></a>
# [0.82.0](https://github.com/Bearer/bearer/compare/v0.81.9...v0.82.0) (2018-12-19)


### Features

* **security:** introduce security package ([#441](https://github.com/Bearer/bearer/issues/441)) ([aae485f](https://github.com/Bearer/bearer/commit/aae485f))





<a name="0.81.9"></a>
## [0.81.9](https://github.com/Bearer/bearer/compare/v0.81.8...v0.81.9) (2018-12-19)


### Bug Fixes

* tslint config as a public package ([6ad4ab4](https://github.com/Bearer/bearer/commit/6ad4ab4))





<a name="0.81.8"></a>
## [0.81.8](https://github.com/Bearer/bearer/compare/v0.81.7...v0.81.8) (2018-12-19)


### Bug Fixes

* **core:** filter falsy query parameters ([1700d84](https://github.com/Bearer/bearer/commit/1700d84))





<a name="0.81.7"></a>
## [0.81.7](https://github.com/Bearer/bearer/compare/v0.81.6...v0.81.7) (2018-12-19)


### Bug Fixes

* **react:** make the package public ([d45c6a6](https://github.com/Bearer/bearer/commit/d45c6a6))





<a name="0.81.6"></a>
## [0.81.6](https://github.com/Bearer/bearer/compare/v0.81.5...v0.81.6) (2018-12-19)


### Bug Fixes

* **react:** fix issue where context would not update + expose context ([#437](https://github.com/Bearer/bearer/issues/437)) ([811c4b4](https://github.com/Bearer/bearer/commit/811c4b4))





<a name="0.81.5"></a>
## [0.81.5](https://github.com/Bearer/bearer/compare/v0.81.4...v0.81.5) (2018-12-18)


### Bug Fixes

* **react:** fix build scripts and package naming ([#436](https://github.com/Bearer/bearer/issues/436)) ([3b3ab50](https://github.com/Bearer/bearer/commit/3b3ab50))





<a name="0.81.4"></a>
## [0.81.4](https://github.com/Bearer/bearer/compare/v0.81.3...v0.81.4) (2018-12-18)

**Note:** Version bump only for package bearer-master





<a name="0.81.3"></a>
## [0.81.3](https://github.com/Bearer/bearer/compare/v0.81.2...v0.81.3) (2018-12-18)

**Note:** Version bump only for package bearer-master





<a name="0.81.2"></a>
## [0.81.2](https://github.com/Bearer/bearer/compare/v0.81.1...v0.81.2) (2018-12-17)


### Bug Fixes

* **logger:** add prepare script ([31c7b0b](https://github.com/Bearer/bearer/commit/31c7b0b))





<a name="0.81.1"></a>
## [0.81.1](https://github.com/Bearer/bearer/compare/v0.81.0...v0.81.1) (2018-12-17)


### Bug Fixes

* display better error in a local development env ([520903d](https://github.com/Bearer/bearer/commit/520903d))
* **logger:** fix package files ([98b279c](https://github.com/Bearer/bearer/commit/98b279c))





<a name="0.81.0"></a>
# [0.81.0](https://github.com/Bearer/bearer/compare/v0.80.1...v0.81.0) (2018-12-17)


### Features

* create a logger package ([0aafd77](https://github.com/Bearer/bearer/commit/0aafd77))





<a name="0.80.1"></a>
## [0.80.1](https://github.com/Bearer/bearer/compare/v0.80.0...v0.80.1) (2018-12-17)


### Bug Fixes

* **tslint-config:** fix typo into tslint.json ([4c461c6](https://github.com/Bearer/bearer/commit/4c461c6))





<a name="0.80.0"></a>
# [0.80.0](https://github.com/Bearer/bearer/compare/v0.79.0...v0.80.0) (2018-12-17)


### Features

* **tslint-config:** Bearer tslint config ([#427](https://github.com/Bearer/bearer/issues/427)) ([32cb172](https://github.com/Bearer/bearer/commit/32cb172))





<a name="0.79.0"></a>
# [0.79.0](https://github.com/Bearer/bearer/compare/v0.78.0...v0.79.0) (2018-12-14)


### Features

* **transpiler:** add events to mutable prop changes for root components ([4912aa8](https://github.com/Bearer/bearer/commit/4912aa8))





<a name="0.78.0"></a>
# [0.78.0](https://github.com/Bearer/bearer/compare/v0.77.0...v0.78.0) (2018-12-14)


### Features

* **package-init:** create package init ([#423](https://github.com/Bearer/bearer/issues/423)) ([d57f4f6](https://github.com/Bearer/bearer/commit/d57f4f6))





<a name="0.77.0"></a>
# [0.77.0](https://github.com/Bearer/bearer/compare/v0.76.0...v0.77.0) (2018-12-13)


### Bug Fixes

* **cli:** use authv2 ([#422](https://github.com/Bearer/bearer/issues/422)) ([cbdba34](https://github.com/Bearer/bearer/commit/cbdba34))


### Features

* **transpiler:** change event name seperator to - ([#416](https://github.com/Bearer/bearer/issues/416)) ([666847e](https://github.com/Bearer/bearer/commit/666847e))





<a name="0.76.0"></a>
# [0.76.0](https://github.com/Bearer/bearer/compare/v0.75.0...v0.76.0) (2018-12-12)


### Features

* **express:** pass req down to webhook handlers ([#418](https://github.com/Bearer/bearer/issues/418)) ([9eed3d3](https://github.com/Bearer/bearer/commit/9eed3d3))





<a name="0.75.0"></a>
# [0.75.0](https://github.com/Bearer/bearer/compare/v0.74.7...v0.75.0) (2018-12-12)


### Bug Fixes

* **legacy-cli:** add missing v2 endpoints to local dev server ([#417](https://github.com/Bearer/bearer/issues/417)) ([6bf9350](https://github.com/Bearer/bearer/commit/6bf9350))


### Features

* **node:** add bearer node client  ([#410](https://github.com/Bearer/bearer/issues/410)) ([95c26de](https://github.com/Bearer/bearer/commit/95c26de))





<a name="0.74.7"></a>
## [0.74.7](https://github.com/Bearer/bearer/compare/v0.74.6...v0.74.7) (2018-12-11)

**Note:** Version bump only for package bearer-master





<a name="0.74.6"></a>
## [0.74.6](https://github.com/Bearer/bearer/compare/v0.74.5...v0.74.6) (2018-12-07)


### Bug Fixes

* **cli:** log body with a json style ([#413](https://github.com/Bearer/bearer/issues/413)) ([e78ee4b](https://github.com/Bearer/bearer/commit/e78ee4b))
* don't create new refId when it's present in request ([#414](https://github.com/Bearer/bearer/issues/414)) ([4cb09da](https://github.com/Bearer/bearer/commit/4cb09da))





<a name="0.74.5"></a>
## [0.74.5](https://github.com/Bearer/bearer/compare/v0.74.4...v0.74.5) (2018-12-07)


### Bug Fixes

* **cli:** get correct configuration back ([#412](https://github.com/Bearer/bearer/issues/412)) ([984145a](https://github.com/Bearer/bearer/commit/984145a))





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
