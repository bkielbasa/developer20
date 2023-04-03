---
title: "New project: ecommerce"
publishdate: 2023-04-10
categories:
    - Golang
    - GoInPractice
tags:
    - golang
    - reactjs
---


I'm excited to share with you my new project - an open-source e-commerce platform. The frontend is built with ReactJS and the backend is written in Go.

## Project Goals
The goal of this project is to:

* Continuously improve and develop the platform.
* Provide an opportunity for less experienced programmers to gain experience working on a real project.
* Experiment with various methodologies and tools that may not be available in other work settings (event-driven architecture, DDD, event sourcing, and more).
* Build an application from scratch, including the entire ecosystem (frontend, backend, database, infrastructure, monitoring, observability).

## Good Development Practices
Additionally, I strive to follow best practices in software development, such as writing Architecture Decision Records (ADRs) and thorough documentation. I believe that these practices help ensure the long-term success of the project and facilitate collaboration among contributors.

I also prioritize clean, maintainable code and regularly incorporate automated testing and code reviews into my workflow to ensure the highest level of quality.

## Join the Project
If you share these values and are interested in contributing to this project, please don't hesitate to reach out. Let's work together to build something great!
The source code is available on Github at https://github.com/golang-app/ecommerce.

I prepared [a dashboard](https://github.com/orgs/golang-app/projects/1) where I try to give a general overview over what I'm working on right now as well as my future plans. I split it into a few views so you can easily find the area you may want to help.

How can you help?

* Star the project
* submit your code
* join or start any [discussion](https://github.com/golang-app/ecommerce/discussions)
* tell what you'd like to learn so we can suggest you a part of the app you can try to implement it.
* create an issue with a description about how to make the project more developer-friendly.
* tell us which feature you'd like to see but other e-commerce platform don't have
* share the project on Twitter or any other social media platform

## My plans

I have a few goals that I want to achieve. Firstly, I want to set up automatic deployments to a k8s cluster with everything I have right now (OTEL, frontend, backend, the database). My plan is to use the Oracle free tier and launch three k8s nodes. I've attempted to follow the instructions provided in [k8s the hard way on OCI](https://github.com/dansimone/kubernetes-the-hard-way-on-oci), but unfortunately, my attempts have failed thus far. I hope to succeed the next time. If you're interested in helping me with this, please let me know! I do have one requirement, and that is that the servers must be free. I don't want to use any paid services, such as [EKS](https://aws.amazon.com/eks/), [Container Engine for Kubernetes](https://www.oracle.com/pl/cloud/cloud-native/container-engine-kubernetes/), or [GKE](https://cloud.google.com/kubernetes-engine/) as long as I'm not earning any money from it.

Secondly, I want to prepare a minimal e-commerce platform that a customer can use. This will involve adding the ability to place orders, generate invoices, and make simple (mocked) payments. Once this is done, I plan to add more sophisticated features, such as product galleries, discounts, and more.

Lastly, I plan to use this platform in my own shop, which will be the first production use of it.

{{< info title="Updates will be posted on the blog" msg="I will share any significant decisions that I have made which do not fit into an ADR but are worth mentioning. I have noticed a lack of projects that document their decision-making process in a detailed manner, so my aim is to provide a verbose account of how and why certain decisions were made." >}}

If you're interested in helping me achieve these goals or contributing to the project in any other way, please reach out to me. I welcome all contributions and look forward to working together to build an amazing e-commerce platform.
