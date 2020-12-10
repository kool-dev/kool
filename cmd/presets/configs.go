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
  - key: app
    default_answer: PHP 7.4
    message: What app service do you want to use
    options:
      - name: PHP 7.4
        template: php74.yml
      - name: PHP 8.0
        template: php8.yml
  - key: database
    default_answer: MySQL 5.7
    message: What database service do you want to use
    options:
      - name: MySQL 8.0
        template: mysql8.yml
      - name: MySQL 5.7
        template: mysql57.yml
      - name: PostgreSQL 13.0
        template: postgresql13.yml
      - name: none
        template: none
  - key: cache
    default_answer: Redis 6.0
    message: What cache service do you want to use
    options:
      - name: Redis 6.0
        template: redis6.yml
      - name: Memcached 1.6
        template: memcached16.yml
      - name: none
        template: none
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
  - key: app
    default_answer: PHP 7.4
    message: What app service do you want to use
    options:
      - name: PHP 7.4
        template: php74.yml
      - name: PHP 8.0
        template: php8.yml
  - key: database
    default_answer: MySQL 5.7
    message: What database service do you want to use
    options:
      - name: MySQL 8.0
        template: mysql8.yml
      - name: MySQL 5.7
        template: mysql57.yml
      - name: PostgreSQL 13.0
        template: postgresql13.yml
      - name: none
        template: none
  - key: cache
    default_answer: Redis 6.0
    message: What cache service do you want to use
    options:
      - name: Redis 6.0
        template: redis6.yml
      - name: Memcached 1.6
        template: memcached16.yml
      - name: none
        template: none

`
	configs["wordpress"] = `language: php
commands:
  create:
  - kool docker kooldev/php:7.4 composer create-project --prefer-dist laravel/laravel
questions:
  - key: database
    default_answer: MySQL 5.7
    message: What database service do you want to use
    options:
      - name: MySQL 8.0
        template: mysql8.yml
      - name: MySQL 5.7
        template: mysql57.yml
      - name: PostgreSQL 13.0
        template: postgresql13.yml
      - name: none
        template: none
  - key: cache
    default_answer: Redis 6.0
    message: What cache service do you want to use
    options:
      - name: Redis 6.0
        template: redis6.yml
      - name: Memcached 1.6
        template: memcached16.yml
      - name: none
        template: none


`
	return configs
}
