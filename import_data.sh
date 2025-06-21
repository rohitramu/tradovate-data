#!/bin/bash
set -e -o pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

mkdir -p $SCRIPT_DIR/.data
mkdir -p $SCRIPT_DIR/.db

go run tradovatedataimport $SCRIPT_DIR/.data/ $SCRIPT_DIR/.db/data.db
