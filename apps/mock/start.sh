#!/bin/bash

printf "./build/nsqd\n"

printf "./build/mock\n"

printf "./build/nsq_to_http --topic=test --channel=test_chan --nsqd-tcp-address 127.0.0.1:4150 --http-client-connect-timeout 50ms --consumer-opt low_rdy_idle_timeout,1s --consumer-opt rdy_redistribute_interval,100ms --consumer-opt max_backoff_duration,100ms --max-in-flight 1000 -n 150 --status-every 100 --post=http://localhost:8080\n"

print "sudo tcpkill -i any port 4150"
print "sudo tcpkill -i any port 8080"
