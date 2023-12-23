---
title: "How to gather GC metrics in NodeJS"
publishdate: 2023-12-22
categories: 
    - Programming
tags:
  - nodejs
  - gc
  - observability
---

We can trace NodeJS GC by using
```shell
  node --trace-gc app.js
```
  
And use the [performance tool](https://nodejs.org/api/perf_hooks.html) to get the data in runtime.
  
```node
  const { PerformanceObserver } = require('perf_hooks');
  
  // Create a performance observer
  const obs = new PerformanceObserver((list) => {
    const entry = list.getEntries()[0]
    /* 
    The entry would be an instance of PerformanceEntry containing
    metrics of garbage collection.
    For example:
    PerformanceEntry {
      name: 'gc',
      entryType: 'gc',
      startTime: 2820.567669,
      duration: 1.315709,
      kind: 1
    }
    */
  });
  
  // Subscribe notifications of GCs
  obs.observe({ entryTypes: ['gc'] });
  
  // Stop subscription
  obs.disconnect();
```

To gather memory usage we can use the following code

```node
  console.log(process.memoryUsage());
  // Prints:
  // {
  //  rss: 4935680,
  //  heapTotal: 1826816,
  //  heapUsed: 650472,
  //  external: 49879,
  //  arrayBuffers: 9386
  // }
```
