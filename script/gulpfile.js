/**
 * Copyright 2020 - KhulnaSoft Authors <admin@khulnasoft.com>
 * SPDX-License-Identifier: Apache-2.0
 */

var fs = require('fs')
var gulp = require('gulp')
var clean = require('gulp-clean')
var gap = require('gulp-append-prepend')
var buffer = require('vinyl-buffer')
var source = require('vinyl-source-stream')
var browserify = require('browserify')
var linguasFile = require('linguas-file')

var pkg = require('./package.json')

var defaultLocale = 'en'
var linguas = !process.env.SKIP_LOCALES
  ? linguasFile.parse(fs.readFileSync('./locales/LINGUAS', 'utf-8'))
  : []

gulp.task('clean:pre', function () {
  return gulp
    .src('./dist', { read: false, allowEmpty: true })
    .pipe(clean())
})

gulp.task('clean:post', function () {
  return gulp
    .src('./dist/**/*.json', { read: false, allowEmpty: true })
    .pipe(clean())
})

gulp.task('default', gulp.series(
  'clean:pre',
  gulp.series([defaultLocale].concat(linguas).map(function (locale) {
    return createLocalizedBundle(locale)
  })),
  'clean:post'
))

function createLocalizedBundle (locale) {
  var dest = './dist/' + locale + '/'
  var scriptTask = makeScriptTask(dest, locale)
  scriptTask.displayName = 'script:' + locale

  return scriptTask
}

function makeScriptTask (dest, locale) {
  return function () {
    var transforms = JSON.parse(JSON.stringify(pkg.browserify.transform))
    // we are setting this at process level so that it propagates to
    // dependencies that also require setting it
    process.env.LOCALE = locale
    var b = browserify({
      entries: './index.js',
      // See: https://github.com/nikku/karma-browserify/issues/130#issuecomment-120036815
      postFilter: function (id, file, currentPkg) {
        if (currentPkg.name === pkg.name) {
          currentPkg.browserify.transform = []
        }
        return true
      },
      transform: transforms.map(function (transform) {
        if (transform === '@khulnasoft/l10nify' || (Array.isArray(transform) && transform[0] === '@khulnasoft/l10nify')) {
          return ['@khulnasoft/l10nify']
        }
        if (transform === 'envify' || (Array.isArray(transform) && transform[0] === 'envify')) {
          return ['envify', { LOCALE: locale }]
        }
        return transform
      })
    })

    return b
      .plugin('tinyify')
      .bundle()
      .pipe(source('script.js'))
      .pipe(buffer())
      .pipe(gap.prependText('*/'))
      .pipe(gap.prependFile('./../banner.txt'))
      .pipe(gap.prependText('/**'))
      .pipe(gulp.dest(dest))
  }
}
