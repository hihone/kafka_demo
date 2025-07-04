#!/bin/sh
# shellcheck disable=SC2063
# shellcheck disable=SC1001

# 定义全局变量
current_branch=$(git branch --show-current|awk -F/ '{print $NF}')
current_user=$(whoami)
current_date=$(date +%Y%m%d)
current_time=$(date +%H%M%S)
current_datetime=$(date +%Y%m%d%H%M%S)

# 部署
function deploy() {
    PLATFORM=$1
    ENV=$2
    IS_INIT=0
    LATEST_TAG=$(git tag -l "v*$PLATFORM-$ENV" --sort=-v:refname | head -n 1)
    if [ -z "$LATEST_TAG" ];then
        LATEST_TAG="v1.0.0-$PLATFORM-$ENV"
        ((IS_INIT = 1))
    fi

    # 解析版本号的各个部分
    cleaned_tag="${LATEST_TAG//-internal-prod/}"
    cleaned_tag="${cleaned_tag//-internal-test/}"

    # 使用冒号连接（在Makefile中隐含完成）
    tag_with_colon="${cleaned_tag}"

    IFS='.' read -ra parts <<< "$tag_with_colon"
    MAJOR_VERSION="${parts[0]#v}"
    MINOR_VERSION="${parts[1]}"
    PATCH_VERSION="${parts[2]}"

    if (($IS_INIT == 0)); then
        if (($PATCH_VERSION >= 9)); then
            ((PATCH_VERSION = 0))
            ((MINOR_VERSION = MINOR_VERSION + 1))
        else
            ((PATCH_VERSION = PATCH_VERSION + 1))
        fi
        if (($MINOR_VERSION >= 9)); then
            ((MINOR_VERSION = 0))
            ((MAJOR_VERSION = MAJOR_VERSION + 1))
        fi
    fi

    echo "当前平台: $PLATFORM"
    echo "当前环境: $ENV"
  	echo "当前版本: $LATEST_TAG"
  	echo "主版本号: v$MAJOR_VERSION"
  	echo "次版本号: $MINOR_VERSION"
  	echo "修订号: $PATCH_VERSION"
  	NEW_TAG="v$MAJOR_VERSION.$MINOR_VERSION.$PATCH_VERSION-$PLATFORM-$ENV"
  	echo "new_tag: $NEW_TAG"

  	git tag "$NEW_TAG"
    git push origin -f "$NEW_TAG"
}

if [[ x$1 = xrelease ]];then
  deploy internal prod;
elif [[ x$1 = xtest ]];then
  deploy internal test;
fi
date '+%Y-%m-%d %H:%M:%S'