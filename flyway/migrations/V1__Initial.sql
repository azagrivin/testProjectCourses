CREATE TABLE codes (
    id    SERIAL PRIMARY KEY,
    value TEXT UNIQUE
);

INSERT INTO codes (value) VALUES('RUB');

CREATE TABLE currencies (
    id      SERIAL PRIMARY KEY,
    code_id INT NOT NULL REFERENCES codes(id),
    date    DATE,
    value   REAL,

    UNIQUE (code_id, date)
);

CREATE TABLE btc_usdt (
    id            SERIAL PRIMARY KEY,
    timestamp     TIMESTAMP UNIQUE,
    buy           REAL,
    sell          REAL,
    high          REAL,
    low           REAL,
    last          REAL,
    average_price REAL
);

CREATE TABLE btc (
    id            SERIAL PRIMARY KEY,
    code_id       INT NOT NULL REFERENCES codes(id),
    timestamp     TIMESTAMP,
    buy           REAL,
    sell          REAL,
    high          REAL,
    low           REAL,
    last          REAL,
    average_price REAL,

    UNIQUE (code_id, timestamp)
);