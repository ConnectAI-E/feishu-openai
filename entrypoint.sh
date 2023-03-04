#!/bin/bash
set -e

APP_ID=${APP_ID:-""}
APP_SECRET=${APP_SECRET:-""}
APP_ENCRYPT_KEY=${APP_ENCRYPT_KEY:-""}
APP_VERIFICATION_TOKEN=${APP_VERIFICATION_TOKEN:-""}
BOT_NAME=${BOT_NAME:-""}
OPENAI_KEY=${OPENAI_KEY:-""}
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
fi

if [ "$APP_VERIFICATION_TOKEN" != "" ] ; then
    sed -i "5c   APP_VERIFICATION_TOKEN: $APP_VERIFICATION_TOKEN" $CONFIG_PATH
else
    echo -e "\033[31m[Warning] You need to set APP_VERIFICATION_TOKEN before running!\033[0m"
fi

if [ "$BOT_NAME" != "" ] ; then
    sed -i "7c   BOT_NAME: $BOT_NAME" $CONFIG_PATH
fi


if [ "$OPENAI_KEY" != "" ] ; then
    sed -i "9c   OPENAI_KEY: $OPENAI_KEY" $CONFIG_PATH
else
    echo -e "\033[31m[Warning] You need to set OPENAI_KEY before running!\033[0m"
fi

/dist/feishu_chatgpt
