#################################
# @author K10s Open Source Team #
# @since 01/01/2025             #
#################################
FROM openjdk:17-jdk-slim

WORKDIR /app

COPY app/k10s.jar /app

CMD ["java", "-jar", "k10s.jar"]