FROM golang:1.18
ENV APP_HOME /app
ADD . $APP_HOME
WORKDIR $APP_HOME

RUN make dep
RUN make build
COPY . .

EXPOSE 8080
CMD ["warehouse"]