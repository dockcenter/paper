# PaperMC Automatically Built Docker Image

[![Watch paper release](https://github.com/dockcenter/paper/actions/workflows/watch-releases.yaml/badge.svg?branch=main&event=schedule)](https://github.com/dockcenter/paper/actions/workflows/watch-releases.yaml)
[![GitHub](https://img.shields.io/github/license/dockcenter/paper?color=informational)](https://github.com/dockcenter/paper/blob/main/LICENSE)

This is a PaperMC docker image with optimized flag provided by official [docs](https://docs.papermc.io/paper/aikars-flags).

We use [GitHub Actions](https://github.com/dockcenter/paper/actions) to track PaperMC builds and automatically build Docker image.

## What is PaperMC?

Paper is a high performance fork of the Spigot Minecraft Server that aims to fix gameplay and mechanics inconsistencies as well as to improve performance. 
Paper contains numerous features, bug fixes, exploit preventions and major performance improvements not found in Spigot.

For more information, please reach to [PaperMC official documentation](https://docs.papermc.io/paper/getting-started).

![PaperMC](assets/paper.png)

## How to use this image

### Start a PaperMC server

With this image, you can create a new PaperMC Minecraft server with one command (note that running said command indicates agreement to the [Minecraft EULA](https://www.minecraft.net/en-us/eula)). 
Here is an example:

```bash
sudo docker run -p 25565:25565 dockcenter/paper
```

While this command will work just fine in many cases, it is only the bare minimum required to start a functional server and can be vastly improved by specifying some options.

## How to extend this image

There are many ways to extend the `dockcenter/paper` image. Without trying to support every possible use case, here are just a few that we have found useful.

### Environment Variables

The `dockcenter/paper` image uses several environment variables which are easy to miss. 
`JAVA_MEMORY` environment variable is not required, but it is highly recommended to set an appropriate value according to your usage.

#### `JAVA_MEMORY`

This variable is not required, but is highly recommended.
By setting this value, you set the java `-Xms` and `Xmx` flag. 
For more information about JVM memory size, refer to this [Oracle guide](https://docs.oracle.com/cd/E21764_01/web.1111/e13814/jvm_tuning.htm#PERFM160).

Default: `6G`

#### `JAVA_FLAGS`

This optional environment variable is used in conjunction with `JAVA_MEMORY` to provide additional java flag. 
We use [PaperMC officially recommended value](https://docs.papermc.io/paper/aikars-flags) as the default value.

Default: `-XX:+UseStringDeduplication -XX:+UseG1GC -XX:+ParallelRefProcEnabled -XX:MaxGCPauseMillis=200 -XX:+UnlockExperimentalVMOptions -XX:+DisableExplicitGC -XX:+AlwaysPreTouch -XX:G1NewSizePercent=30 -XX:G1MaxNewSizePercent=40 -XX:G1HeapRegionSize=8M -XX:G1ReservePercent=20 -XX:G1HeapWastePercent=5 -XX:G1MixedGCCountTarget=4 -XX:InitiatingHeapOccupancyPercent=15 -XX:G1MixedGCLiveThresholdPercent=90 -XX:G1RSetUpdatingPauseTimePercent=5 -XX:SurvivorRatio=32 -XX:+PerfDisableSharedMem -XX:MaxTenuringThreshold=1 -Dusing.aikars.flags=https://mcflags.emc.gs -Daikars.new.flags=true"`

### Volume

The server data is stored in `/data` folder, and we create a volume for you.
To use your host directory to store data, please mount volume by adding the following options:

Using volume:
```bash
-v <my_volume_name>:/data
```

Using bind mount:
```bash
-v </path/to/folders>:/data
```

## LICENSE

Be careful using this container image as you must meet the obligations and conditions of the [Minecraft EULA](https://www.minecraft.net/en-us/eula) as not doing so will be subject you or your organization to penalty under US Federal and International copyright law.

The code for the [project](https://github.com/dockcenter/paper) that builds the [`dockcenter/paper`](https://hub.docker.com/r/dockcenter/paper) image and pushes it to Docker Hub is distributed under the [MIT License](https://github.com/dockcenter/paper/blob/main/LICENSE).

Please, don't confuse the two licenses.