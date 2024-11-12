FROM rabbitmq:4.0.3-management-alpine

RUN rabbitmq-plugins enable rabbitmq_tracing

EXPOSE 5672 15672s