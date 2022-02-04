FROM golang:1.17

ENV USER=postgres
ENV PASSWORD=password
ENV HOST=database
ENV PORT=5432
ENV DBNAME=ip2location