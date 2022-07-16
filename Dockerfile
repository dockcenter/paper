FROM eclipse-temurin:17-jre-alpine AS builder

ARG DOWNLOAD_URL

WORKDIR /data

RUN apk add --no-cache wget

RUN wget -P /opt/papermc/ $DOWNLOAD_URL && \
    java -jar /opt/papermc/paper-*.jar

# Sign eula and remove logs folder
RUN sed -i 's/false/true/g' eula.txt && \
    rm -R logs/

FROM eclipse-temurin:17-jre-alpine

LABEL org.opencontainers.image.vendor="Dockcenter"
LABEL org.opencontainers.image.title="PaperMC"
LABEL org.opencontainers.image.description="Dockcenter PaperMC Docker image"
LABEL org.opencontainers.image.documentation="https://github.com/dockcenter/paper/blob/main/README.md"
LABEL org.opencontainers.image.authors="Chao Tzu-Hsien <danny900714@gmail.com>"
LABEL org.opencontainers.image.licenses="MIT"

ENV JAVA_MEMORY="6G"
ENV JAVA_FLAGS="-XX:+UseStringDeduplication -XX:+UseG1GC -XX:+ParallelRefProcEnabled -XX:MaxGCPauseMillis=200 -XX:+UnlockExperimentalVMOptions -XX:+DisableExplicitGC -XX:+AlwaysPreTouch -XX:G1NewSizePercent=30 -XX:G1MaxNewSizePercent=40 -XX:G1HeapRegionSize=8M -XX:G1ReservePercent=20 -XX:G1HeapWastePercent=5 -XX:G1MixedGCCountTarget=4 -XX:InitiatingHeapOccupancyPercent=15 -XX:G1MixedGCLiveThresholdPercent=90 -XX:G1RSetUpdatingPauseTimePercent=5 -XX:SurvivorRatio=32 -XX:+PerfDisableSharedMem -XX:MaxTenuringThreshold=1 -Dusing.aikars.flags=https://mcflags.emc.gs -Daikars.new.flags=true"

WORKDIR /data

RUN apk add --upgrade --no-cache openssl && \
    addgroup -S paper && \
    adduser -S paper -G paper && \
    chown paper:paper /data

USER paper

VOLUME /data

EXPOSE 25565

COPY --chown=paper --from=builder /data/ /data/

COPY --chown=paper --from=builder /opt/papermc/paper-*.jar /opt/papermc/paper.jar

ENTRYPOINT java -Xms$JAVA_MEMORY -Xmx$JAVA_MEMORY $JAVA_FLAGS -jar /opt/papermc/paper.jar --nogui