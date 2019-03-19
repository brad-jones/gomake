## [2.3.1](https://github.com/brad-jones/gomake/compare/v2.3.0...v2.3.1) (2019-03-19)


### Bug Fixes

* **proxy-args:** in some cases arguments were not being parsed correctly ([18e9c7d](https://github.com/brad-jones/gomake/commit/18e9c7d))

# [2.3.0](https://github.com/brad-jones/gomake/compare/v2.2.2...v2.3.0) (2019-03-10)


### Features

* **executor:** working dir now set to one level above the ".gomake" dir ([55c33e5](https://github.com/brad-jones/gomake/commit/55c33e5))

## [2.2.2](https://github.com/brad-jones/gomake/compare/v2.2.1...v2.2.2) (2019-03-07)


### Bug Fixes

* **missing-parent-cmds:** now get filled in automatically by a noop func ([a7421dd](https://github.com/brad-jones/gomake/commit/a7421dd)), closes [#2](https://github.com/brad-jones/gomake/issues/2)
* **proxy-args:** "addCmdWithArgs declared and not used" error ([cb743ef](https://github.com/brad-jones/gomake/commit/cb743ef))

## [2.2.1](https://github.com/brad-jones/gomake/compare/v2.2.0...v2.2.1) (2019-03-06)


### Bug Fixes

* **go-node-versions:** bump versions of go and node used in pipeline ([7119070](https://github.com/brad-jones/gomake/commit/7119070))

# [2.2.0](https://github.com/brad-jones/gomake/compare/v2.1.3...v2.2.0) (2019-03-06)


### Features

* **proxy-args:** positional arguments are now not parsed by cobra ([665ae33](https://github.com/brad-jones/gomake/commit/665ae33))

## [2.1.3](https://github.com/brad-jones/gomake/compare/v2.1.2...v2.1.3) (2019-02-25)


### Bug Fixes

* **cmd-names:** now can have numbers in them ([4dc15cd](https://github.com/brad-jones/gomake/commit/4dc15cd))
* **error-handling:** improvements to allow for easier debugging ([d990cfa](https://github.com/brad-jones/gomake/commit/d990cfa))

## [2.1.2](https://github.com/brad-jones/gomake/compare/v2.1.1...v2.1.2) (2019-02-14)


### Bug Fixes

* **brew:** must download a tarball and extract, can not download direct ([6526e07](https://github.com/brad-jones/gomake/commit/6526e07))

## [2.1.1](https://github.com/brad-jones/gomake/compare/v2.1.0...v2.1.1) (2019-02-13)


### Bug Fixes

* **golang-modules:** use gopkg.in service now that we are to v2 ([fa348a9](https://github.com/brad-jones/gomake/commit/fa348a9))

# [2.1.0](https://github.com/brad-jones/gomake/compare/v2.0.0...v2.1.0) (2019-02-13)


### Features

* **default-env-values:** options can now consume env vars ([c3b424d](https://github.com/brad-jones/gomake/commit/c3b424d)), closes [#5](https://github.com/brad-jones/gomake/issues/5)

# [2.0.0](https://github.com/brad-jones/gomake/compare/v1.6.3...v2.0.0) (2019-02-13)


### Bug Fixes

* **option-names:** are now chain cased as one would expect for a cli app ([d8137f7](https://github.com/brad-jones/gomake/commit/d8137f7)), closes [#6](https://github.com/brad-jones/gomake/issues/6)


### BREAKING CHANGES

* **option-names:** we no longer support options with camelCase, all options are chain-case. To be clear the function parameters are still camelCase but the option comments and cobra option names are to be chain-case.

## [1.6.3](https://github.com/brad-jones/gomake/compare/v1.6.2...v1.6.3) (2019-02-13)


### Bug Fixes

* **long-description:** comments can be spaced out for readability ([e60d57e](https://github.com/brad-jones/gomake/commit/e60d57e)), closes [#4](https://github.com/brad-jones/gomake/issues/4)
* **long-description:** if only whitespace truncate to zero length string ([873973d](https://github.com/brad-jones/gomake/commit/873973d))
* **option-comments:** can now have colons in their descriptions ([2c63203](https://github.com/brad-jones/gomake/commit/2c63203)), closes [#1](https://github.com/brad-jones/gomake/issues/1)

## [1.6.2](https://github.com/brad-jones/gomake/compare/v1.6.1...v1.6.2) (2019-02-13)


### Bug Fixes

* **cache:** any go file inside a .gomake folder is added to the hash ([ae70702](https://github.com/brad-jones/gomake/commit/ae70702)), closes [#7](https://github.com/brad-jones/gomake/issues/7)

## [1.6.1](https://github.com/brad-jones/gomake/compare/v1.6.0...v1.6.1) (2019-02-13)


### Bug Fixes

* **help:** exit code now returns as 1 for unknown commands or flags ([e1edc9f](https://github.com/brad-jones/gomake/commit/e1edc9f)), closes [#8](https://github.com/brad-jones/gomake/issues/8)

# [1.6.0](https://github.com/brad-jones/gomake/compare/v1.5.0...v1.6.0) (2019-02-12)


### Features

* **rpm-deb:** added rpm and deb  package generation using nfpm ([1d33261](https://github.com/brad-jones/gomake/commit/1d33261))

# [1.5.0](https://github.com/brad-jones/gomake/compare/v1.4.1...v1.5.0) (2019-02-12)


### Features

* **scoop:** added scoop support ([764ab79](https://github.com/brad-jones/gomake/commit/764ab79))

## [1.4.1](https://github.com/brad-jones/gomake/compare/v1.4.0...v1.4.1) (2019-02-12)


### Bug Fixes

* **brew:** download link needed the 'v' prefixed to the version number ([08c4887](https://github.com/brad-jones/gomake/commit/08c4887))

# [1.4.0](https://github.com/brad-jones/gomake/compare/v1.3.1...v1.4.0) (2019-02-07)


### Features

* **brew:** added homebrew support ([5daf063](https://github.com/brad-jones/gomake/commit/5daf063))

## [1.3.1](https://github.com/brad-jones/gomake/compare/v1.3.0...v1.3.1) (2019-02-05)


### Bug Fixes

* **windows-executor:** exit from windows executor workflow correctly ([0e1a73c](https://github.com/brad-jones/gomake/commit/0e1a73c))
* **windows-runner-path:** instead of runner/.exe we now have runner.exe ([ea09f94](https://github.com/brad-jones/gomake/commit/ea09f94))

# [1.3.0](https://github.com/brad-jones/gomake/compare/v1.2.0...v1.3.0) (2019-01-28)


### Features

* **closures:** some functions in the runtime lib now return closures ([ddbd1af](https://github.com/brad-jones/gomake/commit/ddbd1af))

# [1.2.0](https://github.com/brad-jones/gomake/compare/v1.1.0...v1.2.0) (2018-12-20)


### Features

* **runtime:** added some helpers that make writing tasks easier ([5fc98b5](https://github.com/brad-jones/gomake/commit/5fc98b5))

# [1.1.0](https://github.com/brad-jones/gomake/compare/v1.0.4...v1.1.0) (2018-12-20)


### Bug Fixes

* **quotes:** now get escaped correctly when used in function doc blocks ([793b7cd](https://github.com/brad-jones/gomake/commit/793b7cd))


### Features

* **use-short-version:** these variables can now be set by a makefile ([3d913ae](https://github.com/brad-jones/gomake/commit/3d913ae))

## [1.0.4](https://github.com/brad-jones/gomake/compare/v1.0.3...v1.0.4) (2018-12-20)


### Bug Fixes

* **azure-pipeline:** moved ci over to azure dev ops for windows and mac ([4a0281e](https://github.com/brad-jones/gomake/commit/4a0281e))

## [1.0.3](https://github.com/brad-jones/gomake/compare/v1.0.2...v1.0.3) (2018-11-29)


### Bug Fixes

* **findgomakefolder:** add windows root file system detection ([7fdda94](https://github.com/brad-jones/gomake/commit/7fdda94))

## [1.0.2](https://github.com/brad-jones/gomake/compare/v1.0.1...v1.0.2) (2018-11-21)


### Bug Fixes

* **gmv:** display titles for the version information ([84017ab](https://github.com/brad-jones/gomake/commit/84017ab))

## [1.0.1](https://github.com/brad-jones/gomake/compare/v1.0.0...v1.0.1) (2018-11-21)


### Bug Fixes

* **gmv:** version injection should now work correctly ([a974287](https://github.com/brad-jones/gomake/commit/a974287))

# 1.0.0 (2018-11-21)


### Features

* **initial:** release ([ba6791b](https://github.com/brad-jones/gomake/commit/ba6791b))
