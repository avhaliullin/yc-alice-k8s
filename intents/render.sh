#!/bin/bash

export NAMESPACE="неймспейс|нейм спейс|namespace|name space"
export K8S="кубернетис|купернетис|кубер|к8с"

mkdir -p "dist"

for file in *.txt
do
  cat "$file" | envsubst '${NAMESPACE}${K8S}' > "dist/$file"
done
