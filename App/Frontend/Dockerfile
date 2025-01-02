#################################
# @author K10s Open Source Team #
# @since 01/01/2025             #
#################################
# Frontend build
# Stage 1:
FROM nginx:alpine

LABEL maintainer="K10s Open Source Team"
LABEL version="1.0"
LABEL description="This is the frontend image for the app."

COPY ./nginx.conf /etc/nginx/nginx.conf
COPY ./default.conf /etc/nginx/conf.d/default.conf

COPY ./css /usr/share/nginx/html/css
COPY ./images /usr/share/nginx/html/images
COPY ./includes /usr/share/nginx/html/includes
COPY ./js /usr/share/nginx/html/js
COPY ./pages /usr/share/nginx/html/pages
COPY ./manifest.json /usr/share/nginx/html/manifest.json

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]