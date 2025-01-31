/**
 * Copyright 2020 - KhulnaSoft Authors <admin@khulnasoft.com>
 * SPDX-License-Identifier: Apache-2.0
 */

module.exports = router

function router () {
  var registeredRoutes = {}

  function listen (event) {
    function respond (message) {
      if (event.ports && event.ports.length > 0) {
        event.ports[0].postMessage(message)
      }
    }
    respond.selector = event.data.host

    var stack = (registeredRoutes[event.data.type] || []).slice()

    function callNext () {
      function next (err) {
        if (err) {
          if (process.env.NODE_ENV !== 'production') {
            console.error(err)
          }
          respond({
            type: 'ERROR',
            payload: {
              error: err.message,
              stack: err.stack
            }
          })
        } else {
          callNext()
        }
      }
      var nextHandler = stack.shift() || function (event, respond, next) {
        next(new Error('Message of type "' + event.data.type + '" was not handled.'))
      }
      nextHandler(event, respond, next)
    }

    callNext()
  }

  window.addEventListener('message', listen)

  var register = function (messageType /* , ...stack */) {
    var stack = [].slice.call(arguments, 1)
    registeredRoutes[messageType] = stack
  }

  register.unbind = function () {
    window.removeEventListener('message', listen)
  }

  return register
}
