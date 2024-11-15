/**
 * Copyright 2020 - KhulnaSoft Authors <admin@khulnasoft.com>
 * SPDX-License-Identifier: Apache-2.0
 */

/** @jsx h */
const { h } = require('preact')

const Share = require('./../_shared/share')

const ShareAccounts = (props) => {
  return (
    <Share
      {...props}
      headline={__('Share all accounts')}
      subline={__('Share all your KhulnaSoft Fair Web Analytics accounts via email invitation. When granted Admin privilges, invited users have full access to shared accounts and can create new accounts on this instance.')}
    />
  )
}

module.exports = ShareAccounts
