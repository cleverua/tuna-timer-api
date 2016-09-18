CREATE TABLE team_users (
  id serial,
  name varchar(128),
  team_id bigint references teams(id),
  slack_user_id varchar(32),
  CONSTRAINT team_users_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_team_users_slack_user_id ON team_users(slack_user_id);