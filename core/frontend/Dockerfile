# This is a container to build frontend assets with node.js

FROM node:6.9

MAINTAINER AOE's Frontend-Heroes <team-frontend@aoe.com>
LABEL com.aoe.version="1.0.0"

# Global install yarn package manager
RUN npm set progress=false && \
    npm install -g --progress=false yarn

# Create and ensure folder permissions for yarn
RUN mkdir -p /.config && chmod -R ugo+rw /.config && mkdir -p /.cache && chmod -R ugo+rw /.cache

EXPOSE 1337

ENTRYPOINT ["yarn"]
