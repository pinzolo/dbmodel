-- Main schema
CREATE SCHEMA schm;
CREATE SCHEMA other;

SET search_path TO schm;

-- Custom domains
CREATE DOMAIN domain1 varchar(10) CHECK (VALUE IN ('A', 'B', 'C'));

-- Tables
CREATE TABLE tbl1 (
    id serial NOT NULL PRIMARY KEY
  , no_size numeric
  , with_length varchar(50)
  , with_precision numeric(8)
  , with_scale numeric(8, 2)
  , not_null integer NOT NULL DEFAULT 1
  , ts_col timestamp NOT NULL DEFAULT current_timestamp
);

CREATE TABLE tbl2 (
    id serial NOT NULL PRIMARY KEY
  , tbl1_id integer
  , idx_key text
  , chk integer
  , range int4range
  , CONSTRAINT tbl2_fk1 FOREIGN KEY (tbl1_id) REFERENCES tbl1(id)
  , CONSTRAINT tbl2_chk1 CHECK (chk > 0)
  , CONSTRAINT tbl2_ex1 EXCLUDE USING gist(range WITH &&)
);
CREATE INDEX tbl2_idx1 ON tbl2(idx_key);

COMMENT ON TABLE tbl2 IS 'This is table2';
COMMENT ON COLUMN tbl2.id IS 'ID';
COMMENT ON COLUMN tbl2.tbl1_id IS 'tbl1ID';

CREATE TABLE tbl3 (
    tbl1_id integer
  , tbl2_id integer
  , CONSTRAINT tbl3_fk1 FOREIGN KEY (tbl1_id) REFERENCES tbl1(id)
  , CONSTRAINT tbl3_fk2 FOREIGN KEY (tbl2_id) REFERENCES tbl2(id)
  , CONSTRAINT tbl3_pk PRIMARY KEY (tbl1_id, tbl2_id)
);
CREATE INDEX tbl3_idx1 ON tbl3(tbl2_id);

CREATE TABLE other.tbl_other (
    col1 integer
  , col2 domain1
  , tbl1_id integer
  , tbl2_id integer
  , CONSTRAINT tbl_other_fk1 FOREIGN KEY (tbl1_id) REFERENCES schm.tbl1(id)
  , CONSTRAINT tbl_other_fk2 FOREIGN KEY (tbl2_id) REFERENCES schm.tbl2(id)
  , CONSTRAINT tbl_other_fk3 FOREIGN KEY (tbl1_id, tbl2_id) REFERENCES schm.tbl3(tbl1_id, tbl2_id)
);
