FROM eclipse-temurin:17-jre-alpine

WORKDIR /app

# Download official IBKR Client Portal Gateway
RUN apk add --no-cache curl unzip && \
    curl -L https://download2.interactivebrokers.com/portal/clientportal.gw.zip -o clientportal.gw.zip && \
    unzip clientportal.gw.zip && \
    rm clientportal.gw.zip && \
    apk del curl unzip

# Expose Client Portal API port
EXPOSE 5000

# Run the gateway
CMD ["bin/run.sh", "root/conf.yaml"]
