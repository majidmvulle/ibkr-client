FROM eclipse-temurin:17-jre

WORKDIR /app

# Download official IBKR Client Portal Gateway
RUN apt-get update && \
    apt-get install -y --no-install-recommends curl unzip ca-certificates && \
    curl -L https://download2.interactivebrokers.com/portal/clientportal.gw.zip -o clientportal.gw.zip && \
    unzip clientportal.gw.zip && \
    rm clientportal.gw.zip && \
    apt-get purge -y --auto-remove curl unzip && \
    rm -rf /var/lib/apt/lists/*

# Expose Client Portal API port
EXPOSE 5000

# Run the gateway
CMD ["bin/run.sh", "root/conf.yaml"]
