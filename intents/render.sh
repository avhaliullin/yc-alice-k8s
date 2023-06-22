#!/bin/bash

export NAMESPACE="неймспейс|нейм спейс|нам спейс|namespace|name space"
export NAMESPACE_GEN="неймспейса|нейм спейса|нам спейса|namespace|name space"
export NAMESPACE_ADPOS="неймспейсе|нейм спейсе|нам спейсе|namespace|name space"
export NAMESPACE_PLUR="неймспейсы|нейм спейсы|нам спейсы|namespace|name space"
export NAMESPACE_GEN_PLUR="неймспейсов|нейм спейсов|нам спейсов|namespace|name space"
export K8S="кубернетис|купернетис|кубер|к8с"
export DEPLOY="дипло|деплой|диплой|диплей|деплоймент|диплоймент|диплеймент|диплоймант"
export DEPLOY_GEN="деплоя|диплоя|диплея|деплоймента|диплоймента|диплеймента|диплойманта"
export DEPLOY_ADPOS="деплое|диплое|диплее|деплойменте|диплойменте|диплейменте|диплойманте"

mkdir -p "dist"

for file in *.txt
do
  cat "$file" | envsubst '${NAMESPACE}${NAMESPACE_GEN}${NAMESPACE_ADPOS}${NAMESPACE_PLUR}${NAMESPACE_GEN_PLUR}${K8S}${DEPLOY}${DEPLOY_GEN}${DEPLOY_ADPOS}' > "dist/$file"
done
