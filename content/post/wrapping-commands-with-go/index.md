---
title: "Wrapping commands in Go"
publishdate: 2021-04-29

resources:
    - name: header
    - src: featured.jpg

categories: [Golang]
tags:
  - cli
---

You can find a lot of articles about Go that describe general aspects of it. Including the content on this blog. Today, I decided to prepare something different. I’ll tell you about one of my tasks and I’ll show you how I resolved it.

AWS has a feature called [Amazon EC2 Instance Connect](https://aws.amazon.com/blogs/compute/new-using-amazon-ec2-instance-connect-for-ssh-access-to-your-ec2-instances/) that you can use to connect to the EC2 instance using SSH client. The whole process has a few steps:

* get information about the instance (region, instance id and availability zone)
* upload your public key to the EC2 instance
* connect to the instance using the private key and SSH command

The requirement is simple - the usage of the command should be as similar to `ssh` command as possible. In a perfect world - it should be 100% replacement. Let’s try if we can achieve that. An example command will look like this:

```bash
./ec2-ssh HOSTNAME

# or
./ec2-ssh ec2-user@HOSTNAME # the default username is your OS user

# or
./ec2-ssh 192.168.0.1 -4 -k # accepts any parameter that ssh does
```

To upload our public key we need to know the availability zone, EC2 instance ID,  user and the public key itself. From the command parameters we have the IP or the hostname.

Let’s start with the username, host and public key. There’s a `-G` parameter in the `ssh` command that prints all the configuration after evaluating host and match blocks. We can call the ssh command with all the parameters provided by the user and add `-G`. Then, we can parse the output and read all the data from it. In other words, we want to read this command’s output.

```bash
./ec2-ssh 192.168.0.1 -4 -k # from this command

ssh -G 192.168.0.1 -4 -k
```

We have to call the `ssh` command as a subprocess and read it’s output. Go has `exec` package that contains the `Cmd` struct. This stract represents an external command and can be used for this purpose.

```go
	cmd := exec.CommandContext(ctx, "ssh", args...)

	s := ""
	buff := bytes.NewBufferString(s)
	cmd.Stdout = buff
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return nil, err
	}
```

The `Cmd` struct has `cmd.Stdout` field that’s the most important for us. It’s the place where we can choose

ref: https://aws.amazon.com/blogs/compute/new-using-amazon-ec2-instance-connect-for-ssh-access-to-your-ec2-instances/
