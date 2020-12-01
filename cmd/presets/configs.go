package presets

// auto generated file

// GetConfigs get all presets configs
func GetConfigs() map[string]string {
	var configs = make(map[string]string)
	configs["adonis"] = `language: javascript
commands:
  create:
    - kool docker kooldev/node:14-adonis adonis new
`
	configs["golang-cli"] = `language: golang
`
	configs["laravel"] = `language: php
commands:
  create:
  - kool docker kooldev/php:7.4 composer create-project --prefer-dist laravel/laravel
questions:
  database:
    message: What database service do you want to use
    options:
      - MySQL 8.0
      - MySQL 5.7
      - ProstgreSQL 13.0
      - none
  cache:
    message: What cache service do you want to use
    options:
      - Redis 6.0
      - Memcached 1.6
      - none

`
	configs["nestjs"] = `language: javascript
commands:
  create:
  - kool docker kooldev/node:14-nest nest new

`
	configs["nextjs"] = `language: javascript
commands:
  create:
  - kool docker kooldev/node:14 yarn create next-app

`
	configs["nextjs-static"] = `language: javascript
commands:
  create:
  - kool docker kooldev/node:14 yarn create next-app

`
	configs["nuxtjs"] = `language: javascript
commands:
  create:
  - kool docker kooldev/node:14 yarn create nuxt-app

`
	configs["nuxtjs-static"] = `language: javascript
commands:
  create:
  - kool docker kooldev/node:14 yarn create nuxt-app

`
	configs["symfony"] = `language: php
commands:
  create:
  - kool docker kooldev/php:7.4 composer create-project --prefer-dist symfony/website-skeleton
questions:
  database:
    message: What database service do you want to use
    options:
      - MySQL 8.0
      - MySQL 5.7
      - ProstgreSQL 13.0
      - none
  cache:
    message: What cache service do you want to use
    options:
      - Redis 6.0
      - Memcached 1.6
      - none

`
	configs["wordpress"] = `language: php
commands:
  create:
  - kool docker kooldev/php:7.4 composer create-project --prefer-dist laravel/laravel
questions:
  database:
    message: What database service do you want to use
    options:
      - MySQL 8.0
      - MySQL 5.7
      - ProstgreSQL 13.0
      - none
  cache:
    message: What cache service do you want to use
    options:
      - Redis 6.0
      - Memcached 1.6
      - none

`
	return configs
}
