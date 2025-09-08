-- minimal schema skeleton (full schema will be added later)
create extension if not exists pgcrypto;

create table if not exists users (
    id uuid primary key default gen_random_uuid(),
    email text unique not null,
    password_hash text not null,
    created_at timestamptz not null default now()
);


