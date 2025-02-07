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

COPY ./css /usr/share/nginx/html/css
COPY ./fonts /usr/share/nginx/html/fonts
COPY src/assets/images /usr/share/nginx/html/images
COPY ./includes /usr/share/nginx/html/includes
COPY ./js /usr/share/nginx/html/js
COPY ./pages /usr/share/nginx/html/pages
COPY src/assets/svg /usr/share/nginx/html/svg
COPY ./manifest.json /usr/share/nginx/html/manifest.json

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]
