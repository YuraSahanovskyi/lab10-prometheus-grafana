CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

INSERT INTO students (name)
VALUES('Андрій'), ('Богдан'), ('Юрій');