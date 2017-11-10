#!/bin/bash

# ./start.sh m3u8 http://url.com/video.m3u8

if [ "$1" = "m3u8" ]; then
  echo "Running m3u8 downloading demo"
  if [ -n "$2" ]; then
    make m3u8-simple PATH_INFO="$2"
  else
    echo "We need a m3u8 url"
  fi
fi
