-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table posts (
  id bigserial primary key,
  username varchar(255) not null,
  title varchar(100) not null,
  content text not null,
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop table if exists posts;
