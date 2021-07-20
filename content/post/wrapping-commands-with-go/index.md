---
title: "Wrapping commands in Go"
publishdate: 2021-07-14

resources:
    - name: header
    - src: featured.jpg

categories: [Golang]
tags:
  - cli
  - aws
  - ec2
  - ssh
---

You can find a lot of articles about Go that describe general aspects of it. Including the content on this blog. Today, I decided to prepare something different. I’ll tell you about one of my tasks and I’ll show you how I resolved it using Go. I thought it'd be useful to show the `exec` package and to tell a bit about the `ssh` command and learn AWS EE2 a bit better.

AWS has a feature called [Amazon EC2 Instance Connect](https://aws.amazon.com/blogs/compute/new-using-amazon-ec2-instance-connect-for-ssh-access-to-your-ec2-instances/). You can use it to connect to your EC2 instance using an SSH client. The whole process has a few steps:

* get information about the instance (region, instance id, an availability zone)
* upload your public key to the EC2 instance
* Connect to the instance using the private key and SSH command

The problem we're solving today is automating this process. After uploading the SSH key we have 60 seconds to connect to the EC2 instance. If you connect to a lot of EC2 instances and you have to repeat the same steps over and over you want to automate it.

My goal was to create an `ssh` replacement that accepts the same parameters and behaves as a regular `ssh` command but automates the whole setup process.

The requirement is simple - the usage of the command should be as similar to the `ssh` command as possible. In a perfect world - it should be 100% replacement. Let’s try if we can achieve that. An example command will look like this:

```bash
./ec2-ssh HOSTNAME

# or
./ec2-ssh ec2-user@HOSTNAME # the default username is your OS user

# or
./ec2-ssh 192.168.0.1 -4 -k # accepts any parameter that ssh does
```

To upload our public key we need to know the availability zone, EC2 instance ID, user, and the public key itself. From the command parameters, we have the IP or the hostname.

Let’s start with the username, host, and public key. There’s a `-G` parameter in the `ssh` command that prints all the configurations after evaluating host and match blocks. We can call the `ssh` command with all the parameters provided by the user and add `-G`. Then, we can parse the output and read all the data from it. In other words, we want to read this command’s output.

```bash
./ec2-ssh 192.168.0.1 -4 -k # from this command

ssh -G 192.168.0.1 -4 -k # we'll translate to this
```

We have to call the `ssh` command as a subprocess and read its output. Go has an `exec` package that contains the `Cmd` struct. This struct represents an external command and can be used for this purpose.

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

The `Cmd` struct has `cmd.Stdout` field that’s the most important for us. It’s the place where we can forward the output of the command. This field accepts any `io.Writer` type so we put our buffer there. The next step is to put the parameters into the map from where we'll retrieve values.

```go
	res := map[string][]string{}

	scanner := bufio.NewScanner(buff)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		if len(parts) < 1 {
			continue
		}

		if _, exists := res[parts[0]]; !exists {
			res[parts[0]] = []string{}
			continue
		}

		res[parts[0]] = append(res[parts[0]], strings.Join(parts[1:], " "))
	}

	return res, nil
```

We go line by line and put the data on the map. In the next function, we get the required information from the map: we need its IPv4 as well as the username.

```go
func instanceInfoFromString(hostname, user string) (*instanceInfo, error) {
	info := &instanceInfo{
		username: user,
		host:     hostname,
	}

	err := info.resolveIP()
	if err != nil {
		return nil, err
	}
	return info, nil
}
```

We need the IP address because in later steps. We'll need it to filter out irrelevant EC2 instances.
The next step is to find the public key that `ssh` will use to connect us to the EC2 instance. A list of possible SSH keys is available under the `identityfile` key in our map. We iterate over every item and check if it exists. If yes, then we return it.

```go
func existingKey(paths []string) (string, error) {
	for _, path := range paths {
		path, err := expandHomeDirectoryTilde(path)
		if err != nil {
			return "", err
		}

		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			continue
		}

		return path, nil
	}

	return "", errors.New("cannot find any ssh key")
}
```

Every key's path starts (in general) with a tilde (`~`) that's means the user's home directory. We had to write a function that expands the tilde to a full path. Why? The tilde is expanded by your shell's `HOME` value. You can read how it works in more detail in [bash's docs](https://www.gnu.org/software/bash/manual/html_node/Tilde-Expansion.html) or this [SO answer](https://unix.stackexchange.com/questions/146671/does-always-equal-home/146697#146697). Let's get back to the code.

```go
	publicKey, err := getPublicKey(pk)
	if err != nil {
		return fmt.Errorf("cannot read the public key %s.pub. If you want to provide a custom key location, use the `-i` parameter", pk)
	}
```

In the listing below we attempt to read the public key. We need it to upload to the EC2 instance. This public key will be used to authenticate us. It means the EC2 instance has to know it before we'll attempt to connect to it. AWS will put our public key to `~/.ssh/authorized_keys` file. We have only 60 seconds to connect to the instance. For more details on how the SSH authorization works, you can [visit this description](https://www.ssh.com/academy/ssh/public-key-authentication).

We have almost everything we need to connect to the EC2 instance. The only thing missing is the AWS region. We have the requirement that we are ass `ssh` command compatible as possible we cannot just add another parameter to our command. Instead, we'll iterate over all regions and try to connect to every single instance. I know it's not the most optimal way. If you have any idea how I can improve it - let me know in the comments section below.

```go
func setupEC2Instance(ctx context.Context, instance *instanceInfo, publicKey, region string) (bool, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return false, fmt.Errorf("cannot get config for AWS: %w", err)
	}

	client := ec2.NewFromConfig(cfg)

	ec2Instance, err := findEC2Instance(ctx, client, instance)
	if err != nil {
		return false, err
	}

	if ec2Instance == nil {
		return false, nil
	}

	status, err := instanceStatus(ctx, client, *ec2Instance)
	if err != nil {
		return false, fmt.Errorf("cannot get the instance status: %w", err)
	}

	connect := ec2instanceconnect.NewFromConfig(cfg)
	out, err := connect.SendSSHPublicKey(ctx, &ec2instanceconnect.SendSSHPublicKeyInput{
		AvailabilityZone: status.AvailabilityZone,
		InstanceId:       ec2Instance.InstanceId,
		InstanceOSUser:   &instance.username,
		SSHPublicKey:     &publicKey,
	})

	if err != nil {
		return false, fmt.Errorf("cannot upload the public key: %w", err)
	}

	if !out.Success {
		return false, fmt.Errorf("unsuccessful uploaded the public key")
	}

	return true, nil
}
```

In the code above, we're configuring the AWS client, trying to find our EC2 instance in the selected region. If everything goes fine, we're uploading our public key. If it succeeds as well, we're ready to connect. Two functions are new here: `findEC2Instance` and `instanceStatus`.

The first one is quite obvious - it finds our EC2 instance using the IP address we retrieved earlier.

```go
func findEC2Instance(ctx context.Context, client *ec2.Client, info *instanceInfo) (*types.Instance, error) {
	resp, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   strp("private-ip-address"),
				Values: []string{info.ipAddress},
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("cannot contact with AWS API: %w", err)
	}

	for _, r := range resp.Reservations {
		for _, inst := range r.Instances {
			if *inst.PrivateIpAddress == info.ipAddress {
				return &inst, nil
			}
		}
	}
	return nil, nil
}
```

When we know that the instance exists and we have its reference, we can check its status and this is when the `instanceStatus` comes into play.

```go
func instanceStatus(ctx context.Context, client *ec2.Client, instance types.Instance) (types.InstanceStatus, error) {
	descResp, err := client.DescribeInstanceStatus(ctx, &ec2.DescribeInstanceStatusInput{
		InstanceIds: []string{*instance.InstanceId},
	})

	if err != nil {
		return types.InstanceStatus{}, err
	}

	status := descResp.InstanceStatuses[0]
	return status, nil
}
```

The `client.DescribeInstanceStatus` returns a few very valuable information for us: the instance's available zone and the instance's ID. Both values are required while uploading the SSH public key.

At this point, we are ready to connect to the EC2 instance! That's quite simple - had to execute the `ssh` command with all our parameters. We forward all the output to the standard output and do the same with the std input. Thanks to this, we'll be able to interact with the `ssh` command as usual.

```go
func connectToInstance(ctx context.Context, params []string) error {
	cmd := exec.CommandContext(ctx, "ssh", params...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// terminated by Control-C so ignoring
			if exiterr.ExitCode() == 130 {
				return nil
			}
		}

		return fmt.Errorf("error while connecting to the instance: %w", err)
	}

	return nil
}
```

And that's all! The whole source code is [available on Github](https://github.com/bkielbasa/ec2-ssh). From now, you can replace your `ssh` command with `ec2-ssh` while working with AWS EC2 instances. If you have any questions or suggestions, feel free to use the comments section below.
