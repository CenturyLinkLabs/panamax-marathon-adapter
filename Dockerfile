FROM scratch
MAINTAINER CenturyLink Labs <clt-labs-futuretech@centurylink.com>
EXPOSE 8001
COPY panamax-marathon-adapter /
ENTRYPOINT ["/panamax-marathon-adapter"]
