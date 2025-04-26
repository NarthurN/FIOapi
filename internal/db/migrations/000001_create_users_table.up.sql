CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    patronymic VARCHAR(100),
    age INTEGER NOT NULL,
    sex VARCHAR(20) NOT NULL,
    nationality VARCHAR(100) NOT NULL
);

-- Индексы для отдельных полей
CREATE INDEX idx_users_age ON users(age);
CREATE INDEX idx_users_sex ON users(sex);
CREATE INDEX idx_users_nationality ON users(nationality);

-- Составной индекс для ФИО
CREATE INDEX idx_users_full_name ON users(name, surname, patronymic);