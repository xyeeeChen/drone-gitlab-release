# Drone Gitlab release

Drone plugin for creating a GitLab release.

## Build

Build the binary

```sh
go build
```

Build the docker image

```sh
docker build -t drone-gitlab-release
```

## Run

```sh
$ ./drone-gitlab-release -h
NAME:
   Gitlab Release - Create a release to gitlab

USAGE:
   drone-gitlab-release [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --access_token value  Gitlab access token [$DRONE_ACCESS_TOKEN]
   --domain value        Gitlab domain name [$DRONE_DOMAIN]
   --repo value          The gitlab project ID or URL-encoded path of the project [$DRONE_REPO]
   --release value       The release name [$DRONE_RELEASE]
   --tag value           The tag name [$DRONE_TAG]
   --description value   The description of the release that supports Markdown [$DRONE_DESCRIPTION]
   --ref value           The release is created from ref and tagged with tag_name. It can be a commit SHA, another tag name, or a branch name. [$DRONE_REF]
   --assets value        Optional path for a direct asset link. [$DRONE_ASSETS]
   --help, -h            show help
```

## Usage

```sh
docker run -e DRONE_ACCESS_TOKEN="your-access-token" \
           -e DRONE_TAG="v0.0.1" \
           -e DRONE_DESCRIPTION="# your-description" \
           -e DRONE_REPO="example/your-repo" \
           -e DRONE_REF="master" \
           -e DRONE_DOMAIN="your-gitlab-domain"\
           -e DRONE_RELEASE="v0.0.1" \
           -e DRONE_ASSETS="your-assets-file-path" \
           gitlab-release
```