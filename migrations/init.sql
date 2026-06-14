CREATE TABLE IF NOT EXISTS subscriptions
(
    id serial NOT NULL PRIMARY KEY,
    name text NOT NULL,
    price integer NOT NULL,
    user_id uuid NOT NULL,
    start_date date NOT NULL,
    end_date date,
    CONSTRAINT valid_end_date CHECK (end_date IS NULL OR end_date >= start_date)
)