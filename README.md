# gomake
[![Build Status](https://dev.azure.com/brad-jones/gomake/_apis/build/status/brad-jones.gomake?branchName=master)](https://dev.azure.com/brad-jones/gomake/_build/latest?definitionId=1?branchName=master)
[![codecov](https://codecov.io/gh/brad-jones/gomake/branch/master/graph/badge.svg)](https://codecov.io/gh/brad-jones/gomake)

A cross platform build tool / task runner that scales.

## Instalation

### Direct download

Go to <https://github.com/brad-jones/gomake/releases> and download the archive for your Operating System, extract the gomake binary and and add it to your `$PATH`.

#### Curl Bash

```
curl -L https://github.com/brad-jones/gomake/releases/download/v2.1.2/gomake_linux_amd64.tar.gz -o- | sudo tar -xz -C /usr/bin gomake
```

### RPM package

```
sudo rpm -i https://github.com/brad-jones/gomake/releases/download/v2.1.2/gomake_linux_amd64.rpm
```

### DEB package

```
curl -sLO https://github.com/brad-jones/gomake/releases/download/v2.1.2/gomake_linux_amd64.deb && sudo dpkg -i gomake_linux_amd64.deb && rm gomake_linux_amd64.deb
```

### Homebrew

<https://brew.sh>

```
brew install brad-jones/tap/gomake
```

### Scoop

<https://scoop.sh>

```
scoop bucket add brad-jones https://github.com/brad-jones/scoop-bucket.git;
scoop install gomake;
```
