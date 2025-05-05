FROM golang

WORKDIR /app

COPY . .
COPY ../sv ./

EXPOSE 5000
CMD ["go", "run", "start_serv.go"]
