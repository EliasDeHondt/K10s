#################################
# @author K10s Open Source Team #
# @since 01/01/2025             #
#################################
# Frontend build
FROM nginx:alpine

LABEL maintainer "K10s Open Source Team"
LABEL version "1.0"
LABEL description "This is the frontend image for the app."
LABEL org.opencontainers.image.description "This is the frontend image for the app."

COPY ./nginx.conf /etc/nginx/nginx.conf
COPY ./default.conf /etc/nginx/conf.d/default.conf

COPY ./dist/frontend /usr/share/nginx/html

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]
