---
layout: default
title: Generating usage data for development
nav_order: 7
description: "How to generate usage data while developing KhulnaSoft Fair Web Analytics."
permalink: /developing-khulnasoft/generating-usage-data/
parent: For developers
---

<!--
Copyright 2020-2021 - KhulnaSoft Authors <admin@khulnasoft.com>
SPDX-License-Identifier: Apache-2.0
-->

# Generating usage data for development

As all usage data is encrypted and bound to user keys it is not possible to easily provide a predefined set of seed data to be used for development.

In the default development setup, a dummy site served at <http://localhost:8081> deploys the KhulnaSoft Fair Web Analytics script using the "develop" account that is available using the development login. This means you can manually create usage data by visiting and navigating this site.
