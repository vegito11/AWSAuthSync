#!/bin/bash -e

function format_cygwin_path(){
    CYG_PATH=$1
    if [[ $CYG_PATH == "/cygdrive"* ]];then
        DRIVE_PATH=$( echo $CYG_PATH | cut -d"/" -f3 | tr [a-z] [A-Z] )
        WORKSPACE_PATH=$(echo "$CYG_PATH" | cut -b 13-)
        CYG_FMT_PATH="${DRIVE_PATH}:/${WORKSPACE_PATH}"
    else
        echo "Path does not container cygwin prefix $CYG_PATH"
    fi
}


CODE_PARENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)"

cd $CODE_PARENT_DIR
pwd

PROJECT_MODULE="github.com/vegito11/AWSAuthSync"
IMAGE_NAME="kubernetes-codegen:latest"

CUSTOM_RESOURCE_GRP_NAME="vegito11.io"
CUSTOM_RESOURCE_VERSION="v1beta"

echo "Building codegen Docker image..."
# docker build -f "hacks/Dockerfile.generator" -t "${IMAGE_NAME}" .

cmd="./generate-groups.sh all \
    "$PROJECT_MODULE/pkg/client" \
    "$PROJECT_MODULE/pkg/apis" \
    $CUSTOM_RESOURCE_GRP_NAME:$CUSTOM_RESOURCE_VERSION --go-header-file hack/boilerplate.go.txt"

echo "Generating client codes..."

format_cygwin_path $(pwd)
echo " Debugging - docker run --rm -it -v ${CYG_FMT_PATH}:/go/src/${PROJECT_MODULE} ${IMAGE_NAME} sh "

docker run --rm -v "${CYG_FMT_PATH}:/go/src/${PROJECT_MODULE}" "${IMAGE_NAME}" $cmd

echo "Generating CRD manifests ${PROJECT_MODULE}/pkg/apis/${CUSTOM_RESOURCE_GRP_NAME}/${CUSTOM_RESOURCE_VERSION}"

crd_cmd="
./controller-gen paths=${PROJECT_MODULE}/pkg/apis/${CUSTOM_RESOURCE_GRP_NAME}/${CUSTOM_RESOURCE_VERSION} \
crd:crdVersions=v1 output:crd:artifacts:config=manifests"

$crd_cmd