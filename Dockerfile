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

ENV JAVA_OPTIONS="-XX:+UseStringDeduplication -XX:+AlwaysPreTouch"

WORKDIR /data

RUN addgroup -S paper && \
    adduser -S paper -G paper && \
    chown paper:paper /data

USER paper

VOLUME /data

EXPOSE 25565

COPY --chown=paper --from=builder /data/ /data/

COPY --chown=paper --from=builder /opt/papermc/paper-*.jar /opt/papermc/paper.jar

ENTRYPOINT java $JAVA_OPTIONS -jar /opt/papermc/paper.jar --nogui