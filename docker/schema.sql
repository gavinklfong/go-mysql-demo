USE forex;

-- DROP TABLE demo.people IF EXISTS;

CREATE TABLE customer  (
    id BIGINT NOT NULL AUTO_INCREMENT,
    name VARCHAR(50),
    tier INT,
    PRIMARY KEY (id)
);

CREATE TABLE forex_rate_booking (
    id BIGINT NOT NULL AUTO_INCREMENT,
    timestamp timestamp,
    base_currency VARCHAR(3),
    counter_currency VARCHAR(3),
    rate decimal(5,4),
    trade_action VARCHAR(50),
    base_currency_amount decimal(14,4),
    PRIMARY KEY (id)
);