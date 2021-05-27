package presets

// auto generated file

// GetConfigs get all presets configs
func GetConfigs() map[string]string {
	var configs = make(map[string]string)
	configs["adonis"] = `language: javascript
commands:
  create:
    - kool docker kooldev/node:14-adonis adonis new $CREATE_DIRECTORY
questions:
  compose:
    - key: database
      default_answer: MySQL 8.0
      message: What database service do you want to use
      options:
        - name: MySQL 8.0
          template: mysql8-adonis.yml
        - name: MySQL 5.7
          template: mysql57-adonis.yml
        - name: PostgreSQL 13.0
          template: postgresql13-adonis.yml
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
  kool:
    - key: scripts
      default_answer: npm
      message: What javascript package manager do you want to use
      options:
        - name: npm
          template: npm-adonis.yml
        - name: yarn
          template: yarn-adonis.yml
templates:
  - key: app
    template: node14-adonis.yml
`
	configs["golang-cli"] = `language: golang
`
	configs["hugo"] = `language: static
commands:
  create:
  - kool docker klakegg/hugo new site $CREATE_DIRECTORY
#questions:
#  compose:
#    - key: app
#      template: hugo.yml
#   - key: comments
#     default_answer: none
#     message: What comments service do you want to use
#     options:
#       - name: Commento
#         template: commento.yml
#       - name: none
#         template: none
#templates:
#  - key: scripts
#    template: hugo.yml
`
	configs["laravel"] = `language: php
commands:
  create:
  - kool docker kooldev/php:7.4 composer create-project --no-install --no-scripts --prefer-dist laravel/laravel $CREATE_DIRECTORY
questions:
  compose:
    - key: app
      default_answer: PHP 8.0
      message: What app service do you want to use
      options:
        - name: PHP 8.0
          template: php8.yml
        - name: PHP 7.4
          template: php74.yml
    - key: database
      default_answer: MySQL 8.0
      message: What database service do you want to use
      options:
        - name: MySQL 8.0
          template: mysql8.yml
        - name: MySQL 5.7
          template: mysql57.yml
        - name: MariaDB 10.5
          template: mariadb105.yml
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
  kool:
    - key: scripts
      default_answer: npm
      message: What javascript package manager do you want to use
      options:
        - name: npm
          template: npm-laravel.yml
        - name: yarn
          template: yarn-laravel.yml
templates:
  - key: scripts
    template: laravel.yml
`
	configs["nestjs"] = `language: javascript
commands:
  create:
  - kool docker kooldev/node:14-nest nest new $CREATE_DIRECTORY
questions:
  compose:
    - key: database
      default_answer: MySQL 8.0
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
  kool:
    - key: scripts
      default_answer: npm
      message: Which package manager did you choose during Nest setup
      options:
        - name: npm
          template: npm-nestjs.yml
        - name: yarn
          template: yarn-nestjs.yml
templates:
  - key: app
    template: node14-nestjs.yml
`
	configs["nextjs"] = `language: javascript
commands:
  create:
  - kool docker kooldev/node:14 yarn create next-app $CREATE_DIRECTORY
questions:
  kool:
    - key: scripts
      default_answer: npm
      message: What javascript package manager do you want to use
      options:
        - name: npm
          template: npm-nextjs.yml
        - name: yarn
          template: yarn-nextjs.yml
`
	configs["nodejs"] = `language: javascript
commands:
  create:
    - mkdir $CREATE_DIRECTORY
questions:
  kool:
    - key: scripts
      default_answer: npm
      message: What javascript package manager do you want to use
      options:
        - name: npm
          template: npm-nodejs.yml
        - name: yarn
          template: yarn-nodejs.yml
`
	configs["nuxtjs"] = `language: javascript
commands:
  create:
  - kool docker kooldev/node:14 yarn create nuxt-app $CREATE_DIRECTORY
questions:
  kool:
    - key: scripts
      default_answer: npm
      message: Which package manager did you choose during NuxtJS setup
      options:
        - name: npm
          template: npm-nuxtjs.yml
        - name: yarn
          template: yarn-nuxtjs.yml
`
	configs["php"] = `language: php
commands:
  create:
    - mkdir -p $CREATE_DIRECTORY/public
    - echo "<?php echo 'Hello World!';" > $CREATE_DIRECTORY/public/index.php
questions:
  compose:
    - key: app
      default_answer: PHP 7.4
      message: Which version of PHP do you want to use
      options:
        - name: PHP 7.4
          template: php74.yml
        - name: PHP 8.0
          template: php8.yml
  kool:
    - key: scripts
      default_answer: 1.x
      message: What composer version do you want to use
      options:
        - name: 1.x
          template: composer.yml
        - name: 2.x
          template: composer2.yml
templates:
  - key: scripts
    template: php.yml
`
	configs["symfony"] = `language: php
commands:
  create:
  - kool docker kooldev/php:7.4 composer create-project --prefer-dist symfony/website-skeleton $CREATE_DIRECTORY
questions:
  compose:
    - key: app
      default_answer: PHP 8.0
      message: What app service do you want to use
      options:
        - name: PHP 8.0
          template: php8.yml
        - name: PHP 7.4
          template: php74.yml
    - key: database
      default_answer: MySQL 8.0
      message: What database service do you want to use
      options:
        - name: MySQL 8.0
          template: mysql8.yml
        - name: MySQL 5.7
          template: mysql57.yml
        - name: MariaDB 10.5
          template: mariadb105.yml
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
  kool:
    - key: scripts
      default_answer: npm
      message: What javascript package manager do you want to use
      options:
        - name: npm
          template: npm.yml
        - name: yarn
          template: yarn.yml
templates:
  - key: scripts
    template: symfony.yml
`
	configs["wordpress"] = `language: php
commands:
  create:
    - mkdir $CREATE_DIRECTORY
questions:
  compose:
    - key: app
      default_answer: PHP 8.0
      message: What PHP version do you want to use
      options:
        - name: PHP 8.0
          template: wordpress80.yml
        - name: PHP 7.4
          template: wordpress74.yml
    - key: database
      default_answer: MySQL 8.0
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
  kool:
    - key: scripts
      default_answer: npm
      message: What javascript package manager do you want to use
      options:
        - name: npm
          template: npm.yml
        - name: yarn
          template: yarn.yml
templates:
  - key: scripts
    template: wordpress.yml
`
	return configs
}
