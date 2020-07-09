---
title: "Unavailability is fine. Prepare for it"
publishdate: 2020-07-09
resources:
    - name: header
    - src: featured.jpg
categories:
    - Programming
tags:
    - system design
---

When I started my career as a software developer and published the first production application what I did was staring at logs and look for some fatal errors. It was a monolith application. Every log saying that something's wrong had to be fixed. ASAP. This approach worked for some time. However, when the scale increased, and I started building microservices, I couldn't get rid of all of them. Network issues, database failures, and more - it happens all the time.

No matter what - your application will experience outages. You have to accept it. What you should do is to try to use those errors as your tool. To be able to do that, we have to measure them in the first place.


Very often, you have a service level agreement (SLA).  The SLA is an agreement between users of our service and the provider of it. It is metrics that measure the current uptime, response time, or other things we can observe and are important from the user's point of view.

SLA should be discussed with stakeholders, and our goal is to meet those expectations. Let's say your SLA for uptime is 99.9 %. That's a lot. On the other hand, it means your service can be unavailable for 1 minute and 26 seconds every day or almost 9 hours yearly.

When you meet your SLA, you have a bit of space for experimentation. It means you can deploy new features, and even if it fails, you can rollback, fix the issue, and deploy it again. As long as the number of HTTP 5xx errors, response time or other metrics meets his SLA, you can continue experimenting with the service.

When you're service's SLO is much better than accepted SLA, other users of your API can rely on it. Users build their services of what you offer, rather than what you say you'll supply. It means that when the latency/error rate/throughput gets worst, it can have a significant impact on other services.

A way of dealing with over-dependence on the current high SLO is taking the system offline occasionally and introduce a planned outage. This will give a clear signal to consumers of your API that they shouldn't rely on the service too much and prepare for its outages.

Some tools that can help with it. Google has [Chubby](https://static.googleusercontent.com/media/research.google.com/en//archive/chubby-osdi06.pdf), Netflix open-sourced [Chaos Monkey](https://github.com/Netflix/chaosmonkey), but that's not all. There's a list with [Chaos Engineering companies, people tools and practices](https://coggle.it/diagram/WiKceGDAwgABrmyv/t/chaos-engineeringcompanies%2C-people%2C-tools-practices/0a2d4968c94723e48e1256e67df51d0f4217027143924b23517832f53c536e62) you can examine. I advise you to experiment with it

## Summary

Writing high reliable and scalable applications is a part of our job. You put a lot of effort to build the best piece of software we can and that's good. You rely on other services as well. But, do you design those services that are ready for the unavailability of your dependencies? Are you ready for the DB or event bus outage? Some time ago, I wrote an article where I describe [some of our problems with similar scenarios](https://developer20.com/learning-on-mistakes/). What's your experience on this topic? Please let me know in the comments section below.
