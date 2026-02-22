create table if not exists users (
  id serial primary key,
  name varchar(255) not null,
  email varchar(255) not null unique,
  age int not null default 0 check (age >= 0),
  created_at timestamptz not null default now()
);

insert into users (name, email, age)
values ('Aidos', 'aidos@example.com', 20)
on conflict do nothing;
