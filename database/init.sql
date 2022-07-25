CREATE TABLE IF NOT EXISTS users (
  user_id SERIAL PRIMARY KEY,
  username varchar(50) NOT NULL UNIQUE,
  password_hash varchar(250) NOT NULL,
  roles text ARRAY
);

CREATE INDEX ON users USING HASH(username);

INSERT INTO users (username, password_hash, roles) VALUES ('John Mclane', 'asdsadkeltllsas', '{"admin","user"}'), ('Jack Ryan', 'esjsjgjkglslsl', '{"user"}');