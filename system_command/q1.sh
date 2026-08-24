#!/usr/bin/env bash

awk '$2>25{print $1}' input.txt
