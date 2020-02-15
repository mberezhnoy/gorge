#!/bin/sh -l

printenv
ls -la /github
ls -la /github/home
ls -la /github/workspace

mkdir -p /github/workspace/libproj
cp /usr/lib/libproj.*a /github/workspace/libproj/
