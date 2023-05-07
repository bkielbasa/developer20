---
title: "Terraform can do more than you think"
publishdate: 2023-05-05
categories: 
    - Devops
tags:
  - terraform
---

The article showcases different use cases of Terraform, a popular Infrastructure as Code tool. It highlights how versatile Terraform is by demonstrating its capability to manage SSH keys, GitHub repositories, Spotify playlists, DNS, and even ordering pizza from Domino's. Additionally, it features some other providers like HTTP and Infracost.

## SSH keys

I'll start with one of my favorits that is generating SSH keys that can be used to connect to your servers.

```terraform
resource "tls_private_key" "ssh" {
  algorithm = "RSA"
  rsa_bits  = 4096
}
```

And output it:

```terraform
output "ssh_private_key" {
  value = tls_private_key.ssh.private_key_openssh
}
```

The way to get the private key is to run following commands

```sh
cat terraform.tfstate | jq -r '.resources[] | select (.module == "module.networking" and .type == "tls_private_key" and .name == "ssh") | .instances[0].attributes.private_key_openssh' > ~/.ssh/k8s.key
chmod 400 ~/.ssh/k8s.key
ssh -i ~/.ssh/k8s.key ubuntu@<public_ip>
```

Or if you use the `output` feature:

```sh
terraform output -json | jq -r '.ssh_private_key.value' > ~/.my_ssh_key
```

The key is stored inside the Terraform state so you can write a script that will store it into your `.ssh` folder automatically.

```terraform
resource "local_file" "ssh_key" {
  content  = tls_private_key.ssh.private_key_pem
  filename = "~/.ssh/server"
}
```

I'm wondering if someone use Terraform this way.

## GithHub

I don't know if you know that you can manage github repositories and organizations as well

```terraform
terraform {
  required_providers {
    github = {
      source  = "integrations/github"
      version = "~> 5.0"
    }
  }
}

# Configure the GitHub Provider
provider "github" {}

# Add a user to the organization
resource "github_membership" "membership_for_user_x" {
  # ...
}
```

It looks handy when you generate repositories or want to give developers the abbility to change their's repositories settings without giving them admin rights. I found it useful when I wanted to change repositories setting to delete a branch after merge. I could do it everywhere in a quite simple way.

Link: https://registry.terraform.io/providers/integrations/github/latest/docs

## Manage spotify playlist

That one sounds funny and, to be honest, I cannot find a good usage for but it only shows how flexible terraform is.

https://github.com/conradludgate/terraform-provider-spotify

```terraform
resource "spotify_playlist" "playlist" {
  name        = "My playlist"
  description = "My playlist is so awesome"
  public      = false

  tracks = flatten([
    data.spotify_track.overkill.id,
    data.spotify_track.blackwater.id,
    data.spotify_track.overkill.id,
    data.spotify_search_track.search.tracks[*].id,
  ])
}

data "spotify_track" "overkill" {
  url = "https://open.spotify.com/track/4XdaaDFE881SlIaz31pTAG"
}
data "spotify_track" "blackwater" {
  spotify_id = "4lE6N1E0L8CssgKEUCgdbA"
}

data "spotify_search_track" "search" {
  name   = "Somebody Told Me"
  artist = "The Killers"
  album  = "Hot Fuss"
}

output "test" {
  value = data.spotify_search_track.search.tracks
}
```

## Terratest

You can write test for your infrastructure. Tests are written in Go.

Link: https://github.com/gruntwork-io/terratest

PS. If you want me to write a tutorial about it, let me know in commends section below.


## infracost

Infracost is a tool that can tell you how much money you'll spend when you apply your plan. Works great with AWD and GCP.

Link: https://github.com/infracost/infracost

## Dominos

You can order pizza using terraform!

Link: https://ndmckinley.github.io/terraform-provider-dominos/

PS. A funny quote

> As far as I know, there is no programmatic way to destroy an existing pizza. terraform destroy is implemented on the client side, by consuming the pizza.

I don't know if this works because I live in Poland but maybe some geeks will find it useful, somehow.

## HTTP

This provider looks like a very handy one when you have some data exposed in an HTTP API.

```terraform
data "http" "example" {
  url = "https://checkpoint-api.hashicorp.com/v1/check/terraform"

  # Optional request headers
  request_headers = {
    Accept = "application/json"
  }
}
```

Link: https://registry.terraform.io/providers/hashicorp/http/latest

It sounds like a good work-around for some information that are available using a HTTP API or it's easy to expose them this way.


## DNS

I didn't know I can manage my DNS using terraform but when you think about it - it make sense!

```terraform
# Configure the DNS Provider
provider "dns" {
  update {
    server        = "192.168.0.1"
    key_name      = "example.com."
    key_algorithm = "hmac-md5"
    key_secret    = "3VwZXJzZWNyZXQ="
  }
}

# Create a DNS A record set
resource "dns_a_record_set" "www" {
  # ...
}
```

If you use [cloudflare](https://registry.terraform.io/providers/cloudflare/cloudflare/latest/docs), [GoDaddy](https://github.com/n3integration/terraform-provider-godaddy) or any other DNS provider, you have one place to manage everything.

## ChatGPT

Another suprising example is the [provider for ChatGPT](https://registry.terraform.io/providers/develeap/chatgpt/latest/docs).

```terraform
resource "chatgpt_prompt" "example" {
  max_tokens = 256
  query      = "Who is the best cloud provider?"
}

output "example_result" {
  value = chatgpt_prompt.example.result
}
```

To be honest, I cannot see any good useges for it in the infrastructure part. If you have any idea how it can be useful, please let me know in the comments section below (only bad answers.)

## Minecraft

You can build your Minecraft map using terraform as well.

```terraform
resource "minecraft_block" "stone" {
  material = "minecraft:stone"

  position = {
    x = -198
    y = 66
    z = -195
  }
}
```

* https://registry.terraform.io/providers/fourplusone/jira/0.1.20

Do you want to write a tutorial for terratest or writing own providers?
