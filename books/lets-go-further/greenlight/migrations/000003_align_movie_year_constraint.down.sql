ALTER TABLE movies DROP CONSTRAINT movies_year_check;
ALTER TABLE movies ADD CONSTRAINT movies_year_check CHECK (year BETWEEN 1888 AND date_part('year', now()));
