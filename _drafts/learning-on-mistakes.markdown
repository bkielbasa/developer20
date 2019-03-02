---
layout: post
title:  Learning on mistakes
---

There are many situations which can cause a unavailability. One of them can be a bug in a software, bad architecture design decisions or even [a human error](https://aws.amazon.com/message/41926/). Depending on how the numbers are calculated, from [22%](https://cloudscene.com/news/2017/07/datacenterdowntime/) to even [70% outages](https://www.cw.com.hk/it-hk/uptime-institute-70-dc-outages-due-to-human-error) are caused by human error. Software engineers, devops or administrators cannot prevent all the outages but we can learn from them to improve the stability of systems we are creating. In this article, I will present how Brainly learns on his mistakes and works on stability and latency of its infrastructure.

## RabbitMQ outages

RabbitMQ is the heart in our platform. Most of our microservices produce or consumes events using this event bus.
Within a week, we had 3 outages of our RabbitMQ cluster from 3 different reasons. We learned that even though RabbitMQ supports [high availability](https://www.rabbitmq.com/ha.html) some unexpected situations still may happen and cause its downtime.

During maintenance, we updated RabbitMQ configuration we added an env parameter to our infrastructure to configure a value of `ha-params` parameter dynamically. Unfortunately, by mistake, the value was put into double quotes instead of without any quotes. It led to invalid configuration schema. A person who reviewed the change did not notice the mistake and the configuration was applied on production RabbitMQ cluster.
The mistake was noticed very quickly and the patch was deployed as soon as it was possible but it led to the outage of the whole cluster for about 20 minutes.
We had no automatic checking the configuration on a development environment and this is what we fixed.

One week later, we wanted to scale up a number of RabbitMQ nodes to be able to handle bigger traffic. We wanted to increase the number of nodes from 3 to 9 and remove old nodes. It wasn’t possible to simply increase the number of nodes because the “old” nodes were running Ubuntu 14 which is quite old now. The update required their reinstallation.

Every new node was added selectively to the cluster. After forming properly running cluster, we started removing the old clusters one by one. The cluster was running without any errors, our monitoring was not showing any abnormal behavior of the RabbitMQ cluster so we assumed that cluster was running correctly without any downtime.

After a while, one of our developers noticed that his microservice stopped publishing events. A quick investigation showed that queues disappeared. All off the microservices was restarted and after the restart, everything went back to normal. But what did actually happen?

We’re developing a library written in Go which helps us deal with events - Fred. Unfortunately, Fred did not react on the RabbitMQ’s structure change and tried to use non-existing nodes. Fred should notice the change and update the available node’s list. But it did not. We fixed the bug quickly to prevent from similar situations in the future. Unfortunately, it wasn’t our last RabbitMQ outage.

At big scale, you have to deal with more than one server because only one machine cannot handle huge scale and vertical scaling becomes unprofitable. The more machines you have at your disposal, the greater the chance that something will go wrong. A failure is inevitable and it’s only a matter of time when it will happen. Our responsibility is to make sure that we are ready for such situations.

What can go wrong? Almost everything. The network connection between nodes can be slowed down or damaged. Or even someone can stumble and unplug power from the server. In our case, it was a disk failure on one of the RabbitMQ nodes. When the damaged node was removed from the cluster, another node experienced the same hard disk failure. To keep [the quorum](https://en.wikipedia.org/wiki/Quorum), we started node recovery.

When the node was back in the cluster it was unable to process any request. After a few attempts to bring back the broken node, the whole cluster of RabbitMQ crashed and each node started working independently. To prevent losing messages, our DevOps decided to restart the whole cluster. This operation solved the problem and everything went back to normal.

## Why is the story important to us?

To be able to scale fast, we use RabbitMQ as our event bus. As I mentioned before, this event bus stores all events in our infrastructure. Unfortunately, the outage of this service was not handled correctly in all microservices. It led to never-ending rebooting or losing the data in events. This situation has pushed us to action.

## What did we do to prevent similar problems in the future?

RabbitMQ is only an example of a dependency which can become unavailable at any time. It can be a database server, Varnish installation or even the whole Kafka cluster. Because we treat our job very seriously, we decided to find out the best possible solution to prevent similar situations in the future. Fixing bugs was a short-term solution but we certainly had some issues in the architecture. We write an outage note every time when something goes wrong to be able to come back to our experiences later. Our idea to fix that issue in long term was system design.

### System design
System Design is a brainstorming type of session which helped us to find solutions for problems we faced, that is: cascading failures and improving operations in case of the network partition, resulting in weekly 99.99 uptime. We found the root cause of problems and tried to find the best solution to fix it. The root

To discuss ideas we had a set of brainstorming meetings very similar to [Architectural Katas](http://nealford.com/katas/). Then, we tried to find the best and easy way to apply solutions in short and long-term.

### Ideas to find the solution

The first idea was [green-yellow-red](https://www.leadershipthoughts.com/rag-status-definition/) - approach inspired by ElasticSearch.

You can read more about system design [from many resources](https://hackernoon.com/top-10-system-design-interview-questions-for-software-engineers-8561290f0444).


