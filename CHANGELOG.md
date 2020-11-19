# Fischer Data Import Service Changelog

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).

Please note, that this project, while following numbering syntax, it DOES NOT
adhere to [Semantic Versioning](http://semver.org/spec/v2.0.0.html) rules.

## Types of changes

* ```Added``` for new features.
* ```Changed``` for changes in existing functionality.
* ```Deprecated``` for soon-to-be removed features.
* ```Removed``` for now removed features.
* ```Fixed``` for any bug fixes.
* ```Security``` in case of vulnerabilities.


## [2020.4.2.19] - 2020-11-19

### Changed
- added updating product with ID

## [2020.4.2.14] - 2020-11-14

### Changed
- latest libraries
- latest go 1.15.5
- updated dockerfile
- updated create.sh

## [2020.4.1.29] - 2020-10-29

###  Changed
- enabled updating and creating data in zapsi database

## [2020.4.1.26] - 2020-10-26

### Fixed
- fixed leaking goroutine bug when opening sql connections, the right way is this way

## [2020.3.3.19] - 2020-9-19

### Changed
- maps are created with their initial size

## [2020.3.3.16] - 2020-9-16

### Added
- initial commit, fully working, not tested