#!/bin/bash


#用来从环境变量中获取配置信息，将其写入到配置文件
#默认值已在config层完成，这里不再重复
set -e

APP_ID=${APP_ID:-""}
APP_SECRET=${APP_SECRET:-""}
APP_ENCRYPT_KEY=${APP_ENCRYPT_KEY:-""}
APP_VERIFICATION_TOKEN=${APP_VERIFICATION_TOKEN:-""}
BOT_NAME=${BOT_NAME:-""}
OPENAI_KEY=${OPENAI_KEY:-""}
HTTP_PORT=${HTTP_PORT:-""}
HTTPS_PORT=${HTTPS_PORT:-""}
USE_HTTPS=${USE_HTTPS:-""}
CERT_FILE=${CERT_FILE:-""}
KEY_FILE=${KEY_FILE:-""}
API_URL=${API_URL:-""}
HTTP_PROXY=${HTTP_PROXY:-""}
CONFIG_PATH=${CONFIG_PATH:-"config.yaml"}



# modify content in config.yaml
if [ "$APP_ID" != "" ] ; then
    sed -i "2c   APP_ID: $APP_ID" $CONFIG_PATH
else
    echo -e "\033[31m[Warning] You need to set APP_ID before running!\033[0m"
fi

if [ "$APP_SECRET" != "" ] ; then
    sed -i "3c   APP_SECRET: $APP_SECRET" $CONFIG_PATH
else
    echo -e "\033[31m[Warning] You need to set APP_SECRET before running!\033[0m"
fi

if [ "$APP_ENCRYPT_KEY" != "" ] ; then
    sed -i "4c   APP_ENCRYPT_KEY: $APP_ENCRYPT_KEY" $CONFIG_PATH
else
    echo -e "\033[31m[Warning] You need to set APP_ENCRYPT_KEY before running!\033[0m"
fi

if [ "$APP_VERIFICATION_TOKEN" != "" ] ; then
    sed -i "5c   APP_VERIFICATION_TOKEN: $APP_VERIFICATION_TOKEN" $CONFIG_PATH
else
    echo -e "\033[31m[Warning] You need to set APP_VERIFICATION_TOKEN before running!\033[0m"
fi

if [ "$BOT_NAME" != "" ] ; then
    sed -i "7c   BOT_NAME: $BOT_NAME" $CONFIG_PATH
else
    echo -e "\033[31m[Warning] You need to set BOT_NAME before running!\033[0m"
fi

if [ "$OPENAI_KEY" != "" ] ; then
    sed -i "9c   OPENAI_KEY: $OPENAI_KEY" $CONFIG_PATH
else
    echo -e "\033[31m[Warning] You need to set OPENAI_KEY before running!\033[0m"
fi


# 以下为可选配置
if [ "$HTTP_PORT" != "" ] ; then
sed -i "11c HTTP_PORT: $HTTP_PORT" $CONFIG_PATH
fi

if [ "$HTTPS_PORT" != "" ] ; then
sed -i "12c HTTPS_PORT: $HTTPS_PORT" $CONFIG_PATH
fi

if [ "$USE_HTTPS" != "" ] ; then
sed -i "13c USE_HTTPS: $USE_HTTPS" $CONFIG_PATH
fi

if [ "$CERT_FILE" != "" ] ; then
sed -i "14c CERT_FILE: $CERT_FILE" $CONFIG_PATH
fi

if [ "$KEY_FILE" != "" ] ; then
sed -i "15c KEY_FILE: $KEY_FILE" $CONFIG_PATH
fi

if [ "$API_URL" != "" ] ; then
    sed -i "17c   API_URL: $API_URL" $CONFIG_PATH
fi

if [ "$HTTP_PROXY" != "" ] ; then
    sed -i "19c   HTTP_PROXY: $HTTP_PROXY" $CONFIG_PATH
fi

echo -e "\033[32m[Success] Configuration file has been generated!\033[0m"

/dist/feishu_chatgpt
