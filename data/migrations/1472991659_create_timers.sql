CREATE TABLE public.timers (
  id serial,
  team_user_id bigint references team_users(id),
  task_id bigint references tasks(id),
  started_at timestamp with time zone,
  finished_at timestamp with time zone,
  minutes integer default 0,
  deleted_at timestamp with time zone,
  CONSTRAINT timers_pkey PRIMARY KEY (id)
)