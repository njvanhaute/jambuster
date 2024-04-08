ALTER TABLE tunes ADD CONSTRAINT styles_length_check CHECK (array_length(styles, 1) BETWEEN 1 AND 5);
ALTER TABLE tunes ADD CONSTRAINT keys_length_check CHECK (array_length(keys, 1) BETWEEN 1 AND 10);