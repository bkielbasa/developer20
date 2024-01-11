---
title: "Kubernestes cluster for free on OCI"
publishdate: 2024-01-08
categories: 
    - DevOps
tags:
  - kubernetes
  - ckad
  - oracle-cloud
  - terraform
  - ansible
---

I'm preparing for the [CKAD](https://trainingportal.linuxfoundation.org/courses/certified-kubernetes-application-developer-ckad) and I needed a Kubernetes cluster to practice on. I checked many options, but none seemed to work as I wanted. Oracle Cloud offers a [free tier](https://www.oracle.com/pl/cloud/free/) tier where we have available 4 cores, 24 GB RAM, and 100 GB disk space. That's a lot. Why not use it to create a custom Kubernetes cluster that's perfect for learning and maybe running some pet projects?

In this article, I'll share with you how to set up a cluster with 1 master node and 3 workers, install MicroK8s, and configure it into a working Kubernetes cluster. We'll use Terraform and Ansible for creating the infrastructure. Thanks to that, if you break something, you'll be able to recreate it quickly.

# Setup

We have to setup our terraform provider with proper credentials. I won't describe how to get specific values because it's [described in the video](https://www.youtube.com/watch?v=GNkv6ta7VEw).

```terraform
terraform {
  required_providers {
    oci = {
      source  = "oracle/oci"
      version = "~> 4.84.0"
    }
  }
}

provider "oci" {
  tenancy_ocid         = var.tenancy_ocid
  user_ocid            = var.user_ocid
  fingerprint          = var.key_fingerprint
  private_key_path     = var.private_key_path
  private_key_password = var.private_key_password
  region               = var.region
}

module "k8s" {
  source         = "./k8s"
  compartment_id = var.tenancy_ocid
  ssh_public_key = tls_private_key.ssh.public_key_openssh
  ssh_private_key = tls_private_key.ssh.private_key_openssh
}
```
In the root module, we generate an SSH key that we will use to connect to our nodes.


```terraform
resource "tls_private_key" "ssh" {
  algorithm = "RSA"
  rsa_bits  = 4096
}
```

# Networking

Let's begin with configuring the networking. We'll network security lists to allow:

* communication between nodes
* ssh connection outside of the cluster
* ingress and egress access

```terraform
# allow communication between nodes directly using private IPs
resource "oci_core_security_list" "intra" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_virtual_network.vcn.id
  display_name   = "intra-vcn"

  ingress_security_rules {
    stateless   = true
    protocol    = "all"
    source      = "10.240.0.0/24"
    source_type = "CIDR_BLOCK"
  }
}

# allow to connect to all nodes using SSH from any IP
resource "oci_core_security_list" "ssh" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_virtual_network.vcn.id
  display_name   = "ssh"

  ingress_security_rules {
    stateless   = false
    protocol    = "6"
    source      = "0.0.0.0/0"
    source_type = "CIDR_BLOCK"

    tcp_options {
      min = 22
      max = 22
    }
  }
}

resource "oci_core_security_list" "ingress-access" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_virtual_network.vcn.id
  display_name   = "worker"

  ingress_security_rules {
    stateless   = false
    protocol    = "6"
    source      = "0.0.0.0/0"
    source_type = "CIDR_BLOCK"

    # NodePort Services
    tcp_options {
      min = 30000
      max = 32767
    }
  }
}

# allow all nodes to communicate with Internet
resource "oci_core_security_list" "egress-access" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_virtual_network.vcn.id
  display_name   = "egress access"

  egress_security_rules {
    stateless   = false
    protocol    = "all"
    destination = "0.0.0.0/0"
  }
}
```

It's the time for our new virtual network with a single subnet.

```terraform
resource "oci_core_virtual_network" "vcn" {
  cidr_block     = "10.240.0.0/24"
  dns_label      = "vcn"
  compartment_id = var.compartment_id
  display_name   = "k8s-cluster"
}

resource "oci_core_subnet" "subnet" {
  cidr_block     = "10.240.0.0/24"
  compartment_id = var.compartment_id 
  vcn_id         = oci_core_virtual_network.vcn.id

  display_name      = "kubernetes"
  dns_label         = "subnet"
  route_table_id    = oci_core_route_table.route_table.id
  security_list_ids = [oci_core_security_list.intra.id, oci_core_security_list.ingress-access.id, oci_core_security_list.egress-access.id, oci_core_security_list.ssh.id]
}

resource "oci_core_internet_gateway" "internet_gateway" {
  compartment_id = var.compartment_id 
    vcn_id         = oci_core_virtual_network.vcn.id
    display_name   = "k8s-cluster"
}
```

We pass our security lists to the subnet so our configuration will be applied there. Now it's time for a route table.


```terraform
resource "oci_core_route_table" "route_table" {
  compartment_id = var.compartment_id 
  vcn_id         = oci_core_virtual_network.vcn.id

  display_name = "k8s-cluster"
  route_rules {
    network_entity_id = oci_core_internet_gateway.internet_gateway.id
    destination       = "0.0.0.0/0"
  }
}
```

At this point, we should have all the networking set up correctly. Now it's time to create our compute instances.

# Compute instances (k8s nodes)

To create our compute instances, we need to know the ID of an Ubuntu Image we're going to use.

```terraform
data "oci_core_images" "base-image" {
    compartment_id = var.compartment_id

    operating_system = "Canonical Ubuntu"

    filter {
      name   = "operating_system_version"
      values = ["22.04"]
    }
}
```

The next step is to get availability zones. They are also required to create a new compute instance.

```terraform
data "oci_identity_availability_domains" "availability_domains" {
    compartment_id = var.compartment_id
}
```

We'll create a new node with 1 CPU and 6 GB of memory available.

```terraform
resource "oci_core_instance" "control-plane" {
  compartment_id      = var.compartment_id
  shape               = "VM.Standard.A1.Flex"
  availability_domain = lookup(data.oci_identity_availability_domains.availability_domains.availability_domains[1], "name")
  display_name        = "control-plane"
  
  source_details {
    source_id   = data.oci_core_images.base-image.images[0].id
    source_type = "image"
  }

  create_vnic_details {
    assign_public_ip = true
    subnet_id        = oci_core_subnet.subnet.id
  }

  shape_config {
    memory_in_gbs = 6
    ocpus = 1
  }

  connection {
    type        = "ssh"
    user        = "ubuntu"
    host        = self.public_ip
    private_key = var.ssh_private_key
  }

  provisioner "remote-exec" {
    inline = [
      "sudo ufw allow from 10.240.0.0/24;sudo iptables -A INPUT -i ens3 -s 10.240.0.0/24 -j ACCEPT;sudo iptables -F;sudo iptables --flush;sudo iptables -tnat --flush",
    ]
  }

  metadata = {
    "ssh_authorized_keys" = var.ssh_public_key
  }
}
```

There are a few things that need clarification. Firstly, we use the second availability zone instead of the first one. I'm doing this because in the first one, those instances that we want to use are often not available.
Secondly, we connect to the instance, just after creating, to set up the local firewall. We allow for inside-cluster communication. Last but not least, we add our freshly created public key to the instance. Thanks to this, we'll be able to SSH into it.

Creating workers is almost identical. The only difference is that I added the `count` parameter to create 3 instances of them.


```terraform
resource "oci_core_instance" "worker" {
  count = 3
  compartment_id      = var.compartment_id
  shape               = "VM.Standard.A1.Flex"
  availability_domain = lookup(data.oci_identity_availability_domains.availability_domains.availability_domains[2], "name")
  display_name        = "worker.${count.index}"
  
  source_details {
    source_id   = data.oci_core_images.base-image.images[0].id
    source_type = "image"
  }

  create_vnic_details {
    assign_public_ip = true
    subnet_id        = oci_core_subnet.subnet.id
  }

  shape_config {
    memory_in_gbs = 6
    ocpus = 1
  }

  connection {
    type        = "ssh"
    user        = "ubuntu"
    host        = self.public_ip
    private_key = var.ssh_private_key
  }

  provisioner "remote-exec" {
    inline = [
      "sudo ufw allow from 10.240.0.0/24;sudo iptables -A INPUT -i ens3 -s 10.240.0.0/24 -j ACCEPT;sudo iptables -F;sudo iptables --flush;sudo iptables -tnat --flush",
    ]
  }

  metadata = {
    "ssh_authorized_keys" = var.ssh_public_key
  }
}

```

After running `terraform init` and `terraform apply`,` you should see the plan for creating our new Kubernetes cluster.

When everything is ready, you can use the following command to copy the private key. We'll need it to SSH into our nodes.


```sh
cat terraform.tfstate | jq -r '.resources[] | select (.type == "tls_private_key" and .name == "ssh") | .instances[0].attributes.private_key_openssh' > ~/.ssh/k8s.key
```

You can test if it works by using the command:

```sh
ssh ubuntu@PUBLIC_IP -i ~/.ssh/k8s.key
```

# Setting up k8s

To automate configuring the k8s cluster, we'll use Ansible.

Let's create requirements file for our playbook.

```yaml
---
roles:
- src: https://github.com/istvano/ansible_role_microk8s
```

In this role, we will help you install MicroK8s and connect everything into a single cluster.

We need our inventory where you have to put the public and private IPs of your nodes.

```yaml
all:
  children:
    microk8s_HA:
      hosts:
        control-plane:
          ansible_host: X.X.X.X
          ansible_user: ubuntu
          ansible_ssh_private_key_file: ~/.ssh/k8s.key

    microk8s_WORKERS:
      hosts:
        worker-0:
          ansible_host: X.X.X.X
          private_ip: Y.Y.Y.Y
          ansible_user: ubuntu
          ansible_ssh_private_key_file: ~/.ssh/k8s.key
        worker-1:
          ansible_host: X.X.X.X
          private_ip: Y.Y.Y.Y
          ansible_user: ubuntu
          ansible_ssh_private_key_file: ~/.ssh/k8s.key

        worker-2:
          ansible_host: X.X.X.X
          private_ip: Y.Y.Y.Y
          ansible_user: ubuntu
          ansible_ssh_private_key_file: ~/.ssh/k8s.key
```

The private IP is needed to set up the `/etc/hosts` file on every node so they use correct IPs internally. Without it, we'll experience challenges during joining the cluster.

The playbook is divided into a few sections. The first one sets up the `/etc/hosts` file on all nodes as I described earlier. The second section installs MicroK8s for us and sets up the whole cluster. The last two parts add our ubuntu user to a proper group so we can use the `microk8s` command without a permission error and set up some aliases to make our life a bit easier.


```yaml
---
- hosts: all
  become: true
  tasks:
    - name: Add IP address of all hosts to all hosts
      lineinfile:
        dest: /etc/hosts
        regexp: '.*{{ item }}$'
        line: "{{ hostvars[item].private_ip }} {{item}}"
        state: present
      when: hostvars[item].private_ip is defined
      with_items: "{{ groups.all }}"

- hosts: all
  roles:
    - role: istvano.microk8s
      vars:
        microk8s_plugins:
          istio: true
          cert-manager: true
          ingress: true
          dns: true
          ha-cluster: true

- hosts: all
  tasks:
    - name: Add users to microk8s group
      ansible.builtin.user:
        name: ubuntu
        group: microk8s

- hosts: all
  become: true
  vars:
    bashrc: /etc/bash.bashrc
  tasks:
    - name: k alias
      lineinfile:
        path: "{{ bashrc }}"
        line: alias k='microk8s kubectl'
    - name: kubectl alias
      lineinfile:
        path: "{{ bashrc }}"
        line: alias kubectl='microk8s kubectl'
    - name: helm alias
      lineinfile:
        path: "{{ bashrc }}"
        line: alias helm='microk8s helm'
```

After appling the playbook, you should have fully configured cluster.

```sh
ansible-galaxy install -r requirements.yml
ansible-playbook playbook.yaml -i inventory.yaml
```

When it succeeds, you should be able to SSH into the control plane and see all nodes in the cluster.


```sh
ubuntu@control-plane:~$ k get nodes
NAME            STATUS   ROLES    AGE   VERSION
worker-1        Ready    <none>   13d   v1.27.8
worker-2        Ready    <none>   13d   v1.27.8
control-plane   Ready    <none>   13d   v1.27.8
worker-0        Ready    <none>   13d   v1.27.8
```

The full source code will available on [github](https://github.com/bkielbasa/oci-k8s). If there are any better ways to do some operations, pls let me know in the comments section below or open a PR.

I hope you enjoyed the article and have a running kubernetes cluster **for free**!
