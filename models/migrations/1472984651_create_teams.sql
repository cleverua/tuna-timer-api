CREATE TABLE teams (
  id serial,
  slack_team_id varchar(32),
  created_at timestamp with time zone,
  CONSTRAINT teams_pkey PRIMARY KEY (id)
);
