#################################
# @author K10s Open Source Team #
# @since 01/01/2025             #
#################################
FROM python:3.11-slim

WORKDIR /app

COPY app/ /app/

EXPOSE 80

CMD ["python3", "-m", "http.server", "80", "--directory", "/app", "--bind", "0.0.0.0"]