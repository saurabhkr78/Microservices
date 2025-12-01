FROM postgres:15.2-alpine3.18
COPY ./up.sql /docker-entrypoint-initdb.d/1.sql
CMD ["postgres"]