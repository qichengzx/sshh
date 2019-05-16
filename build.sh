#!/usr/bin/env bash

package=sshh
version="0.2.1"

platforms=("darwin/amd64" "linux/amd64" "linux/arm" "linux/386")

echo "build..."

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })

    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    output_name="${package}-${GOOS}-${GOARCH}_${version}"

    echo "build ${output_name}"

    CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o "./releases/${output_name}/${package}" main.go

    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi

    cd ./releases

    zip -r -m "./${output_name}.zip" "./${output_name}"

    cd ../

done

echo "build done"
