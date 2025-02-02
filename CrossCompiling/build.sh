#!/bin/bash

platforms=("windows/amd64" "linux/amd64" "darwin/amd64" "darwin/arm64")

for platform in "${platforms[@]}"
do
    osarch=(${platform//\// })
    GOOS=${osarch[0]}
    GOARCH=${osarch[1]}
    output_name="main-${GOOS}-${GOARCH}"

    if [ "$GOOS" == "windows" ]; then
        output_name+=".exe"
    fi

    echo "Building for $GOOS/$GOARCH..."
    env CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name
done

echo "Done!"
