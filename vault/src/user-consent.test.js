/**
 * Copyright 2020 - KhulnaSoft Authors <admin@khulnasoft.com>
 * SPDX-License-Identifier: Apache-2.0
 */

var assert = require('assert')

var consentStatus = require('./user-consent')

describe('src/user-consent.js', function () {
  describe('get()', function () {
    context('with no cookies', function () {
      it('returns null', function () {
        var result = consentStatus.get()
        assert.strictEqual(result, null)
      })
    })
  })
})
