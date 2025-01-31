---
layout: default
title: Tagging campaigns and sources (UTM)
nav_order: 14
description: "How to better understand the traffic to your website by tagging campaigns and sources."
permalink: /running-khulnasoft/campaigns-sources/
parent: For operators
---

<!--
Copyright 2020 - KhulnaSoft Authors <admin@khulnasoft.com>
SPDX-License-Identifier: Apache-2.0
-->

# Tagging campaigns and sources (UTM)

To better understand the traffic to your website coming from links you have shared, you can add special tags to your URLs. These are then recognized by KhulnaSoft Fair Web Analytics as distinct referrers and presented separately in the Auditorium.

The two categories **Campaigns** and **Sources** are supported. They are listed in the **Top Pages** section in addition to **Referrers** and specify Sessions and Page depth in addition to the tag name.

## Setting up tags

For example, links that are shared in e-mails or text messages often do not transmit any referrers. Therefore, the list of referrers does not show any information about them. To fix this, you can simply add tags containing `utm_campaign` or `utm_source` values to the URLs you are sharing.

```
https://www.yourpage.org/?utm_campaign=Newsletter-October
```
```
https://www.yourpage.org/blog/article/?utm_source=Messenger
```
The traffic resulting from these URLs would then displayed in the respective category as "Newsletter-October" or "Messenger".

<span class="label label-green">Note</span>

For more information on how [Referrers,][Referrers] [Campaigns][Campaigns] and [Sources][Sources] are displayed in the Auditorium, head to the ["Metrics explained"][Metrics explained] section.

[Metrics explained]: /running-khulnasoft/metrics-explained/
[Referrers]: /running-khulnasoft/metrics-explained/#referrers
[Campaigns]: /running-khulnasoft/metrics-explained/#campaigns
[Sources]: /running-khulnasoft/metrics-explained/#sources
