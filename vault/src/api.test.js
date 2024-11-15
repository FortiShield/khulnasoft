/**
 * Copyright 2020 - KhulnaSoft Authors <admin@khulnasoft.com>
 * SPDX-License-Identifier: Apache-2.0
 */

var assert = require('assert')
var fetchMock = require('fetch-mock')

var api = require('./api')

describe('src/api.js', function () {
  describe('getAccount', function () {
    before(function () {
      fetchMock.get('https://server.khulnasoft.com/accounts/foo-bar', {
        status: 200,
        body: { accountId: 'foo-bar', data: 'ok' }
      })
    })

    after(function () {
      fetchMock.restore()
    })

    it('calls the given endpoint with the correct parameters', function () {
      var get = api.getAccountWith('https://server.khulnasoft.com/accounts')
      return get('foo-bar')
        .then(function (result) {
          assert.deepStrictEqual(result, { accountId: 'foo-bar', data: 'ok' })
        })
    })
  })

  describe('getEvents', function () {
    before(function () {
      fetchMock.get('https://server.khulnasoft.com/events', {
        status: 200,
        body: { events: ['a', 'b', 'c'] }
      })
    })

    after(function () {
      fetchMock.restore()
    })

    it('calls the given endpoint with the correct parameters', function () {
      var get = api.getEventsWith('https://server.khulnasoft.com/events')
      return get()
        .then(function (result) {
          assert.deepStrictEqual(result, { events: ['a', 'b', 'c'] })
        })
    })
  })
})
