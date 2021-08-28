#!/bin/sh
echo `which network-health-server` | entr -nr `which network-health-server` $@
