---
title: "Automatic restarting apps on Kubernetes (periodically)"
publishdate: 2021-04-08
resources:
    - name: header
    - src: featured.jpg
categories:
    - DevOps
tags:
    - kubernetes
    - devops
	- service account
	- cron job
---

[Failures and downtime](https://developer20.com/unavailability-is-fine/) are part of our day-to-day life. I had a problem with one of the services that started crashing a few times a week. We noticed that it crashes because the memory usage reaches its limits no matter how high the limit is. Debugging memory leaks is hard and time-consuming. As a temporary fix[^Nothing is more permanent than a temporary solution] we decided to restart the application once a day. That that bought us time. How did I do it in Kubernetes?

## Creating the CronJob
Deployment, in Kubernetes (k8s), is a description of the desired state of pods or replica sets. If we want to restart the application then we restart the deployment. Kubernetes will run start a new deployment and wait for all pods to be healthy. It will route the traffic to the new deployment and stop all pods in the previous deployment. K8s has support for [CronJobs](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/). We will use this functionality.

```yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: deployment-restart-${DEPLOYMENT}
  namespace: ${NAMESPACE}
spec:
  concurrencyPolicy: Forbid
  schedule: '${SCHEDULE}'
  jobTemplate:
    spec:
      backoffLimit: 1
      activeDeadlineSeconds: 600
      template:
        metadata:
          annotations:
            sidecar.istio.io/inject: "false"
        spec:
          restartPolicy: Never
          containers:
            - name: kubectl
              image: bitnami/kubectl
              command:
                - 'kubectl'
                - 'rollout'
                - 'restart'
                - 'deployment/${DEPLOYMENT}'
```

Apply changes using
```
kubectl apply -f cron.yaml
```

Replace all `${DEPLOYMENT}` with the name of the deployment you want to restart and `${NAMESPACE}` of the namespace where the deployment is running. In the `schedule` field we put how the restarting process should be executed. We use the [Cron schedule syntax](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/#cron-schedule-syntax).

Notice that we disable the [istio sidecar](https://istio.io/latest/docs/setup/additional-setup/sidecar-injection/).

## The rollout command
We use `bitnami/kubectl` image where the `kubectl` command is available. We use the `rollout restart` command that can restart one of those resources: deployments, daemonsets, and statefulsets.

The `rollout` command has a set of useful subcommands:
* history - View rollout history
* pause - Mark the provided resource as paused
* restart - Restart a resource
* resume - Resume a paused resource
* status - Show the status of the rollout
* undo - Undo a previous rollout

Using this command we can manage rollouts

```sh
kubectl rollout status deployment/my-deployment
kubectl rollout pause deployment/my-deployment

; after some time
kubectl rollout resume deployment/my-deployment
```

## The service account

Kubernetes uses service accounts to authenticate pods inside the cluster to interact with the k8s API. The default service account has limited permissions. Our cron job works but when we want to add a call, for example, `kubectl rollout status` in it, it will fail. To fix that, we have to create a new role and attach (bind) it to our service account.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: deployment-restart-${DEPLOYMENT}
  namespace: ${NAMESPACE}
rules:
  - apiGroups: ["apps", "extensions"]
    resources: ["deployments"]
    resourceNames: ["${DEPLOYMENT}"]
    verbs: ["get", "patch", "list", "watch"] # extended list of operations
```

Apply changes using
```
kubectl apply -f role.yaml
```
Our new Service Account definition is available below/

```yaml
kind: ServiceAccount
apiVersion: v1
metadata:
  name: deployment-restart
  namespace: ${NAMESPACE}
```

Apply changes using
```
kubectl apply -f serviceAccount.yaml
```

We have to do two more things: bind the role to the service account and select the new service account in the cron job. Binding looks like this:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: deployment-restart
  namespace: ${NAMESPACE}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: deployment-restart
subjects:
  - kind: ServiceAccount
    name: deployment-restart # important - the name of the service account
    namespace: ${NAMESPACE}
```

Apply changes using
```
kubectl apply -f roleBinding.yaml
```

And the last change - updating the cron job to use the new service account.

```yaml
spec:
  jobTemplate:
    spec:
      template:
        spec:
          serviceAccountName: deployment-restart # name of the service
```

Apply changes using
```
kubectl apply -f cron.yaml
```

That’s it. I hope you like it and you found it useful. The technique I showed you is hacky. I know it. However, it may help to buy time and limit the money loss introduced by a bug. If you’d like to read more about Kubernetes, let me know in the comments section below.
