FROM debian

ADD ./bin/orderservice /app/bin/orderservice
RUN chmod +x /app/bin/orderservice
WORKDIR /app

CMD ["/app/bin/orderservice"]