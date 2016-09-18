CREATE TABLE slack_commands (
  id serial,
  team_id varchar(32),
  team_domain varchar(128),
  channel_id varchar(32),
  channel_name varchar(128),
  user_id varchar(32),
  user_name varchar(128),
  command varchar(32),
  text text,
  response_url varchar(255),
  error text,
  created_at timestamp with time zone,
  CONSTRAINT slack_commands_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_slack_commands_team_id ON slack_commands(team_id);
CREATE INDEX idx_slack_commands_user_id ON slack_commands(user_id);
