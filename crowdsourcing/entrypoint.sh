#!/bin/bash
set -e

EXECUTABLE=$1
shift


exec /app/$EXECUTABLE
