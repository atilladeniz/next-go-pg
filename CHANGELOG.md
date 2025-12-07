# Changelog

## [1.3.0](https://github.com/atilladeniz/next-go-pg/compare/v1.2.0...v1.3.0) (2025-12-07)


### Features

* automate database migrations for dev and production ([#28](https://github.com/atilladeniz/next-go-pg/issues/28)) ([2fe7236](https://github.com/atilladeniz/next-go-pg/commit/2fe72369a67494d733ccf9355a7fdfd5cedb5c1d))

## [1.2.0](https://github.com/atilladeniz/next-go-pg/compare/v1.1.0...v1.2.0) (2025-12-07)


### Features

* add database migrations and Prometheus metrics ([#26](https://github.com/atilladeniz/next-go-pg/issues/26)) ([d699f2f](https://github.com/atilladeniz/next-go-pg/commit/d699f2f5d90a139430b61ff235a6683a9ac61aeb))
* add rate limiting middleware for API protection ([#24](https://github.com/atilladeniz/next-go-pg/issues/24)) ([b68ae8f](https://github.com/atilladeniz/next-go-pg/commit/b68ae8f38d29b28d6526e3586d6df38afca415ad))

## [1.1.0](https://github.com/atilladeniz/next-go-pg/compare/v1.0.5...v1.1.0) (2025-12-07)


### Features

* add testing infrastructure with Vitest and Playwright ([#23](https://github.com/atilladeniz/next-go-pg/issues/23)) ([df1b5a1](https://github.com/atilladeniz/next-go-pg/commit/df1b5a18d4ef9fd99cab1d5a7e4e0fb72db3c886))

## [1.0.5](https://github.com/atilladeniz/next-go-pg/compare/v1.0.4...v1.0.5) (2025-12-07)


### Bug Fixes

* skip Docker build check on main ([#21](https://github.com/atilladeniz/next-go-pg/issues/21)) ([6be33c5](https://github.com/atilladeniz/next-go-pg/commit/6be33c52a214e59493b4fe0c7505158ad1a02470))

## [1.0.4](https://github.com/atilladeniz/next-go-pg/compare/v1.0.3...v1.0.4) (2025-12-07)


### Bug Fixes

* add packages write permission for GHCR ([#19](https://github.com/atilladeniz/next-go-pg/issues/19)) ([cf0ec07](https://github.com/atilladeniz/next-go-pg/commit/cf0ec075aa8d6eb38fcb24d9db810201afea2a5e))

## [1.0.3](https://github.com/atilladeniz/next-go-pg/compare/v1.0.2...v1.0.3) (2025-12-07)


### Bug Fixes

* remove codecov integration ([#17](https://github.com/atilladeniz/next-go-pg/issues/17)) ([9b24d6e](https://github.com/atilladeniz/next-go-pg/commit/9b24d6e753c936c72c6d78312654cc3c4d1de4ad))

## [1.0.2](https://github.com/atilladeniz/next-go-pg/compare/v1.0.1...v1.0.2) (2025-12-07)


### Bug Fixes

* add comment to backend Dockerfile ([#14](https://github.com/atilladeniz/next-go-pg/issues/14)) ([604cb8d](https://github.com/atilladeniz/next-go-pg/commit/604cb8da5ae1a3538b22bd5331e9f153355a70c8))
* frontend Dockerfile and add Docker build check in CI ([#16](https://github.com/atilladeniz/next-go-pg/issues/16)) ([4f81ff9](https://github.com/atilladeniz/next-go-pg/commit/4f81ff96b74247856263d8afee1ff88b729d6f35))

## [1.0.1](https://github.com/atilladeniz/next-go-pg/compare/v1.0.0...v1.0.1) (2025-12-07)


### Bug Fixes

* resolve Docker multi-arch build for ARM64 ([#10](https://github.com/atilladeniz/next-go-pg/issues/10)) ([cb3677d](https://github.com/atilladeniz/next-go-pg/commit/cb3677dcf04c7ade7d3730c342d2b8d0edba2193))
* trigger CI for all PRs including release-please ([#12](https://github.com/atilladeniz/next-go-pg/issues/12)) ([5928b86](https://github.com/atilladeniz/next-go-pg/commit/5928b86b46906d21f741c11c4f0b6704146215f3))

## 1.0.0 (2025-12-07)


### Features

* add CI/CD pipeline with GitHub Actions ([#7](https://github.com/atilladeniz/next-go-pg/issues/7)) ([585f9be](https://github.com/atilladeniz/next-go-pg/commit/585f9be0f93fb91874d1761a14f9519bf122e4c5))
* add comprehensive documentation for Goca  ([#4](https://github.com/atilladeniz/next-go-pg/issues/4)) ([96af5eb](https://github.com/atilladeniz/next-go-pg/commit/96af5ebf311147f90696d36789fc213dab22b0a0))
* add comprehensive system design documentation and concepts ([bdf15df](https://github.com/atilladeniz/next-go-pg/commit/bdf15dfa369401afec9e39a4081ad87a61717f6a))
* add database migration command and update dependencies ([#2](https://github.com/atilladeniz/next-go-pg/issues/2)) ([9032160](https://github.com/atilladeniz/next-go-pg/commit/9032160063887ea518a1866359eff0e2896503c4))
* add documentation for API generation and database migration commands ([06eccb5](https://github.com/atilladeniz/next-go-pg/commit/06eccb5fab471bef3e29e1e9a4b9500ee1a9baed))
* add environment validation and security test suite ([#6](https://github.com/atilladeniz/next-go-pg/issues/6)) ([9e039bf](https://github.com/atilladeniz/next-go-pg/commit/9e039bf293697add607c3c1ae636c17703d05e82))
* add pre-commit hooks with husky ([1db54ac](https://github.com/atilladeniz/next-go-pg/commit/1db54ac37a606bcd94f05f61c547a5f617c64acc))
* add search-docs functionality and update Makefile ([b58dadc](https://github.com/atilladeniz/next-go-pg/commit/b58dadcbb7c9d28847e6e4eeecbc1416d2f420a4))
* add security scanning with gitleaks and update Makefile and README ([f4bdeea](https://github.com/atilladeniz/next-go-pg/commit/f4bdeeabaa3acb23e88fe9396126b08fb33bf37c))
* enhance API functionality and improve frontend structure ([98c2662](https://github.com/atilladeniz/next-go-pg/commit/98c26623eb4e3317b609078c8cb44a750ded6b33))
* enhance buildIndex function for incremental updates ([3ab5868](https://github.com/atilladeniz/next-go-pg/commit/3ab5868871e25b7d91634c953370ffe3302c96ba))
* enhance documentation and examples for cross-tab authentication synchronization ([970c77d](https://github.com/atilladeniz/next-go-pg/commit/970c77d96dac6f74372fd7f0433f18b6f2669d4b))
* enhance Makefile and add Kamal deployment documentation ([26d2901](https://github.com/atilladeniz/next-go-pg/commit/26d2901490457cbf5561eb0ba3567e48e23f4704))
* enhance Makefile help section with branding and project details ([f7467d1](https://github.com/atilladeniz/next-go-pg/commit/f7467d1c4cb744919a88ca8fc670b40428dee6b6))
* enhance search functionality with semantic search and caching ([2ca1687](https://github.com/atilladeniz/next-go-pg/commit/2ca1687ea685a62ba8299a6dc7442ab856917c4a))
* implement centralized logging with Grafana + Loki ([#3](https://github.com/atilladeniz/next-go-pg/issues/3)) ([e050679](https://github.com/atilladeniz/next-go-pg/commit/e0506796a38cd3ea9f1a5489c310808bee8344ca))
* implement cross-tab authentication synchronization ([a095c02](https://github.com/atilladeniz/next-go-pg/commit/a095c026fcbd7518c9cfd564cc531ab9e5e18ce1))
* implement Feature-Sliced Design (FSD) architecture for frontend ([#1](https://github.com/atilladeniz/next-go-pg/issues/1)) ([2d893fc](https://github.com/atilladeniz/next-go-pg/commit/2d893fcb9adc8895f7139079f31a4d251b297e83))
* implement HydrationBoundary pattern for improved data fetching ([15d5749](https://github.com/atilladeniz/next-go-pg/commit/15d5749a2bb6ff0de4ad1f83ce421c4da01044cc))
* implement Kamal deployment configuration and documentation ([15a5902](https://github.com/atilladeniz/next-go-pg/commit/15a5902577ddfbfd772f60f522f0d4ae727aaf0b))
* implement protected layout and API testing page ([0bbcb29](https://github.com/atilladeniz/next-go-pg/commit/0bbcb29595a1dd26edcbcf7581d0d2febc208940))
* implement user statistics updates with SSE support ([6f0b415](https://github.com/atilladeniz/next-go-pg/commit/6f0b415d069dfd5688ee8b74a8be042c9b94dfbb))
* Magic Link Authentication with Enterprise Security Features ([#5](https://github.com/atilladeniz/next-go-pg/issues/5)) ([ebd4973](https://github.com/atilladeniz/next-go-pg/commit/ebd49734cb8296b3fe976d87244e9b7ce287963e))
* update backend documentation and integrate Goca CLI usage ([37774bd](https://github.com/atilladeniz/next-go-pg/commit/37774bd32d0a0f01f4763c8c3d30ff56fbe0fb73))
