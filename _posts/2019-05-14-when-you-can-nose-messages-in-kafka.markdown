---
layout: post
title:  When you can lose messages in Kafka
mainPhoto: kafka.jpg
---

Kafka is speedy and fault-tolerant distributed streaming platform. However, there are some situations when messages can disappear. It can happen due to misconfiguration or misunderstanding Kafka's internals. In this article, I'll explain when the data loss can happen and how to prevent it.

## Publisher: Acknowledgment

When a message is sent to the publisher, the publisher waits for an acknowledgment (ACK) from the broker. There are three configuration options which can be used:

- `acks = all` - the broker will return the ACK only when all replicas will confirm that they saved the message.
- `acks = 1` - the ACK will be returned when the leader replica will save the message but won't wait for replicas to do the same
- `acks = 0` - the producer won't wait for the confirmation from the leader replica

As the last option is obvious (fire and forget), the second one may lead to less explicit data loss. There is a scenario when the producer will receive confirmation that the message was saved, but just after the ACK the leader replica crashes and doesn't start. Because other replicas don't have the message when a new replica leader is elected the message is gone forever.

By default, the `acks` is set to `1` in JVM and `all` in [golang](https://github.com/segmentio/kafka-go). As you can see, there are differences in various implementations, so it's a better idea to set this value explicitly. 

## Publisher: The buffer

For performance reasons (to reduce the network usage) a buffering can be enabled. Messages aren't sent when the publishing method is called but when the buffer reached its maximum capacity or in a given interval. Those behaviors are controlled by `batch.size` (in bytes) and `linger.ms` parameters. If one of those limitations is reached, the messages are sent at once. What's important to stress, the client will receive information that the message is already sent, but that's not true. If the app crashes before flushing the buffer, the data are irreversibly lost.

Please remember that those parameters may be different depending on the implementation. In JVM the `batch.size` is a number of bytes in the buffer (`16384` bytes by default), but in `kafka-go` the parameter describes the number of messages in the buffer (100 by default). What's more, the JVM users have the `linger.ms` set to `0` by default but `kafka-go` users have set it `1` second.

In the JVM implementation, when batching is disabled (`linger.ms=0`) messages can still be sent together. It happens under heavy load - the messages that arrive close together in time will be batched anyway.

## Subscriber: Offsets

During consuming messages, the consumer (subscriber) sends his current offset to the broker. This is the place where data loss can happen. There are, at least, two plots when it may happen.

The first scenario is consuming in parallel. Imagine a situation where 2 messages come to a consumer: A and B. All the messages are processed in parallel. Processing the messages, B was successful, and the offset was committed. However, handling the message, A produced an error. Because the message B has a larger offset, Kafka will save the latest offset and the message A never comes back to the consumer.

![](/assets/posts/kafka-commit-message.png)

## Broker: Committed doesn't mean saved on the disk

Kafka, on Linux system, saves messages to a [filesystem cache](https://www.thomas-krenn.com/en/wiki/Linux_Page_Cache_Basics) but doesn't wait the message get persisted on the hard drive. It means that if you have only one replica or `acks = 1` it is possible that the broker will go down and the message will be lost even if the broker returned the ACK.

## Broker: Saved on the hard drive doesn't mean it won't disappear

Not all the data that exists on the leader of the partition is available for clients to read. It happens when not all in-sync replicas received the message. For example, when a follower broker is behind the leader but it is still considered as in-sync  (the lag time is configured by `replica.lag.time.max.ms ` parameter, 500 by default) and then the leader crashes. A new leader is elected, but it didn't receive the message. **The message is gone**. 

This situation is the reason why consumers are not allowed to receive unsafe data.

## Summary

Kafka is an excellent tool with high capacity, consistency, and latency. On the other hand, a lot of these things depends on the configuration in the producer, consumer and even on the broker(s). What's even more problematic, the behavior can change depending on the implementation of the producer or the consumer. That's why it is essential to read the docs.

Do you know about any other scenarios where the data loss can occur? Let me know in the comments section below.

Bibliography:

- [The Definitive Guide: Real-Time Data and Stream](https://www.amazon.com/Kafka-Definitive-Real-Time-Stream-Processing-ebook/dp/B0758ZYVVN/ref=sr_1_2?keywords=Real-Time+Data+and+Stream&qid=1555916450&s=gateway&sr=8-2)
- [Kafka documentation](https://kafka.apache.org/documentation/)
