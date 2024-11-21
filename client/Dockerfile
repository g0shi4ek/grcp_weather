FROM golang

WORKDIR /app

COPY . .
COPY ../cl ./

EXPOSE 3000
CMD ["go", "run", "start_cl.go"]
