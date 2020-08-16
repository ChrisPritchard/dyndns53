# dyndns53

A simple program that will update AWS route 53's records for a domain with the local external IP address.

Basically a golang version of [ddclient](https://github.com/ddclient/ddclient) but only targeting AWS Route 53 (at the time of writing).

Relies on https://ipinfo.io/ip to retrieve the current external IP address, so where this is run needs 443 access to the net. The external service address can be changed using the flag `--myip-service`.

Patterned after and inspired by the following two bash equivalents (I wrote in go because *when all you have is a hammer...* etc etc, and I wanted to try out AWS manipulation with go):

- https://willwarren.com/2014/07/03/roll-dynamic-dns-service-using-amazon-route53/
- https://medium.com/@avishayil/dynamic-dns-using-aws-route-53-60a2331a58a4

Also, slight advantage, this tool does not require `dig` or the AWS cli to be installed. It *DOES* require credentials to be configured however (which is easy to do with the cli, but can be done in a few different ways)

After configuring credentials (see [this page on guidance for how](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials)), run with no args or `-h` to see configuration options, or check the setup example below.

## Required packages to compile

Just the following go gets:

- `go get github.com/aws/aws-sdk-go/aws`
- `go get github.com/aws/aws-sdk-go/service/route53`

## Steps to setup as a cronjob on linux, safely (?), no aws cli required

First, log into the AWS portal with your management user. We need to create an IAM user specifically for dyndns53, and gather some info.

### In AWS

- go to IAM and manage users, and create a new user
    - the user will be a 'programmatic user', and does not need access to the portal
    - when it asks for permissions, select 'assign directly'
    - create a new policy to assign to the user
        - the new policy should have a single permission: `route53:ChangeResourceRecordSets`
        - under resources, either select the specific hosted zone or use all, if you might use dyndns53 for more than one domain: `arn:aws:route53:::hostedzone/*`
    - the policy should look like below in json:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": "route53:ChangeResourceRecordSets",
            "Resource": "arn:aws:route53:::hostedzone/*"
        }
    ]
}
```

- at the end of the user creation process, keep a copy of the **access key** and **secret access key**
- finally, go to route53 > hosted zones, and make a note of the **hosted zone id** and **domain name** of the domains you want to update.

### On your linux machine

- on linux, create a new user. This user will run the tool, and will have the aws credentials, but does not need to log in:

    `sudo adduser dyndns53-user --disabled-login`

- navigate to the user's home directory, and put your compiled copy of dyndns53 in there.
- create a new folder `.aws`
- in that folder, create a file named `credentials` and stick in it the following content (updating with the keys for your new IAM user):

```yaml
[default]
aws_access_key_id = <access key>
aws_secret_access_key = <secret access key>
```

note you dont need to quote wrap these keys.

- back in the new users home dir, create a new bash file named `run-dyndns53.sh`, and add the following content (again, updating with your values for hosted zone, target domain):

```bash
#!/bin/bash

/home/dyndns53-user/dyndns53 --hosted-zone-id <hosted zone id> --target-domain <domain name>
```

- run the following commands as sudo to give all this to dyndns53-user:

```bash
sudo chown -R dyndns53-user .aws
sudo chown dyndns53-user ./dyndns53
sudo chown dyndns53-user ./run-dyndns53.sh
sudo -u dyndns53-user chmod +x ./run-dyndns53.sh
```

- at this point you can test all is tickety boo by running `sudo -u dyndns53-user ./run-dyndns53`. The output should be as expected from naked running the tool, and should update AWS as this is the first run.

- run the following to open/create a cron table for the user: `sudo crontab -e -u dyndns53-user`
- in the editor, add the following to setup a job to run every five minutes: `*/5 * * * * /home/dyndns53-user/run-dyndns53.sh`

And that's it! Wait five minutes, then run `sudo grep CRON /var/log/syslog` to see if the tool is running. If you like, you could update the run script to use the `-current-ip` arg so you can overwrite it to something manual, like `127.0.0.1`, and see this be reflected in aws, before stripping back to sourcing the accurate value. NOICE!