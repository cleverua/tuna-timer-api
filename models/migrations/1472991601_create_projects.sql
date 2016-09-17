CREATE TABLE projects (
  id serial,
  name varchar(64),
  slack_channel_id varchar(32),
  slack_channel_name varchar(64),
  team_id bigint references teams(id),
  CONSTRAINT projects_pkey PRIMARY KEY (id)
);