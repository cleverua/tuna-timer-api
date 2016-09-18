CREATE TABLE tasks (
  id serial,
  name varchar(128),
  hash varchar(8),
  team_id bigint references teams(id),
  project_id bigint references projects(id),
  total_minutes integer default 0,
  CONSTRAINT tasks_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_tasks_name ON tasks(name);
CREATE INDEX idx_tasks_hash ON tasks(hash);