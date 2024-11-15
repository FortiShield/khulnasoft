/**
 * Copyright 2020 - KhulnaSoft Authors <admin@khulnasoft.com>
 * SPDX-License-Identifier: Apache-2.0
 */

const assert = require('assert')

const globalError = require('./global-error')

describe('src/reducers/global-error.js', function () {
  describe('globalError(state, action)', function () {
    it('returns the initial state', function () {
      const next = globalError(undefined, {})
      assert.strictEqual(next, null)
    })

    it('handles UNRECOVERABLE_ERROR', function () {
      const next = globalError(null, {
        type: 'UNRECOVERABLE_ERROR',
        payload: 'fake payload'
      })
      assert.strictEqual(next, 'fake payload')
    })

    it('handles NAVIGATE', function () {
      const next = globalError('nope', {
        type: 'NAVIGATE',
        payload: null
      })
      assert.strictEqual(next, null)
    })
  })
})
