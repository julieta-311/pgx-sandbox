create extension if not exists "uuid-ossp";

create table thing (
	id uuid primary key default uuid_generate_v4(),
	name text not null default '',
	labels text[],
	n int not null,
	x numeric,
	created_at timestamp with time zone default (now() at time zone 'utc'),
	stuff jsonb
);
