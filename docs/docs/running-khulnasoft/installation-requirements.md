---
layout: default
title: Requirements for installing
nav_order: 6
description: "Considerations for installing KhulnaSoft Fair Web Analytics in an production environment."
permalink: /running-khulnasoft/installation-requirements/
parent: For operators
---

<!--
Copyright 2020 - KhulnaSoft Authors <admin@khulnasoft.com>
SPDX-License-Identifier: Apache-2.0
-->

# Requirements for installing
{: .no_toc }

In case you want to use KhulnaSoft Fair Web Analytics for collecting usage data in a production setup or similar, there are a few requirements to consider.

---

## Table of contents
{: .no_toc }

1. TOC
{:toc}

---

## Hardware requirements

KhulnaSoft Fair Web Analytics is designed to have a very low hardware footprint, so unless you are serving a lot of traffic, a single CPU core will be enough. The amount of memory consumed depends on the number of backend users logging in to your instance. If this number is rather low (up to 4), 512MB should be enough already, but if you plan on having more users, 1GB or 2GB will be a better choice and ensure a smooth experience for everyone.

## Running the application as a service

KhulnaSoft Fair Web Analytics as an application is a web server that binds to one or multiple TCP ports and listens for incoming traffic. This means that in a production setup you will need to ensure the process is always running and restarts on failure or system restart.

Choice of tools for this task depends heavily on your host OS. Alternatively, you can use the official Docker image that wraps the binary to have a unified interface across operating systems.

## Usage of a subdomain

KhulnaSoft Fair Web Analytics uses the [Same-origin policy][sop] for protecting usage data from being accessed by third party scripts.

In an example scenario where you are using your KhulnaSoft Fair Web Analytics instance to collect usage data for the domain `www.mywebsite.org`, your instance should be bound to a _different_ subdomain of _the same_ host domain, e.g. `khulnasoft.mywebsite.org`. This ensures that
1. usage data is accessible for that domain only
1. third-party cookie restrictions do not apply

In the above scenario you would use KhulnaSoft Fair Web Analytics on your website by embedding the following script:

```html
<script src="https://khulnasoft.mywebsite.org/script.js" data-account-id="<YOUR_ACCOUNT_ID>"></script>
```

In case you would use the script on a website running on an entirely different host domain, usage data would be collected for users that have third party cookies enabled only.

[sop]: https://developer.mozilla.org/en-US/docs/Web/Security/Same-origin_policy

## Serving the application via https

KhulnaSoft Fair Web Analytics itself __requires to be served using SSL__. This enables us to guarantee data is being transmitted without the possibility of third parties intercepting any communication. In case you do not have a SSL certificate for your KhulnaSoft Fair Web Analytics subdomain, you can configure it to automatically request and periodically renew a free certificate from [LetsEncrypt][lets-encrypt] for the subdomain KhulnaSoft Fair Web Analytics is being served from. See the [configuration article][config-article] for how to use this feature.

In case you do have a SSL certificate for the domain you are planning to run KhulnaSoft Fair Web Analytics on (e.g. a wildcard certificate for your top domain), pass the certificate's location to the runtime configuration and it will use the certificate. See the [configuration article][config-article] for information on how to do so.

---

While KhulnaSoft Fair Web Analytics can take care of itself being run using SSL, the protocol used to serve the host document also matters, as it defines whether browsers [consider the execution context secure][secure-context] or not. This means that in case you serve your website using plain `http`, KhulnaSoft Fair Web Analytics will not be able to use native cryptographic methods for encrypting usage data and will fall back to userland implementations instead. This approach is heavy and slow and is __not recommended__.

__Using SSL for your own site will be beneficial regarding lots of other aspects as well.__ You can check the [LetsEncrypt website][lets-encrypt] for plenty of information on how to get free and robust SSL for any setup.

[lets-encrypt]: https://letsencrypt.org/
[secure-context]: https://developer.mozilla.org/en-US/docs/Web/Security/Secure_Contexts/features_restricted_to_secure_contexts
[config-article]: /running-khulnasoft/configuring-the-application/

## Choosing a datastore

KhulnaSoft Fair Web Analytics requires a relational datastore to be available for it to store event and account data. By default, it will use a SQLite database that is stored on the host system. This works well and is a good choice if you do not serve very high amounts of traffic.

In case you want to scale it or need more performance you can also use a MySQL or Postgres database instead. See the [configuration article][config-article] for how to configure dialect and database location in such a setup.

## Additional considerations

### Running the application behind a reverse proxy

KhulnaSoft Fair Web Analytics itself is sufficiently hardened in order to be exposed to the public internet directly. If your setup requires running it behind a reverse proxy, you can do it, but __we actively advise against doing so__ for the reason that the proxy will leave possibly unwanted traces like logs containing user IPs or similar. If the reverse proxy is not a hard requirement of your setup, leave it out.

### Transactional email

KhulnaSoft Fair Web Analytics can email you a link to reset your account's password in case you forgot it. The recommended way of doing so is configuring it with SMTP credentials (you might well be able to use your default mail credentials here).

In case you do not configure this, KhulnaSoft Fair Web Analytics falls back to a local `sendmail` installation if found, yet it is likely that these messages will never arrive at all due to system restrictions or third parties bouncing email sent using that channel. Do note that the `sendmail` fallback is not available in `khulnasoft/khulnasoft` Docker containers.
