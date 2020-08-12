# dyndns53

A simple program that will update AWS route 53's records for a domain with the local external IP address.

Basically a golang version of [ddclient](https://github.com/ddclient/ddclient) but only targeting AWS Route 53 (at the time of writing).

Relies on https://ipinfo.io/ip to retrieve the current external IP address, so where this is run needs 443 access to the net. The external service address can be changed using the flag `--myip-service`.

Patterned after and inspired by the following two bash equivalents (I wrote in go because *when all you have is a hammer...* etc etc, and I wanted to try out AWS manipulation with go):

- https://willwarren.com/2014/07/03/roll-dynamic-dns-service-using-amazon-route53/
- https://medium.com/@avishayil/dynamic-dns-using-aws-route-53-60a2331a58a4

Also, slight advantage, this tool does not require `dig` or the AWS cli to be installed. It *DOES* require credentials to be configured however (which is easy to do with the cli, but can be done in a few different ways)

After configuring credentials (see [this page on guidance for how](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials)), run with no args or `-h` to see configuration options.

## Required packages to compile

Just the following go gets:

- `go get github.com/aws/aws-sdk-go/aws`
- `go get github.com/aws/aws-sdk-go/service/route53`
