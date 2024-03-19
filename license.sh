#!/bin/bash

LICENSE="// Copyright (c) Microsoft Corporation.\n// Licensed under the MIT license.\n\n"

for file in $(find . -name '*.go'); do
    printf "%b%s" "$LICENSE" "$(cat $file)" > temp && mv temp $file
done
