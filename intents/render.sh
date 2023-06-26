#!/bin/bash

export NAMESPACE="эмспейс|неймспейс|нейм спейс|нам спейс|namespace|name space"
export NAMESPACE_GEN="эмспейса|неймспейса|нейм спейса|нам спейса|namespace|name space"
export NAMESPACE_ADPOS="эмспейсе|неймспейсе|нейм спейсе|нам спейсе|namespace|name space"
export NAMESPACE_PLUR="эмспейсы|неймспейсы|нейм спейсы|нам спейсы|namespace|name space"
export NAMESPACE_GEN_PLUR="эмспейсов|неймспейсов|нейм спейсов|нам спейсов|namespace|name space"
export K8S="кубернетис|купернетис|кубер|к8с|купер"
export DEPLOY="дипло|деплой|диплой|диплей|деплоймент|диплоймент|диплеймент|диплоймант"
export DEPLOY_GEN="деплоя|диплоя|диплея|деплоймента|диплоймента|диплеймента|диплойманта"
export DEPLOY_ADPOS="деплое|диплое|диплее|деплойменте|диплойменте|диплейменте|диплойманте"

mkdir -p "dist"

for file in *.txt
do
  cat "$file" | envsubst '${NAMESPACE}${NAMESPACE_GEN}${NAMESPACE_ADPOS}${NAMESPACE_PLUR}${NAMESPACE_GEN_PLUR}${K8S}${DEPLOY}${DEPLOY_GEN}${DEPLOY_ADPOS}' > "dist/$file"
done
