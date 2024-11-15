/**
 * Copyright 2020 - KhulnaSoft Authors <admin@khulnasoft.com>
 * SPDX-License-Identifier: Apache-2.0
 */

/** @jsx h */
const { h } = require('preact')
const { connect } = require('react-redux')

const HighlightBox = require('./../_shared/highlight-box')
const Dots = require('./../_shared/dots')
const authentication = require('./../../../action-creators/authentication')

const withAuth = (redirectTo) => (OriginalComponent) => {
  const WrappedComponent = (props) => {
    if (!props.authenticatedUser) {
      props.login(
        null, null,
        __('Please log in using your credentials.')
      )
      return (
        <HighlightBox error={props.error} flash={props.flash}>
          {__('Checking authentication')}
          <Dots />
        </HighlightBox>
      )
    }
    return <OriginalComponent {...props} />
  }

  const mapStateToProps = (state) => ({
    authenticatedUser: state.authenticatedUser,
    error: state.globalError,
    flash: state.flash
  })

  const mapDispatchToProps = {
    login: authentication.login
  }

  return connect(mapStateToProps, mapDispatchToProps)(WrappedComponent)
}

module.exports = withAuth
