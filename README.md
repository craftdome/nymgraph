# Nymgraph
A graphical chat client for NYM (nym-client)

## Developer's mixnodes

[![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fexplorer.nymtech.net%2Fapi%2Fv1%2Fmix-node%2F895&query=mix_node.identity_key&style=plastic&logo=numpy&logoColor=white&label=Advanced%20Engineering%201&color=%23136401&cacheSeconds=60&link=https%3A%2F%2Fexplorer.nymtech.net%2Fnetwork-components%2Fmixnode%2F895)](https://explorer.nymtech.net/network-components/mixnode/895)
[![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fexplorer.nymtech.net%2Fapi%2Fv1%2Fmix-node%2F895&query=stake_saturation&style=plastic&logo=myspace&logoColor=white&label=Stake&cacheSeconds=60&link=https%3A%2F%2Fexplorer.nymtech.net%2Fnetwork-components%2Fmixnode%2F895)](https://explorer.nymtech.net/network-components/mixnode/895)
[![Dynamic JSON Badge](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fexplorer.nymtech.net%2Fapi%2Fv1%2Fmix-node%2F895&query=mix_node.version&style=plastic&logo=git&logoColor=white&label=Version&cacheSeconds=60&link=https%3A%2F%2Fexplorer.nymtech.net%2Fnetwork-components%2Fmixnode%2F895)](https://explorer.nymtech.net/network-components/mixnode/895)

# Compiling from source (Ubuntu & Debian (amd64))
## Step 1. Installing Golang
1. Visit the official page [go.dev](https://go.dev/doc/install) to download and install Go.

## Step 2. Installing the required packages for compilation
1. Packages.
```bash
apt install git build-essential libxinerama-dev libgl1-mesa-dev xorg-dev libx11-dev pkg-config
```

2. Setting up PKG_CONFIG_PATH for pkgconfig dir.
```sh
export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/usr/lib/pkgconfig
```

## Step 3. Download and compile nymgraph

1. Clonning nymgraph project via git clone.
```sh
git clone https://github.com/craftdome/nymgraph.git && cd nymgraph
```

2. Syncronize go dependencies (downloading all imports)
```sh
go mod tidy
```

3. Compiling to output dir `./bin`.
```sh
CGO_ENABLED=1 go build -o ./bin/nymgraph-amd64 -ldflags="-s -w" -trimpath github.com/craftdome/nymgraph/cmd/app
```

4. Run a result.

![изображение](https://github.com/Tyz3/nymgraph/assets/21179689/2ce9e4f1-117f-4475-992f-e8d90f3bc7a1)

![изображение](https://github.com/Tyz3/nymgraph/assets/21179689/a0f9b6b3-fb9b-4ef3-b85c-c905b47fd725)


# Screenshots (Windows 11)

![изображение](https://github.com/Tyz3/nymgraph/assets/21179689/b36347f8-c673-4bec-a79b-32eed22b7115) 

|||
|---|---|
|![изображение](https://github.com/Tyz3/nymgraph/assets/21179689/2f7595c0-35c9-4817-909d-9c9099245f6d)|![изображение](https://github.com/Tyz3/nymgraph/assets/21179689/cc16fd7d-0edf-42b2-a943-7065796419fa)|
|![изображение](https://github.com/Tyz3/nymgraph/assets/21179689/501f5e6d-aa6c-4f14-9e19-f99c5f874477)|![изображение](https://github.com/Tyz3/nymgraph/assets/21179689/a60516b4-aa3d-477f-a24d-56fe648e55ed)|
