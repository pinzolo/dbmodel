-- arranged from: https://github.com/lorint/AdventureWorks-for-Postgres

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
    memo TEXT,
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

CREATE SCHEMA production;

CREATE TABLE production.location(
  location_id SERIAL NOT NULL, -- smallint
  name "Name" NOT NULL,
  cost_rate numeric NOT NULL CONSTRAINT "df_location_cost_rate" DEFAULT (0.00), -- smallmoney -- money
  availability decimal(8, 2) NOT NULL CONSTRAINT "df_location_availability" DEFAULT (0.00),
  modified_date TIMESTAMP NOT NULL CONSTRAINT "df_location_modified_date" DEFAULT (NOW()),
  CONSTRAINT "ck_location_cost_rate" CHECK (cost_rate >= 0.00),
  CONSTRAINT "ck_location_availability" CHECK (availability >= 0.00)
);

COMMENT ON TABLE production.location IS 'Product inventory and manufacturing locations.';
COMMENT ON COLUMN production.location.location_id IS 'Primary key for Location records.';
COMMENT ON COLUMN production.location.cost_rate IS 'Standard hourly cost of the manufacturing location.';
COMMENT ON COLUMN production.location.availability IS 'Work capacity (in hours) of the manufacturing location.';
ALTER TABLE production.location ADD
    CONSTRAINT "pk_location_location_id" PRIMARY KEY
    (location_id);
CLUSTER production.location USING "pk_location_location_id";

