#!/bin/bash

# ./start.sh m3u8 http://url.com/video.m3u8

if [ "$1" = "m3u8" ]; then
  echo "Running m3u8 downloading demo"
  if [ ! -n "$2" ]; then
    echo "We need a m3u8 url"
    exit 1
  fi

  target="m3u8.ts"
  if [ -n "$3" ]; then
    target=$3
  fi
  make m3u8-simple PATH_INFO="$2" TARGET="$target"
fi
