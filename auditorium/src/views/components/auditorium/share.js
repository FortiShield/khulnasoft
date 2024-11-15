/**
 * Copyright 2020 - KhulnaSoft Authors <admin@khulnasoft.com>
 * SPDX-License-Identifier: Apache-2.0
 */

/** @jsx h */
const { h } = require('preact')

const Share = require('./../_shared/share')

const ShareAccount = (props) => {
  return (
    <Share
      {...props}
      headline={__('Share this account')}
      subline={__(
        'Share your KhulnaSoft Fair Web Analytics account <em class="%s">%s</em> via email invitation. If granted admin privileges, invited users have full access to a shared account and can also create further accounts for this instance.', 'i tracked',
        props.accountName
      )}
      collapsible
    />
  )
}

module.exports = ShareAccount
