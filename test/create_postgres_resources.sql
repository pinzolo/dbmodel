-- pick from: https://github.com/lorint/AdventureWorks-for-Postgres

CREATE DOMAIN "Name" varchar(50) NULL;

CREATE SCHEMA person;

CREATE TABLE person.country_region(
    country_region_code varchar(3) NOT NULL,
    name "Name" NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_country_region_modified_date" DEFAULT (NOW())
);

COMMENT ON TABLE person.country_region IS 'Lookup table containing the ISO standard codes for countries and regions.';
COMMENT ON COLUMN person.country_region.country_region_code IS 'ISO standard code for countries and regions.';
COMMENT ON COLUMN person.country_region.name IS 'Country or region name.';

ALTER TABLE person.country_region ADD
    CONSTRAINT "pk_country_region_country_region_code" PRIMARY KEY
    (country_region_code);
CLUSTER person.country_region USING "pk_country_region_country_region_code";

CREATE SCHEMA sales;

CREATE TABLE sales.currency(
    currency_code char(3) NOT NULL,
    name "Name" NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_currency_modified_date" DEFAULT (NOW())
);
COMMENT ON TABLE sales.currency IS 'Lookup table containing standard ISO currencies.';
COMMENT ON COLUMN sales.currency.currency_code IS 'The ISO code for the currency.';
ALTER TABLE sales.currency ADD
    CONSTRAINT "pk_currency_currency_code" PRIMARY KEY
    (currency_code);
CLUSTER sales.currency USING "pk_currency_currency_code";

CREATE TABLE sales.country_region_currency(
    country_region_code varchar(3) NOT NULL,
    currency_code char(3) NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_country_region_currency_modified_date" DEFAULT (NOW())
);

COMMENT ON COLUMN sales.country_region_currency.country_region_code IS 'ISO code for countries and regions. Foreign key to country_region.country_region_code.';
COMMENT ON COLUMN sales.country_region_currency.currency_code IS 'ISO standard currency code. Foreign key to currency.currency_code.';
ALTER TABLE sales.country_region_currency ADD
    CONSTRAINT "pk_country_region_currency_country_region_code_currency_code" PRIMARY KEY
    (country_region_code, currency_code);
CLUSTER sales.country_region_currency USING "pk_country_region_currency_country_region_code_currency_code";
ALTER TABLE sales.country_region_currency ADD
    CONSTRAINT "fk_country_region_currency_country_region_country_region_code" FOREIGN KEY
    (country_region_code) REFERENCES person.country_region(country_region_code);
ALTER TABLE sales.country_region_currency ADD
    CONSTRAINT "fk_country_region_currency_currency_currency_code" FOREIGN KEY
    (currency_code) REFERENCES sales.currency(currency_code);
