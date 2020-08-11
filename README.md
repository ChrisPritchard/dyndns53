# dyndns53

A simple program that will update AWS route 53's records for a domain with the local external IP address.

Basically a golang version of [ddclient](https://github.com/ddclient/ddclient) but only targeting AWS Route 53 (at the time of writing), and in Golang.

Relies on https://ipinfo.io/ip to retrieve the current external IP address, so where this is run needs 443 access to the net. The external service address can be changed using the flag `--myip-service`.

Patterned after and inspired by the following two bash equivalents (I wrote in go because *when all you have is a hammer...* etc etc, and I wanted to try out AWS manipulation with go):

- https://willwarren.com/2014/07/03/roll-dynamic-dns-service-using-amazon-route53/
- https://medium.com/@avishayil/dynamic-dns-using-aws-route-53-60a2331a58a4

## Usage

The samples assume you have compiled/built this as `dyndns` for your platform (e.g. dyndns.exe on windows would still work as below from powershell).

