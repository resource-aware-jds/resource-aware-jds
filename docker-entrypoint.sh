#!/bin/sh

nohup dockerd &
sleep 5
$1
