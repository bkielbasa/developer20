---
title: "Golang Tips & Tricks #7 - private repository and proxy"
publishdate: 2019-10-02
categories: [Golang, Programming]
tags:
  - golang
  - modules
  - proxy

resources:
    - name: header
    - src: featured.png
---

In Go 1.13 all modules are provided using a proxy. The proxy caches dependencies what helps to make sure that the version of an external dependencies will never change. If the vendor remove the version and create a new one with the same version, only the first one will be provided.
Proxy improves the performance of downloading dependencies as well so it's useful to have such functionality in the ecosystem.

The problem starts when you use a private dependency which cannot be cached using public cache. To control this behavior and exclude some dependencies from caching, the `GOPRIVATE` env variable should be used.

{{< highlight sh>}}
export GOPRIVATE=mycorp.com/*
{{< / highlight >}}

It says that all dependencies which start with `mycorp.com/` will miss the public cache. If your company host code on GitHub, you can set this variable to:

{{< highlight sh >}}
export GOPRIVATE=github.com/yourcorp/*
{{< / highlight >}}

This will make sure that all private company repos hosted on GitHub will miss the public proxy.