FROM migrate/migrate:v4.17.0

COPY ./migrations /migrations

ENTRYPOINT ["migrate"]
CMD ["-path", "/migrations", "-database", "postgres://user:password@postgres:5432/question_db?sslmode=disable", "up"]
