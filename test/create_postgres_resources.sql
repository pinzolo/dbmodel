
-- Support to auto-generate UUIDs (aka GUIDs)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-------------------------------------
-- Custom data types
-------------------------------------

CREATE DOMAIN "OrderNumber" varchar(25) NULL;
CREATE DOMAIN "AccountNumber" varchar(15) NULL;

CREATE DOMAIN "Flag" boolean NOT NULL;
CREATE DOMAIN "NameStyle" boolean NOT NULL;
CREATE DOMAIN "Name" varchar(50) NULL;
CREATE DOMAIN "Phone" varchar(25) NULL;


-------------------------------------
-- Five schemas, with tables and data
-------------------------------------

CREATE SCHEMA person
  CREATE TABLE business_entity(
    business_entity_id SERIAL, --  NOT FOR REPLICATION
    rowguid uuid NOT NULL CONSTRAINT "df_business_entity_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_business_entity_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE person(
    business_entity_id INT NOT NULL,
    person_type char(2) NOT NULL,
    name_style "NameStyle" NOT NULL CONSTRAINT "df_person_name_style" DEFAULT (false),
    title varchar(8) NULL,
    first_name "Name" NOT NULL,
    middle_name "Name" NULL,
    last_name "Name" NOT NULL,
    suffix varchar(10) NULL,
    email_promotion INT NOT NULL CONSTRAINT "df_person_email_promotion" DEFAULT (0),
    additional_contact_info TEXT NULL, -- XML("additional_contact_info_schema_collection"),
    demographics TEXT NULL, -- XML("individual_survey_schema_collection"),
    rowguid uuid NOT NULL CONSTRAINT "df_person_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_person_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_person_email_promotion" CHECK (email_promotion BETWEEN 0 AND 2),
    CONSTRAINT "ck_person_person_type" CHECK (person_type IS NULL OR UPPER(person_type) IN ('SC', 'VC', 'IN', 'EM', 'SP', 'GC'))
  )
  CREATE TABLE state_province(
    state_province_id SERIAL,
    state_province_code char(3) NOT NULL,
    country_region_code varchar(3) NOT NULL,
    is_only_state_province_flag "Flag" NOT NULL CONSTRAINT "df_state_province_is_only_state_province_flag" DEFAULT (true),
    name "Name" NOT NULL,
    territory_id INT NOT NULL,
    rowguid uuid NOT NULL CONSTRAINT "df_state_province_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_state_province_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE address(
    address_id SERIAL, --  NOT FOR REPLICATION
    address_line1 varchar(60) NOT NULL,
    address_line2 varchar(60) NULL,
    city varchar(30) NOT NULL,
    state_province_id INT NOT NULL,
    postal_code varchar(15) NOT NULL,
    spatial_location varchar(44) NULL,
    rowguid uuid NOT NULL CONSTRAINT "df_address_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_address_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE address_type(
    address_type_id SERIAL,
    name "Name" NOT NULL,
    rowguid uuid NOT NULL CONSTRAINT "df_address_type_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_address_type_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE business_entity_address(
    business_entity_id INT NOT NULL,
    address_id INT NOT NULL,
    address_type_id INT NOT NULL,
    rowguid uuid NOT NULL CONSTRAINT "df_business_entity_address_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_business_entity_address_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE contact_type(
    contact_type_id SERIAL,
    name "Name" NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_contact_type_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE business_entity_contact(
    business_entity_id INT NOT NULL,
    person_id INT NOT NULL,
    contact_type_id INT NOT NULL,
    rowguid uuid NOT NULL CONSTRAINT "df_business_entity_contact_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_business_entity_contact_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE email_address(
    business_entity_id INT NOT NULL,
    email_address_id SERIAL,
    email_address varchar(50) NULL,
    rowguid uuid NOT NULL CONSTRAINT "df_email_address_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_email_address_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE password(
    business_entity_id INT NOT NULL,
    password_hash VARCHAR(128) NOT NULL,
    password_salt VARCHAR(10) NOT NULL,
    rowguid uuid NOT NULL CONSTRAINT "df_password_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_password_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE phone_number_type(
    phone_number_type_id SERIAL,
    name "Name" NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_phone_number_type_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE person_phone(
    business_entity_id INT NOT NULL,
    phone_number "Phone" NOT NULL,
    phone_number_type_id INT NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_person_phone_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE country_region(
    country_region_code varchar(3) NOT NULL,
    name "Name" NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_country_region_modified_date" DEFAULT (NOW())
  );

COMMENT ON SCHEMA person IS 'Contains objects related to names and addresses of customers, vendors, and employees';

CREATE SCHEMA human_resources
  CREATE TABLE department(
    department_id SERIAL NOT NULL, -- smallint
    name "Name" NOT NULL,
    group_name "Name" NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_department_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE employee(
    business_entity_id INT NOT NULL,
    national_id_number varchar(15) NOT NULL,
    login_id varchar(256) NOT NULL,
    org varchar NULL,-- hierarchyid, will become organization_node
    organization_level INT NULL, -- AS organization_node.GetLevel(),
    job_title varchar(50) NOT NULL,
    birth_date DATE NOT NULL,
    marital_status char(1) NOT NULL,
    gender char(1) NOT NULL,
    hire_date DATE NOT NULL,
    salaried_flag "Flag" NOT NULL CONSTRAINT "df_employee_salaried_flag" DEFAULT (true),
    vacation_hours smallint NOT NULL CONSTRAINT "df_employee_vacation_hours" DEFAULT (0),
    sick_leave_hours smallint NOT NULL CONSTRAINT "df_employee_sick_leave_hours" DEFAULT (0),
    current_flag "Flag" NOT NULL CONSTRAINT "df_employee_current_flag" DEFAULT (true),
    rowguid uuid NOT NULL CONSTRAINT "df_employee_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_employee_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_employee_birth_date" CHECK (birth_date BETWEEN '1930-01-01' AND NOW() - INTERVAL '18 years'),
    CONSTRAINT "ck_employee_marital_status" CHECK (UPPER(marital_status) IN ('M', 'S')), -- Married or Single
    CONSTRAINT "ck_employee_hire_date" CHECK (hire_date BETWEEN '1996-07-01' AND NOW() + INTERVAL '1 day'),
    CONSTRAINT "ck_employee_gender" CHECK (UPPER(gender) IN ('M', 'F')), -- Male or Female
    CONSTRAINT "ck_employee_vacation_hours" CHECK (vacation_hours BETWEEN -40 AND 240),
    CONSTRAINT "ck_employee_sick_leave_hours" CHECK (sick_leave_hours BETWEEN 0 AND 120)
  )
  CREATE TABLE employee_department_history(
    business_entity_id INT NOT NULL,
    department_id smallint NOT NULL,
    shift_id smallint NOT NULL, -- tinyint
    start_date DATE NOT NULL,
    end_date DATE NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_employee_department_history_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_employee_department_history_end_date" CHECK ((end_date >= start_date) OR (end_date IS NULL))
  )
  CREATE TABLE employee_pay_history(
    business_entity_id INT NOT NULL,
    rate_change_date TIMESTAMP NOT NULL,
    rate numeric NOT NULL, -- money
    pay_frequency smallint NOT NULL,  -- tinyint
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_employee_pay_history_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_employee_pay_history_pay_frequency" CHECK (pay_frequency IN (1, 2)), -- 1 = monthly salary, 2 = biweekly salary
    CONSTRAINT "ck_employee_pay_history_rate" CHECK (rate BETWEEN 6.50 AND 200.00)
  )
  CREATE TABLE job_candidate(
    job_candidate_id SERIAL NOT NULL, -- int
    business_entity_id INT NULL,
    resume varchar NULL, -- XML(hr_resume_schema_collection)
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_job_candidate_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE shift(
    shift_id SERIAL NOT NULL, -- tinyint
    name "Name" NOT NULL,
    start_time time NOT NULL,
    end_time time NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_shift_modified_date" DEFAULT (NOW())
  );

COMMENT ON SCHEMA human_resources IS 'Contains objects related to employees and departments.';

-- Calculated column that needed to be there just for the CSV import
ALTER TABLE human_resources.employee DROP COLUMN organization_level;

-- Employee HierarchyID column
ALTER TABLE human_resources.employee ADD organization_node VARCHAR DEFAULT '/';
-- Convert from all the hex to a stream of hierarchyid bits
WITH RECURSIVE hier AS (
  SELECT business_entity_id, org, get_byte(decode(substring(org, 1, 2), 'hex'), 0)::bit(8)::varchar AS bits, 2 AS i
    FROM human_resources.employee
  UNION ALL
  SELECT e.business_entity_id, e.org, hier.bits || get_byte(decode(substring(e.org, i + 1, 2), 'hex'), 0)::bit(8)::varchar, i + 2 AS i
    FROM human_resources.employee AS e INNER JOIN
      hier ON e.business_entity_id = hier.business_entity_id AND i < LENGTH(e.org)
)
UPDATE human_resources.employee AS emp
  SET org = COALESCE(trim(trailing '0' FROM hier.bits::TEXT), '')
  FROM hier
  WHERE emp.business_entity_id = hier.business_entity_id
    AND (hier.org IS NULL OR i = LENGTH(hier.org));

-- Convert bits to the real hieararchy paths
CREATE OR REPLACE FUNCTION f_convert_org_nodes()
  RETURNS void AS
$func$
DECLARE
  got_none BOOLEAN;
BEGIN
  LOOP
  got_none := true;
  -- 01 = 0-3
  UPDATE human_resources.employee
   SET organization_node = organization_node || SUBSTRING(org, 3,2)::bit(2)::INTEGER::VARCHAR || CASE SUBSTRING(org, 5, 1) WHEN '0' THEN '.' ELSE '/' END,
     org = SUBSTRING(org, 6, 9999)
    WHERE org LIKE '01%';
  IF FOUND THEN
    got_none := false;
  END IF;

  -- 100 = 4-7
  UPDATE human_resources.employee
   SET organization_node = organization_node || (SUBSTRING(org, 4,2)::bit(2)::INTEGER + 4)::VARCHAR || CASE SUBSTRING(org, 6, 1) WHEN '0' THEN '.' ELSE '/' END,
     org = SUBSTRING(org, 7, 9999)
    WHERE org LIKE '100%';
  IF FOUND THEN
    got_none := false;
  END IF;

  -- 101 = 8-15
  UPDATE human_resources.employee
   SET organization_node = organization_node || (SUBSTRING(org, 4,3)::bit(3)::INTEGER + 8)::VARCHAR || CASE SUBSTRING(org, 7, 1) WHEN '0' THEN '.' ELSE '/' END,
     org = SUBSTRING(org, 8, 9999)
    WHERE org LIKE '101%';
  IF FOUND THEN
    got_none := false;
  END IF;

  -- 110 = 16-79
  UPDATE human_resources.employee
   SET organization_node = organization_node || ((SUBSTRING(org, 4,2)||SUBSTRING(org, 7,1)||SUBSTRING(org, 9,3))::bit(6)::INTEGER + 16)::VARCHAR || CASE SUBSTRING(org, 12, 1) WHEN '0' THEN '.' ELSE '/' END,
     org = SUBSTRING(org, 13, 9999)
    WHERE org LIKE '110%';
  IF FOUND THEN
    got_none := false;
  END IF;

  -- 1110 = 80-1103
  UPDATE human_resources.employee
   SET organization_node = organization_node || ((SUBSTRING(org, 5,3)||SUBSTRING(org, 9,3)||SUBSTRING(org, 13,1)||SUBSTRING(org, 15,3))::bit(10)::INTEGER + 80)::VARCHAR || CASE SUBSTRING(org, 18, 1) WHEN '0' THEN '.' ELSE '/' END,
     org = SUBSTRING(org, 19, 9999)
    WHERE org LIKE '1110%';
  IF FOUND THEN
    got_none := false;
  END IF;
  EXIT WHEN got_none;
  END LOOP;
END
$func$ LANGUAGE plpgsql;

SELECT f_convert_org_nodes();
-- Drop the original binary hierarchyid column
ALTER TABLE human_resources.employee DROP COLUMN org;
DROP FUNCTION f_convert_org_nodes();

CREATE SCHEMA production
  CREATE TABLE bill_of_materials(
    bill_of_materials_id SERIAL NOT NULL, -- int
    product_assembly_id INT NULL,
    component_id INT NOT NULL,
    start_date TIMESTAMP NOT NULL CONSTRAINT "df_bill_of_materials_start_date" DEFAULT (NOW()),
    end_date TIMESTAMP NULL,
    unit_measure_code char(3) NOT NULL,
    bom_level smallint NOT NULL,
    per_assembly_qty decimal(8, 2) NOT NULL CONSTRAINT "df_bill_of_materials_per_assembly_qty" DEFAULT (1.00),
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_bill_of_materials_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_bill_of_materials_end_date" CHECK ((end_date > start_date) OR (end_date IS NULL)),
    CONSTRAINT "ck_bill_of_materials_product_assembly_id" CHECK (product_assembly_id <> component_id),
    CONSTRAINT "ck_bill_of_materials_bom_level" CHECK (((product_assembly_id IS NULL)
        AND (bom_level = 0) AND (per_assembly_qty = 1.00))
        OR ((product_assembly_id IS NOT NULL) AND (bom_level >= 1))),
    CONSTRAINT "ck_bill_of_materials_per_assembly_qty" CHECK (per_assembly_qty >= 1.00)
  )
  CREATE TABLE culture(
    culture_id char(6) NOT NULL,
    name "Name" NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_culture_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE document(
    doc varchar NULL,-- hierarchyid, will become document_node
    document_level INTEGER, -- AS document_node.get_level(),
    title varchar(50) NOT NULL,
    owner INT NOT NULL,
    folder_flag "Flag" NOT NULL CONSTRAINT "df_document_folder_flag" DEFAULT (false),
    file_name varchar(400) NOT NULL,
    file_extension varchar(8) NULL,
    revision char(5) NOT NULL,
    change_number INT NOT NULL CONSTRAINT "df_document_change_number" DEFAULT (0),
    status smallint NOT NULL, -- tinyint
    document_summary text NULL,
    document bytea  NULL, -- varbinary
    rowguid uuid NOT NULL UNIQUE CONSTRAINT "df_document_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_document_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_document_status" CHECK (Status BETWEEN 1 AND 3)
  )
  CREATE TABLE product_category(
    product_category_id SERIAL NOT NULL, -- int
    name "Name" NOT NULL,
    rowguid uuid NOT NULL CONSTRAINT "df_product_category_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_product_category_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE product_subcategory(
    product_subcategory_id SERIAL NOT NULL, -- int
    product_category_id INT NOT NULL,
    name "Name" NOT NULL,
    rowguid uuid NOT NULL CONSTRAINT "df_product_subcategory_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_product_subcategory_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE product_model(
    product_model_id SERIAL NOT NULL, -- int
    name "Name" NOT NULL,
    catalog_description varchar NULL, -- XML(production.product_description_schema_collection)
    instructions varchar NULL, -- XML(production.manu_instructions_schema_collection)
    rowguid uuid NOT NULL CONSTRAINT "df_product_model_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_product_model_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE product(
    product_id SERIAL NOT NULL, -- int
    name "Name" NOT NULL,
    product_number varchar(25) NOT NULL,
    make_flag "Flag" NOT NULL CONSTRAINT "df_product_make_flag" DEFAULT (true),
    finished_goods_flag "Flag" NOT NULL CONSTRAINT "df_product_finished_goods_flag" DEFAULT (true),
    color varchar(15) NULL,
    safety_stock_level smallint NOT NULL,
    reorder_point smallint NOT NULL,
    standard_cost numeric NOT NULL, -- money
    list_price numeric NOT NULL, -- money
    size varchar(5) NULL,
    size_unit_measure_code char(3) NULL,
    weight_unit_measure_code char(3) NULL,
    weight decimal(8, 2) NULL,
    days_to_manufacture INT NOT NULL,
    product_line char(2) NULL,
    class char(2) NULL,
    style char(2) NULL,
    product_subcategory_id INT NULL,
    product_model_id INT NULL,
    sell_start_date TIMESTAMP NOT NULL,
    sell_end_date TIMESTAMP NULL,
    discontinued_date TIMESTAMP NULL,
    rowguid uuid NOT NULL CONSTRAINT "df_product_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_product_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_product_safety_stock_level" CHECK (safety_stock_level > 0),
    CONSTRAINT "ck_product_reorder_point" CHECK (reorder_point > 0),
    CONSTRAINT "ck_product_standard_cost" CHECK (standard_cost >= 0.00),
    CONSTRAINT "ck_product_list_price" CHECK (list_price >= 0.00),
    CONSTRAINT "ck_product_weight" CHECK (weight > 0.00),
    CONSTRAINT "ck_product_days_to_manufacture" CHECK (days_to_manufacture >= 0),
    CONSTRAINT "ck_product_product_line" CHECK (UPPER(product_line) IN ('S', 'T', 'M', 'R') OR product_line IS NULL),
    CONSTRAINT "ck_product_class" CHECK (UPPER(class) IN ('L', 'M', 'H') OR Class IS NULL),
    CONSTRAINT "ck_product_style" CHECK (UPPER(style) IN ('W', 'M', 'U') OR Style IS NULL),
    CONSTRAINT "ck_product_sell_end_date" CHECK ((sell_end_date >= sell_start_date) OR (sell_end_date IS NULL))
  )
  CREATE TABLE product_cost_history(
    product_id INT NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NULL,
    standard_cost numeric NOT NULL,  -- money
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_product_cost_history_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_product_cost_history_end_date" CHECK ((end_date >= start_date) OR (end_date IS NULL)),
    CONSTRAINT "ck_product_cost_history_standard_cost" CHECK (standard_cost >= 0.00)
  )
  CREATE TABLE product_description(
    product_description_id SERIAL NOT NULL, -- int
    description varchar(400) NOT NULL,
    rowguid uuid NOT NULL CONSTRAINT "df_product_description_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_product_description_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE product_document(
    product_id INT NOT NULL,
    doc varchar NOT NULL, -- hierarchy_id, will become document_node
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_product_document_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE location(
    location_id SERIAL NOT NULL, -- smallint
    name "Name" NOT NULL,
    cost_rate numeric NOT NULL CONSTRAINT "df_location_cost_rate" DEFAULT (0.00), -- smallmoney -- money
    availability decimal(8, 2) NOT NULL CONSTRAINT "df_location_availability" DEFAULT (0.00),
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_location_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_location_cost_rate" CHECK (cost_rate >= 0.00),
    CONSTRAINT "ck_location_availability" CHECK (availability >= 0.00)
  )
  CREATE TABLE product_inventory(
    product_id INT NOT NULL,
    location_id smallint NOT NULL,
    shelf varchar(10) NOT NULL,
    bin smallint NOT NULL, -- tinyint
    quantity smallint NOT NULL CONSTRAINT "df_product_inventory_quantity" DEFAULT (0),
    rowguid uuid NOT NULL CONSTRAINT "df_product_inventory_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_product_inventory_modified_date" DEFAULT (NOW()),
--    CONSTRAINT "ck_product_inventory_shelf" CHECK ((shelf LIKE 'AZa-z]') OR (shelf = 'N/A')),
    CONSTRAINT "ck_product_inventory_bin" CHECK (bin BETWEEN 0 AND 100)
  )
  CREATE TABLE product_list_price_history(
    product_id INT NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NULL,
    list_price numeric NOT NULL,  -- money
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_product_list_price_history_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_product_list_price_history_end_date" CHECK ((end_date >= start_date) OR (end_date IS NULL)),
    CONSTRAINT "ck_product_list_price_history_list_price" CHECK (list_price > 0.00)
  )
  CREATE TABLE illustration(
    illustration_id SERIAL NOT NULL, -- int
    diagram varchar NULL,  -- XML
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_illustration_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE product_model_illustration(
    product_model_id INT NOT NULL,
    illustration_id INT NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_product_model_illustration_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE product_model_product_description_culture(
    product_model_id INT NOT NULL,
    product_description_id INT NOT NULL,
    culture_id char(6) NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_product_model_product_description_culture_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE product_photo(
    product_photo_id SERIAL NOT NULL, -- int
    thumb_nail_photo bytea NULL,-- varbinary
    thumbnail_photo_file_name varchar(50) NULL,
    large_photo bytea NULL,-- varbinary
    large_photo_file_name varchar(50) NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_product_photo_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE product_product_photo(
    product_id INT NOT NULL,
    product_photo_id INT NOT NULL,
    "primary" "Flag" NOT NULL CONSTRAINT "df_product_product_photo_primary" DEFAULT (false),
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_product_product_photo_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE product_review(
    product_review_id SERIAL NOT NULL, -- int
    product_id INT NOT NULL,
    reviewer_name "Name" NOT NULL,
    review_date TIMESTAMP NOT NULL CONSTRAINT "df_product_review_review_date" DEFAULT (NOW()),
    email_address varchar(50) NOT NULL,
    rating INT NOT NULL,
    comments varchar(3850),
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_product_review_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_product_review_rating" CHECK (rating BETWEEN 1 AND 5)
  )
  CREATE TABLE scrap_reason(
    scrap_reason_id SERIAL NOT NULL, -- smallint
    name "Name" NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_scrap_reason_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE transaction_history(
    transaction_id SERIAL NOT NULL, -- INT IDENTITY (100000, 1)
    product_id INT NOT NULL,
    reference_order_id INT NOT NULL,
    reference_order_line_id INT NOT NULL CONSTRAINT "df_transaction_history_reference_order_line_id" DEFAULT (0),
    transaction_date TIMESTAMP NOT NULL CONSTRAINT "df_transaction_history_transaction_date" DEFAULT (NOW()),
    transaction_type char(1) NOT NULL,
    quantity INT NOT NULL,
    actual_cost numeric NOT NULL,  -- money
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_transaction_history_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_transaction_history_transaction_type" CHECK (UPPER(transaction_type) IN ('W', 'S', 'P'))
  )
  CREATE TABLE transaction_history_archive(
    transaction_id INT NOT NULL,
    product_id INT NOT NULL,
    reference_order_id INT NOT NULL,
    reference_order_line_id INT NOT NULL CONSTRAINT "df_transaction_history_archive_reference_order_line_id" DEFAULT (0),
    transaction_date TIMESTAMP NOT NULL CONSTRAINT "df_transaction_history_archive_transaction_date" DEFAULT (NOW()),
    transaction_type char(1) NOT NULL,
    quantity INT NOT NULL,
    actual_cost numeric NOT NULL,  -- money
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_transaction_history_archive_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_transaction_history_archive_transaction_type" CHECK (UPPER(transaction_type) IN ('W', 'S', 'P'))
  )
  CREATE TABLE unit_measure(
    unit_measure_code char(3) NOT NULL,
    name "Name" NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_unit_measure_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE work_order(
    work_order_id SERIAL NOT NULL, -- int
    product_id INT NOT NULL,
    order_qty INT NOT NULL,
    stocked_qty int, -- AS ISNULL(order_qty - scrapped_qty, 0),
    scrapped_qty smallint NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NULL,
    due_date TIMESTAMP NOT NULL,
    scrap_reason_id smallint NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_work_order_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_work_order_order_qty" CHECK (order_qty > 0),
    CONSTRAINT "ck_work_order_scrapped_qty" CHECK (scrapped_qty >= 0),
    CONSTRAINT "ck_work_order_end_date" CHECK ((end_date >= start_date) OR (end_date IS NULL))
  )
  CREATE TABLE work_order_routing(
    work_order_id INT NOT NULL,
    product_id INT NOT NULL,
    operation_sequence smallint NOT NULL,
    location_id smallint NOT NULL,
    scheduled_start_date TIMESTAMP NOT NULL,
    scheduled_end_date TIMESTAMP NOT NULL,
    actual_start_date TIMESTAMP NULL,
    actual_end_date TIMESTAMP NULL,
    actual_resource_hrs decimal(9, 4) NULL,
    planned_cost numeric NOT NULL, -- money
    actual_cost numeric NULL,  -- money
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_work_order_routing_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_work_order_routing_scheduled_end_date" CHECK (scheduled_end_date >= scheduled_start_date),
    CONSTRAINT "ck_work_order_routing_actual_end_date" CHECK ((actual_end_date >= actual_start_date)
        OR (actual_end_date IS NULL) OR (actual_start_date IS NULL)),
    CONSTRAINT "ck_work_order_routing_actual_resource_hrs" CHECK (actual_resource_hrs >= 0.0000),
    CONSTRAINT "ck_work_order_routing_planned_cost" CHECK (planned_cost > 0.00),
    CONSTRAINT "ck_work_order_routing_actual_cost" CHECK (actual_cost > 0.00)
  );

COMMENT ON SCHEMA production IS 'Contains objects related to products, inventory, and manufacturing.';

-- Calculated columns that needed to be there just for the CSV import
ALTER TABLE production.work_order DROP COLUMN stocked_qty;
ALTER TABLE production.document DROP COLUMN document_level;

-- Document HierarchyID column
ALTER TABLE production.document ADD document_node VARCHAR DEFAULT '/';
-- Convert from all the hex to a stream of hierarchyid bits
WITH RECURSIVE hier AS (
  SELECT rowguid, doc, get_byte(decode(substring(doc, 1, 2), 'hex'), 0)::bit(8)::varchar AS bits, 2 AS i
    FROM production.document
  UNION ALL
  SELECT e.rowguid, e.doc, hier.bits || get_byte(decode(substring(e.doc, i + 1, 2), 'hex'), 0)::bit(8)::varchar, i + 2 AS i
    FROM production.document AS e INNER JOIN
      hier ON e.rowguid = hier.rowguid AND i < LENGTH(e.doc)
)
UPDATE production.document AS emp
  SET doc = COALESCE(trim(trailing '0' FROM hier.bits::TEXT), '')
  FROM hier
  WHERE emp.rowguid = hier.rowguid
    AND (hier.doc IS NULL OR i = LENGTH(hier.doc));

-- Convert bits to the real hieararchy paths
CREATE OR REPLACE FUNCTION f_convert_doc_nodes()
  RETURNS void AS
$func$
DECLARE
  got_none BOOLEAN;
BEGIN
  LOOP
  got_none := true;
  -- 01 = 0-3
  UPDATE production.document
   SET document_node = document_node || SUBSTRING(doc, 3,2)::bit(2)::INTEGER::VARCHAR || CASE SUBSTRING(doc, 5, 1) WHEN '0' THEN '.' ELSE '/' END,
     doc = SUBSTRING(doc, 6, 9999)
    WHERE doc LIKE '01%';
  IF FOUND THEN
    got_none := false;
  END IF;

  -- 100 = 4-7
  UPDATE production.document
   SET document_node = document_node || (SUBSTRING(doc, 4,2)::bit(2)::INTEGER + 4)::VARCHAR || CASE SUBSTRING(doc, 6, 1) WHEN '0' THEN '.' ELSE '/' END,
     doc = SUBSTRING(doc, 7, 9999)
    WHERE doc LIKE '100%';
  IF FOUND THEN
    got_none := false;
  END IF;

  -- 101 = 8-15
  UPDATE production.document
   SET document_node = document_node || (SUBSTRING(doc, 4,3)::bit(3)::INTEGER + 8)::VARCHAR || CASE SUBSTRING(doc, 7, 1) WHEN '0' THEN '.' ELSE '/' END,
     doc = SUBSTRING(doc, 8, 9999)
    WHERE doc LIKE '101%';
  IF FOUND THEN
    got_none := false;
  END IF;

  -- 110 = 16-79
  UPDATE production.document
   SET document_node = document_node || ((SUBSTRING(doc, 4,2)||SUBSTRING(doc, 7,1)||SUBSTRING(doc, 9,3))::bit(6)::INTEGER + 16)::VARCHAR || CASE SUBSTRING(doc, 12, 1) WHEN '0' THEN '.' ELSE '/' END,
     doc = SUBSTRING(doc, 13, 9999)
    WHERE doc LIKE '110%';
  IF FOUND THEN
    got_none := false;
  END IF;

  -- 1110 = 80-1103
  UPDATE production.document
   SET document_node = document_node || ((SUBSTRING(doc, 5,3)||SUBSTRING(doc, 9,3)||SUBSTRING(doc, 13,1)||SUBSTRING(doc, 15,3))::bit(10)::INTEGER + 80)::VARCHAR || CASE SUBSTRING(doc, 18, 1) WHEN '0' THEN '.' ELSE '/' END,
     doc = SUBSTRING(doc, 19, 9999)
    WHERE doc LIKE '1110%';
  IF FOUND THEN
    got_none := false;
  END IF;
  EXIT WHEN got_none;
  END LOOP;
END
$func$ LANGUAGE plpgsql;

SELECT f_convert_doc_nodes();
-- Drop the original binary hierarchyid column
ALTER TABLE production.document DROP COLUMN Doc;
DROP FUNCTION f_convert_doc_nodes();

-- ProductDocument HierarchyID column
  ALTER TABLE production.product_document ADD document_node VARCHAR DEFAULT '/';
ALTER TABLE production.product_document ADD rowguid uuid NOT NULL CONSTRAINT "df_product_document_rowguid" DEFAULT (uuid_generate_v1());
-- Convert from all the hex to a stream of hierarchyid bits
WITH RECURSIVE hier AS (
  SELECT rowguid, doc, get_byte(decode(substring(doc, 1, 2), 'hex'), 0)::bit(8)::varchar AS bits, 2 AS i
    FROM production.product_document
  UNION ALL
  SELECT e.rowguid, e.doc, hier.bits || get_byte(decode(substring(e.doc, i + 1, 2), 'hex'), 0)::bit(8)::varchar, i + 2 AS i
    FROM production.product_document AS e INNER JOIN
      hier ON e.rowguid = hier.rowguid AND i < LENGTH(e.doc)
)
UPDATE production.product_document AS emp
  SET doc = COALESCE(trim(trailing '0' FROM hier.bits::TEXT), '')
  FROM hier
  WHERE emp.rowguid = hier.rowguid
    AND (hier.doc IS NULL OR i = LENGTH(hier.doc));

-- Convert bits to the real hieararchy paths
CREATE OR REPLACE FUNCTION f_convert_doc_nodes()
  RETURNS void AS
$func$
DECLARE
  got_none BOOLEAN;
BEGIN
  LOOP
  got_none := true;
  -- 01 = 0-3
  UPDATE production.product_document
   SET document_node = document_node || SUBSTRING(doc, 3,2)::bit(2)::INTEGER::VARCHAR || CASE SUBSTRING(doc, 5, 1) WHEN '0' THEN '.' ELSE '/' END,
     doc = SUBSTRING(doc, 6, 9999)
    WHERE doc LIKE '01%';
  IF FOUND THEN
    got_none := false;
  END IF;

  -- 100 = 4-7
  UPDATE production.product_document
   SET document_node = document_node || (SUBSTRING(doc, 4,2)::bit(2)::INTEGER + 4)::VARCHAR || CASE SUBSTRING(doc, 6, 1) WHEN '0' THEN '.' ELSE '/' END,
     doc = SUBSTRING(doc, 7, 9999)
    WHERE doc LIKE '100%';
  IF FOUND THEN
    got_none := false;
  END IF;

  -- 101 = 8-15
  UPDATE production.product_document
   SET document_node = document_node || (SUBSTRING(doc, 4,3)::bit(3)::INTEGER + 8)::VARCHAR || CASE SUBSTRING(doc, 7, 1) WHEN '0' THEN '.' ELSE '/' END,
     doc = SUBSTRING(doc, 8, 9999)
    WHERE doc LIKE '101%';
  IF FOUND THEN
    got_none := false;
  END IF;

  -- 110 = 16-79
  UPDATE production.product_document
   SET document_node = document_node || ((SUBSTRING(doc, 4,2)||SUBSTRING(doc, 7,1)||SUBSTRING(doc, 9,3))::bit(6)::INTEGER + 16)::VARCHAR || CASE SUBSTRING(doc, 12, 1) WHEN '0' THEN '.' ELSE '/' END,
     doc = SUBSTRING(doc, 13, 9999)
    WHERE doc LIKE '110%';
  IF FOUND THEN
    got_none := false;
  END IF;

  -- 1110 = 80-1103
  UPDATE production.product_document
   SET document_node = document_node || ((SUBSTRING(doc, 5,3)||SUBSTRING(doc, 9,3)||SUBSTRING(doc, 13,1)||SUBSTRING(doc, 15,3))::bit(10)::INTEGER + 80)::VARCHAR || CASE SUBSTRING(doc, 18, 1) WHEN '0' THEN '.' ELSE '/' END,
     doc = SUBSTRING(doc, 19, 9999)
    WHERE doc LIKE '1110%';
  IF FOUND THEN
    got_none := false;
  END IF;
  EXIT WHEN got_none;
  END LOOP;
END
$func$ LANGUAGE plpgsql;

SELECT f_convert_doc_nodes();
-- Drop the original binary hierarchyid column
ALTER TABLE production.product_document DROP COLUMN Doc;
DROP FUNCTION f_convert_doc_nodes();
ALTER TABLE production.product_document DROP COLUMN rowguid;





CREATE SCHEMA purchasing
  CREATE TABLE product_vendor(
    product_id INT NOT NULL,
    business_entity_id INT NOT NULL,
    average_lead_time INT NOT NULL,
    standard_price numeric NOT NULL, -- money
    last_receipt_cost numeric NULL, -- money
    last_receipt_date TIMESTAMP NULL,
    min_order_qty INT NOT NULL,
    max_order_qty INT NOT NULL,
    on_order_qty INT NULL,
    unit_measure_code char(3) NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_product_vendor_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_product_vendor_average_lead_time" CHECK (average_lead_time >= 1),
    CONSTRAINT "ck_product_vendor_standard_price" CHECK (standard_price > 0.00),
    CONSTRAINT "ck_product_vendor_last_receipt_cost" CHECK (last_receipt_cost > 0.00),
    CONSTRAINT "ck_product_vendor_min_order_qty" CHECK (min_order_qty >= 1),
    CONSTRAINT "ck_product_vendor_max_order_qty" CHECK (max_order_qty >= 1),
    CONSTRAINT "ck_product_vendor_on_order_qty" CHECK (on_order_qty >= 0)
  )
  CREATE TABLE purchase_order_detail(
    purchase_order_id INT NOT NULL,
    purchase_order_detail_id SERIAL NOT NULL, -- int
    due_date TIMESTAMP NOT NULL,
    order_qty smallint NOT NULL,
    product_id INT NOT NULL,
    unit_price numeric NOT NULL, -- money
    line_total numeric, -- AS ISNULL(OrderQty * UnitPrice, 0.00),
    received_qty decimal(8, 2) NOT NULL,
    rejected_qty decimal(8, 2) NOT NULL,
    stocked_qty numeric, -- AS ISNULL(ReceivedQty - RejectedQty, 0.00),
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_purchase_order_detail_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_purchase_order_detail_order_qty" CHECK (order_qty > 0),
    CONSTRAINT "ck_purchase_order_detail_unit_price" CHECK (unit_price >= 0.00),
    CONSTRAINT "ck_purchase_order_detail_received_qty" CHECK (received_qty >= 0.00),
    CONSTRAINT "ck_purchase_order_detail_rejected_qty" CHECK (rejected_qty >= 0.00)
  )
  CREATE TABLE purchase_order_header(
    purchase_order_id SERIAL NOT NULL,  -- int
    revision_number smallint NOT NULL CONSTRAINT "df_purchase_order_header_revision_number" DEFAULT (0),  -- tinyint
    status smallint NOT NULL CONSTRAINT "df_purchase_order_header_status" DEFAULT (1),  -- tinyint
    employee_id INT NOT NULL,
    vendor_id INT NOT NULL,
    ship_method_id INT NOT NULL,
    order_date TIMESTAMP NOT NULL CONSTRAINT "df_purchase_order_header_order_date" DEFAULT (NOW()),
    ship_date TIMESTAMP NULL,
    sub_total numeric NOT NULL CONSTRAINT "df_purchase_order_header_sub_total" DEFAULT (0.00),  -- money
    tax_amt numeric NOT NULL CONSTRAINT "df_purchase_order_header_tax_amt" DEFAULT (0.00),  -- money
    freight numeric NOT NULL CONSTRAINT "df_purchase_order_header_freight" DEFAULT (0.00),  -- money
    total_due numeric, -- AS ISNULL(sub_total + tax_amt + freight, 0) PERSISTED NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_purchase_order_header_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_purchase_order_header_status" CHECK (status BETWEEN 1 AND 4), -- 1 = Pending; 2 = Approved; 3 = Rejected; 4 = Complete
    CONSTRAINT "ck_purchase_order_header_ship_date" CHECK ((ship_date >= order_date) OR (ship_date IS NULL)),
    CONSTRAINT "ck_purchase_order_header_sub_total" CHECK (sub_total >= 0.00),
    CONSTRAINT "ck_purchase_order_header_tax_amt" CHECK (tax_amt >= 0.00),
    CONSTRAINT "ck_purchase_order_header_freight" CHECK (freight >= 0.00)
  )
  CREATE TABLE ship_method(
    ship_method_id SERIAL NOT NULL, -- int
    name "Name" NOT NULL,
    ship_base numeric NOT NULL CONSTRAINT "df_ship_method_ship_base" DEFAULT (0.00), -- money
    ship_rate numeric NOT NULL CONSTRAINT "df_ship_method_ship_rate" DEFAULT (0.00), -- money
    rowguid uuid NOT NULL CONSTRAINT "df_ship_method_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_ship_method_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_ship_method_ship_base" CHECK (ship_base > 0.00),
    CONSTRAINT "ck_ship_method_ship_rate" CHECK (ship_rate > 0.00)
  )
  CREATE TABLE vendor(
    business_entity_id INT NOT NULL,
    account_number "AccountNumber" NOT NULL,
    name "Name" NOT NULL,
    credit_rating smallint NOT NULL, -- tinyint
    preferred_vendor_status "Flag" NOT NULL CONSTRAINT "df_vendor_preferred_vendor_status" DEFAULT (true),
    active_flag "Flag" NOT NULL CONSTRAINT "df_vendor_active_flag" DEFAULT (true),
    purchasing_web_service_url varchar(1024) NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_vendor_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_vendor_credit_rating" CHECK (credit_rating BETWEEN 1 AND 5)
  );

COMMENT ON SCHEMA purchasing IS 'Contains objects related to vendors and purchase orders.';

-- Calculated columns that needed to be there just for the CSV import
ALTER TABLE purchasing.purchase_order_detail DROP COLUMN line_total;
ALTER TABLE purchasing.purchase_order_detail DROP COLUMN stocked_qty;
ALTER TABLE purchasing.purchase_order_header DROP COLUMN total_due;

CREATE SCHEMA sales
  CREATE TABLE country_region_currency(
    country_region_code varchar(3) NOT NULL,
    currency_code char(3) NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_country_region_currency_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE credit_card(
    credit_card_id SERIAL NOT NULL, -- int
    card_type varchar(50) NOT NULL,
    card_number varchar(25) NOT NULL,
    exp_month smallint NOT NULL, -- tinyint
    exp_year smallint NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_credit_card_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE currency(
    currency_code char(3) NOT NULL,
    name "Name" NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_currency_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE currency_rate(
    currency_rate_id SERIAL NOT NULL, -- int
    currency_rate_date TIMESTAMP NOT NULL,
    from_currency_code char(3) NOT NULL,
    to_currency_code char(3) NOT NULL,
    average_rate numeric NOT NULL, -- money
    end_of_day_rate numeric NOT NULL,  -- money
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_currency_rate_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE customer(
    customer_id SERIAL NOT NULL, --  NOT FOR REPLICATION -- int
    -- A customer may either be a person, a store, or a person who works for a store
    person_id INT NULL, -- If this customer represents a person, this is non-null
    store_id INT NULL,  -- If the customer is a store, or is associated with a store then this is non-null.
    territory_id INT NULL,
    account_number VARCHAR, -- AS ISNULL('AW' + dbo.ufnLeadingZeros(CustomerID), ''),
    rowguid uuid NOT NULL CONSTRAINT "df_customer_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_customer_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE person_credit_card(
    business_entity_id INT NOT NULL,
    credit_card_id INT NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_person_credit_card_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE sales_order_detail(
    sales_order_id INT NOT NULL,
    sales_order_detail_id SERIAL NOT NULL, -- int
    carrier_tracking_number varchar(25) NULL,
    order_qty smallint NOT NULL,
    product_id INT NOT NULL,
    special_offer_id INT NOT NULL,
    unit_price numeric NOT NULL, -- money
    unit_price_discount numeric NOT NULL CONSTRAINT "df_sales_order_detail_unit_price_discount" DEFAULT (0.0), -- money
    line_total numeric, -- AS ISNULL(UnitPrice * (1.0 - UnitPriceDiscount) * OrderQty, 0.0),
    rowguid uuid NOT NULL CONSTRAINT "df_sales_order_detail_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_sales_order_detail_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_sales_order_detail_order_qty" CHECK (order_qty > 0),
    CONSTRAINT "ck_sales_order_detail_unit_price" CHECK (unit_price >= 0.00),
    CONSTRAINT "ck_sales_order_detail_unit_price_discount" CHECK (unit_price_discount >= 0.00)
  )
  CREATE TABLE sales_order_header(
    sales_order_id SERIAL NOT NULL, --  NOT FOR REPLICATION -- int
    revision_number smallint NOT NULL CONSTRAINT "df_sales_order_header_revision_number" DEFAULT (0), -- tinyint
    order_date TIMESTAMP NOT NULL CONSTRAINT "df_sales_order_header_order_date" DEFAULT (NOW()),
    due_date TIMESTAMP NOT NULL,
    ship_date TIMESTAMP NULL,
    status smallint NOT NULL CONSTRAINT "df_sales_order_header_status" DEFAULT (1), -- tinyint
    online_order_flag "Flag" NOT NULL CONSTRAINT "df_sales_order_header_online_order_flag" DEFAULT (true),
    sales_order_number VARCHAR(23), -- AS ISNULL(N'SO' + CONVERT(nvarchar(23), sales_order_id), N'*** ERROR ***'),
    purchase_order_number "OrderNumber" NULL,
    account_number "AccountNumber" NULL,
    customer_id INT NOT NULL,
    sales_person_id INT NULL,
    territory_id INT NULL,
    bill_to_address_id INT NOT NULL,
    ship_to_address_id INT NOT NULL,
    ship_method_id INT NOT NULL,
    credit_card_id INT NULL,
    credit_card_approval_code varchar(15) NULL,
    currency_rate_id INT NULL,
    sub_total numeric NOT NULL CONSTRAINT "df_sales_order_header_sub_total" DEFAULT (0.00), -- money
    tax_amt numeric NOT NULL CONSTRAINT "df_sales_order_header_tax_amt" DEFAULT (0.00), -- money
    freight numeric NOT NULL CONSTRAINT "df_sales_order_header_freight" DEFAULT (0.00), -- money
    total_due numeric, -- AS ISNULL(SubTotal + TaxAmt + Freight, 0),
    comment varchar(128) NULL,
    rowguid uuid NOT NULL CONSTRAINT "df_sales_order_header_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_sales_order_header_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_sales_order_header_status" CHECK (status BETWEEN 0 AND 8),
    CONSTRAINT "ck_sales_order_header_due_date" CHECK (due_date >= order_date),
    CONSTRAINT "ck_sales_order_header_ship_date" CHECK ((ship_date >= order_date) OR (ship_date IS NULL)),
    CONSTRAINT "ck_sales_order_header_sub_total" CHECK (sub_total >= 0.00),
    CONSTRAINT "ck_sales_order_header_tax_amt" CHECK (tax_amt >= 0.00),
    CONSTRAINT "ck_sales_order_header_freight" CHECK (freight >= 0.00)
  )
  CREATE TABLE sales_order_header_sales_reason(
    sales_order_id INT NOT NULL,
    sales_reason_id INT NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_sales_order_header_sales_reason_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE sales_person(
    business_entity_id INT NOT NULL,
    territory_id INT NULL,
    sales_quota numeric NULL, -- money
    bonus numeric NOT NULL CONSTRAINT "df_sales_person_bonus" DEFAULT (0.00), -- money
    commission_pct numeric NOT NULL CONSTRAINT "df_sales_person_commission_pct" DEFAULT (0.00), -- smallmoney -- money
    sales_ytd numeric NOT NULL CONSTRAINT "df_sales_person_sales_ytd" DEFAULT (0.00), -- money
    sales_last_year numeric NOT NULL CONSTRAINT "df_sales_person_sales_last_year" DEFAULT (0.00), -- money
    rowguid uuid NOT NULL CONSTRAINT "df_sales_person_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_sales_person_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_sales_person_sales_quota" CHECK (sales_quota > 0.00),
    CONSTRAINT "ck_sales_person_bonus" CHECK (bonus >= 0.00),
    CONSTRAINT "ck_sales_person_commission_pct" CHECK (commission_pct >= 0.00),
    CONSTRAINT "ck_sales_person_sales_ytd" CHECK (sales_ytd >= 0.00),
    CONSTRAINT "ck_sales_person_sales_last_year" CHECK (sales_last_year >= 0.00)
  )
  CREATE TABLE sales_person_quota_history(
    business_entity_id INT NOT NULL,
    quota_date TIMESTAMP NOT NULL,
    sales_quota numeric NOT NULL, -- money
    rowguid uuid NOT NULL CONSTRAINT "df_sales_person_quota_history_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_sales_person_quota_history_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_sales_person_quota_history_sales_quota" CHECK (sales_quota > 0.00)
  )
  CREATE TABLE sales_reason(
    sales_reason_id SERIAL NOT NULL, -- int
    name "Name" NOT NULL,
    reason_type "Name" NOT NULL,
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_sales_reason_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE sales_tax_rate(
    sales_tax_rate_id SERIAL NOT NULL, -- int
    state_province_id INT NOT NULL,
    tax_type smallint NOT NULL, -- tinyint
    tax_rate numeric NOT NULL CONSTRAINT "df_sales_tax_rate_tax_rate" DEFAULT (0.00), -- smallmoney -- money
    name "Name" NOT NULL,
    rowguid uuid NOT NULL CONSTRAINT "df_sales_tax_rate_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_sales_tax_rate_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_sales_tax_rate_tax_type" CHECK (tax_type BETWEEN 1 AND 3)
  )
  CREATE TABLE sales_territory(
    territory_id SERIAL NOT NULL, -- int
    name "Name" NOT NULL,
    country_region_code varchar(3) NOT NULL,
    "group" varchar(50) NOT NULL, -- Group
    sales_ytd numeric NOT NULL CONSTRAINT "df_sales_territory_sales_ytd" DEFAULT (0.00), -- money
    sales_last_year numeric NOT NULL CONSTRAINT "df_sales_territory_sales_last_year" DEFAULT (0.00), -- money
    cost_ytd numeric NOT NULL CONSTRAINT "df_sales_territory_cost_ytd" DEFAULT (0.00), -- money
    cost_last_year numeric NOT NULL CONSTRAINT "df_sales_territory_cost_last_year" DEFAULT (0.00), -- money
    rowguid uuid NOT NULL CONSTRAINT "df_sales_territory_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_sales_territory_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_sales_territory_sales_ytd" CHECK (sales_ytd >= 0.00),
    CONSTRAINT "ck_sales_territory_sales_last_year" CHECK (sales_last_year >= 0.00),
    CONSTRAINT "ck_sales_territory_cost_ytd" CHECK (cost_ytd >= 0.00),
    CONSTRAINT "ck_sales_territory_cost_last_year" CHECK (cost_last_year >= 0.00)
  )
  CREATE TABLE sales_territory_history(
    business_entity_id INT NOT NULL,  -- A sales person
    territory_id INT NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NULL,
    rowguid uuid NOT NULL CONSTRAINT "df_sales_territory_history_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_sales_territory_history_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_sales_territory_history_end_date" CHECK ((end_date >= start_date) OR (end_date IS NULL))
  )
  CREATE TABLE shopping_cart_item(
    shopping_cart_item_id SERIAL NOT NULL, -- int
    shopping_cart_id varchar(50) NOT NULL,
    quantity INT NOT NULL CONSTRAINT "df_shopping_cart_item_quantity" DEFAULT (1),
    product_id INT NOT NULL,
    date_created TIMESTAMP NOT NULL CONSTRAINT "df_shopping_cart_item_date_created" DEFAULT (NOW()),
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_shopping_cart_item_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_shopping_cart_item_quantity" CHECK (quantity >= 1)
  )
  CREATE TABLE special_offer(
    special_offer_id SERIAL NOT NULL, -- int
    description varchar(255) NOT NULL,
    discount_pct numeric NOT NULL CONSTRAINT "df_special_offer_discount_pct" DEFAULT (0.00), -- smallmoney -- money
    type varchar(50) NOT NULL,
    category varchar(50) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    min_qty INT NOT NULL CONSTRAINT "df_special_offer_min_qty" DEFAULT (0),
    max_qty INT NULL,
    rowguid uuid NOT NULL CONSTRAINT "df_special_offer_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_special_offer_modified_date" DEFAULT (NOW()),
    CONSTRAINT "ck_special_offer_end_date" CHECK (end_date >= start_date),
    CONSTRAINT "ck_special_offer_discount_pct" CHECK (discount_pct >= 0.00),
    CONSTRAINT "ck_special_offer_min_qty" CHECK (min_qty >= 0),
    CONSTRAINT "ck_special_offer_max_qty"  CHECK (max_qty >= 0)
  )
  CREATE TABLE special_offer_product(
    special_offer_id INT NOT NULL,
    product_id INT NOT NULL,
    rowguid uuid NOT NULL CONSTRAINT "df_special_offer_product_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_special_offer_product_modified_date" DEFAULT (NOW())
  )
  CREATE TABLE Store(
    business_entity_id INT NOT NULL,
    Name "Name" NOT NULL,
    sales_person_id INT NULL,
    demographics varchar NULL, -- XML(sales.store_survey_schema_collection)
    rowguid uuid NOT NULL CONSTRAINT "df_store_rowguid" DEFAULT (uuid_generate_v1()), -- ROWGUIDCOL
    modified_date TIMESTAMP NOT NULL CONSTRAINT "df_store_modified_date" DEFAULT (NOW())
  );

COMMENT ON SCHEMA sales IS 'Contains objects related to customers, sales orders, and sales territories.';


-- Calculated columns that needed to be there just for the CSV import
ALTER TABLE sales.customer DROP COLUMN account_number;
ALTER TABLE sales.sales_order_detail DROP COLUMN line_total;
ALTER TABLE sales.sales_order_header DROP COLUMN sales_order_number;



-------------------------------------
-- TABLE AND COLUMN COMMENTS
-------------------------------------

COMMENT ON TABLE person.address IS 'Street address information for customers, employees, and vendors.';
  COMMENT ON COLUMN person.address.address_id IS 'Primary key for address records.';
  COMMENT ON COLUMN person.address.address_line1 IS 'First street address line.';
  COMMENT ON COLUMN person.address.address_line2 IS 'Second street address line.';
  COMMENT ON COLUMN person.address.city IS 'Name of the city.';
  COMMENT ON COLUMN person.address.state_province_id IS 'Unique identification number for the state or province. Foreign key to state_province table.';
  COMMENT ON COLUMN person.address.postal_code IS 'Postal code for the street address.';
  COMMENT ON COLUMN person.address.spatial_location IS 'Latitude and longitude of this address.';

COMMENT ON TABLE person.address_type IS 'Types of addresses stored in the address table.';
  COMMENT ON COLUMN person.address_type.address_type_id IS 'Primary key for address_type records.';
  COMMENT ON COLUMN person.address_type.name IS 'Address type description. For example, Billing, Home, or Shipping.';

COMMENT ON TABLE production.bill_of_materials IS 'Items required to make bicycles and bicycle subassemblies. It identifies the heirarchical relationship between a parent product and its components.';
  COMMENT ON COLUMN production.bill_of_materials.bill_of_materials_id IS 'Primary key for bill_of_materials records.';
  COMMENT ON COLUMN production.bill_of_materials.product_assembly_id IS 'Parent product identification number. Foreign key to product.product_id.';
  COMMENT ON COLUMN production.bill_of_materials.component_id IS 'Component identification number. Foreign key to product.product_id.';
  COMMENT ON COLUMN production.bill_of_materials.start_date IS 'Date the component started being used in the assembly item.';
  COMMENT ON COLUMN production.bill_of_materials.end_date IS 'Date the component stopped being used in the assembly item.';
  COMMENT ON COLUMN production.bill_of_materials.unit_measure_code IS 'Standard code identifying the unit of measure for the quantity.';
  COMMENT ON COLUMN production.bill_of_materials.bom_level IS 'Indicates the depth the component is from its parent (AssemblyID).';
  COMMENT ON COLUMN production.bill_of_materials.per_assembly_qty IS 'Quantity of the component needed to create the assembly.';

COMMENT ON TABLE person.business_entity IS 'Source of the ID that connects vendors, customers, and employees with address and contact information.';
  COMMENT ON COLUMN person.business_entity.business_entity_id IS 'Primary key for all customers, vendors, and employees.';

COMMENT ON TABLE person.business_entity_address IS 'Cross-reference table mapping customers, vendors, and employees to their addresses.';
  COMMENT ON COLUMN person.business_entity_address.business_entity_id IS 'Primary key. Foreign key to business_entity.business_entity_id.';
  COMMENT ON COLUMN person.business_entity_address.address_id IS 'Primary key. Foreign key to address.address_id.';
  COMMENT ON COLUMN person.business_entity_address.address_type_id IS 'Primary key. Foreign key to address_type.address_type_id.';

COMMENT ON TABLE person.business_entity_contact IS 'Cross-reference table mapping stores, vendors, and employees to people';
  COMMENT ON COLUMN person.business_entity_contact.business_entity_id IS 'Primary key. Foreign key to business_entity.business_entity_id.';
  COMMENT ON COLUMN person.business_entity_contact.person_id IS 'Primary key. Foreign key to person.business_entity_id.';
  COMMENT ON COLUMN person.business_entity_contact.contact_type_id IS 'Primary key.  Foreign key to contact_type.contact_type_id.';

COMMENT ON TABLE person.contact_type IS 'Lookup table containing the types of business entity contacts.';
  COMMENT ON COLUMN person.contact_type.contact_type_id IS 'Primary key for ContactType records.';
  COMMENT ON COLUMN person.contact_type.name IS 'Contact type description.';

-- COMMENT ON TABLE sales.country_region_currency IS 'Cross-reference table mapping ISO currency codes to a country or region.';
  COMMENT ON COLUMN sales.country_region_currency.country_region_code IS 'ISO code for countries and regions. Foreign key to country_region.country_region_code.';
  COMMENT ON COLUMN sales.country_region_currency.currency_code IS 'ISO standard currency code. Foreign key to currency.currency_code.';

COMMENT ON TABLE person.country_region IS 'Lookup table containing the ISO standard codes for countries and regions.';
  COMMENT ON COLUMN person.country_region.country_region_code IS 'ISO standard code for countries and regions.';
  COMMENT ON COLUMN person.country_region.name IS 'Country or region name.';

COMMENT ON TABLE sales.credit_card IS 'Customer credit card information.';
  COMMENT ON COLUMN sales.credit_card.credit_card_id IS 'Primary key for credit_card records.';
  COMMENT ON COLUMN sales.credit_card.card_type IS 'Credit card name.';
  COMMENT ON COLUMN sales.credit_card.card_number IS 'Credit card number.';
  COMMENT ON COLUMN sales.credit_card.exp_month IS 'Credit card expiration month.';
  COMMENT ON COLUMN sales.credit_card.exp_year IS 'Credit card expiration year.';

COMMENT ON TABLE production.culture IS 'Lookup table containing the languages in which some AdventureWorks data is stored.';
  COMMENT ON COLUMN production.culture.culture_id IS 'Primary key for culture records.';
  COMMENT ON COLUMN production.culture.name IS 'Culture description.';

COMMENT ON TABLE sales.currency IS 'Lookup table containing standard ISO currencies.';
  COMMENT ON COLUMN sales.currency.currency_code IS 'The ISO code for the Currency.';
  COMMENT ON COLUMN sales.currency.name IS 'Currency name.';

COMMENT ON TABLE sales.currency_rate IS 'Currency exchange rates.';
  COMMENT ON COLUMN sales.currency_rate.currency_rate_id IS 'Primary key for currency_rate records.';
  COMMENT ON COLUMN sales.currency_rate.currency_rate_date IS 'Date and time the exchange rate was obtained.';
  COMMENT ON COLUMN sales.currency_rate.from_currency_code IS 'Exchange rate was converted from this currency code.';
  COMMENT ON COLUMN sales.currency_rate.to_currency_code IS 'Exchange rate was converted to this currency code.';
  COMMENT ON COLUMN sales.currency_rate.average_rate IS 'Average exchange rate for the day.';
  COMMENT ON COLUMN sales.currency_rate.end_of_day_rate IS 'Final exchange rate for the day.';

COMMENT ON TABLE sales.customer IS 'Current customer information. Also see the Person and Store tables.';
  COMMENT ON COLUMN sales.customer.customer_id IS 'Primary key.';
  COMMENT ON COLUMN sales.customer.person_id IS 'Foreign key to person.business_entity_id';
  COMMENT ON COLUMN sales.customer.store_id IS 'Foreign key to Store.business_entity_id';
  COMMENT ON COLUMN sales.customer.territory_id IS 'ID of the territory in which the customer is located. Foreign key to sales_territory.sales_territory_id.';

COMMENT ON TABLE human_resources.department IS 'Lookup table containing the departments within the Adventure Works Cycles company.';
  COMMENT ON COLUMN human_resources.department.department_id IS 'Primary key for department records.';
  COMMENT ON COLUMN human_resources.department.name IS 'Name of the department.';
  COMMENT ON COLUMN human_resources.department.group_name IS 'Name of the group to which the department belongs.';

COMMENT ON TABLE production.document IS 'Product maintenance documents.';
  COMMENT ON COLUMN production.document.document_node IS 'Primary key for document records.';
  COMMENT ON COLUMN production.document.title IS 'Title of the document.';
  COMMENT ON COLUMN production.document.owner IS 'Employee who controls the document.  Foreign key to employee.business_entity_id';
  COMMENT ON COLUMN production.document.folder_flag IS '0 = This is a folder, 1 = This is a document.';
  COMMENT ON COLUMN production.document.file_name IS 'File name of the document';
  COMMENT ON COLUMN production.document.file_extension IS 'File extension indicating the document type. For example, .doc or .txt.';
  COMMENT ON COLUMN production.document.revision IS 'Revision number of the document.';
  COMMENT ON COLUMN production.document.change_number IS 'Engineering change approval number.';
  COMMENT ON COLUMN production.document.status IS '1 = Pending approval, 2 = Approved, 3 = Obsolete';
  COMMENT ON COLUMN production.document.document_summary IS 'Document abstract.';
  COMMENT ON COLUMN production.document.document IS 'Complete document.';
  COMMENT ON COLUMN production.document.rowguid IS 'ROWGUIDCOL number uniquely identifying the record. Required for FileStream.';

COMMENT ON TABLE person.email_address IS 'Where to send a person email.';
  COMMENT ON COLUMN person.email_address.business_entity_id IS 'Primary key. Person associated with this email address.  Foreign key to person.business_entity_id';
  COMMENT ON COLUMN person.email_address.email_address_id IS 'Primary key. ID of this email address.';
  COMMENT ON COLUMN person.email_address.email_address IS 'E-mail address for the person.';

COMMENT ON TABLE human_resources.employee IS 'Employee information such as salary, department, and title.';
  COMMENT ON COLUMN human_resources.employee.business_entity_id IS 'Primary key for employee records.  Foreign key to business_entity.business_entity_id.';
  COMMENT ON COLUMN human_resources.employee.national_id_number IS 'Unique national identification number such as a social security number.';
  COMMENT ON COLUMN human_resources.employee.login_id IS 'Network login.';
  COMMENT ON COLUMN human_resources.employee.organization_node IS 'Where the employee is located in corporate hierarchy.';
  COMMENT ON COLUMN human_resources.employee.job_title IS 'Work title such as Buyer or Sales Representative.';
  COMMENT ON COLUMN human_resources.employee.birth_date IS 'Date of birth.';
  COMMENT ON COLUMN human_resources.employee.marital_status IS 'M = Married, S = Single';
  COMMENT ON COLUMN human_resources.employee.gender IS 'M = Male, F = Female';
  COMMENT ON COLUMN human_resources.employee.hire_date IS 'Employee hired on this date.';
  COMMENT ON COLUMN human_resources.employee.salaried_flag IS 'Job classification. 0 = Hourly, not exempt from collective bargaining. 1 = Salaried, exempt from collective bargaining.';
  COMMENT ON COLUMN human_resources.employee.vacation_hours IS 'Number of available vacation hours.';
  COMMENT ON COLUMN human_resources.employee.sick_leave_hours IS 'Number of available sick leave hours.';
  COMMENT ON COLUMN human_resources.employee.current_flag IS '0 = Inactive, 1 = Active';

COMMENT ON TABLE human_resources.employee_department_history IS 'Employee department transfers.';
  COMMENT ON COLUMN human_resources.employee_department_history.business_entity_id IS 'Employee identification number. Foreign key to employee.business_entity_id.';
  COMMENT ON COLUMN human_resources.employee_department_history.department_id IS 'Department in which the employee worked including currently. Foreign key to department.department_id.';
  COMMENT ON COLUMN human_resources.employee_department_history.shift_id IS 'Identifies which 8-hour shift the employee works. Foreign key to shift.shift.id.';
  COMMENT ON COLUMN human_resources.employee_department_history.start_date IS 'Date the employee started work in the department.';
  COMMENT ON COLUMN human_resources.employee_department_history.end_date IS 'Date the employee left the department. NULL = Current department.';

COMMENT ON TABLE human_resources.employee_pay_history IS 'Employee pay history.';
  COMMENT ON COLUMN human_resources.employee_pay_history.business_entity_id IS 'Employee identification number. Foreign key to employee.business_entity_id.';
  COMMENT ON COLUMN human_resources.employee_pay_history.rate_change_date IS 'Date the change in pay is effective';
  COMMENT ON COLUMN human_resources.employee_pay_history.rate IS 'Salary hourly rate.';
  COMMENT ON COLUMN human_resources.employee_pay_history.pay_frequency IS '1 = Salary received monthly, 2 = Salary received biweekly';

COMMENT ON TABLE production.illustration IS 'Bicycle assembly diagrams.';
  COMMENT ON COLUMN production.illustration.illustration_id IS 'Primary key for illustration records.';
  COMMENT ON COLUMN production.illustration.diagram IS 'Illustrations used in manufacturing instructions. Stored as XML.';

COMMENT ON TABLE human_resources.job_candidate IS 'Résumés submitted to Human Resources by job applicants.';
  COMMENT ON COLUMN human_resources.job_candidate.job_candidate_id IS 'Primary key for job_candidate records.';
  COMMENT ON COLUMN human_resources.job_candidate.business_entity_id IS 'Employee identification number if applicant was hired. Foreign key to employee.business_entity_id.';
  COMMENT ON COLUMN human_resources.job_candidate.resume IS 'Résumé in XML format.';

COMMENT ON TABLE production.location IS 'Product inventory and manufacturing locations.';
  COMMENT ON COLUMN production.location.location_id IS 'Primary key for location records.';
  COMMENT ON COLUMN production.location.name IS 'Location description.';
  COMMENT ON COLUMN production.location.cost_rate IS 'Standard hourly cost of the manufacturing location.';
  COMMENT ON COLUMN production.location.availability IS 'Work capacity (in hours) of the manufacturing location.';

COMMENT ON TABLE person.password IS 'One way hashed authentication information';
  COMMENT ON COLUMN person.password.password_hash IS 'Password for the e-mail account.';
  COMMENT ON COLUMN person.password.password_salt IS 'Random value concatenated with the password string before the password is hashed.';

COMMENT ON TABLE person.person IS 'Human beings involved with AdventureWorks: employees, customer contacts, and vendor contacts.';
  COMMENT ON COLUMN person.person.business_entity_id IS 'Primary key for person records.';
  COMMENT ON COLUMN person.person.person_type IS 'Primary type of person: SC = Store Contact, IN = Individual (retail) customer, SP = Sales person, EM = Employee (non-sales), VC = Vendor contact, GC = General contact';
  COMMENT ON COLUMN person.person.name_style IS '0 = The data in FirstName and LastName are stored in western style (first name, last name) order.  1 = Eastern style (last name, first name) order.';
  COMMENT ON COLUMN person.person.title IS 'A courtesy title. For example, Mr. or Ms.';
  COMMENT ON COLUMN person.person.first_name IS 'First name of the person.';
  COMMENT ON COLUMN person.person.middle_name IS 'Middle name or middle initial of the person.';
  COMMENT ON COLUMN person.person.last_name IS 'Last name of the person.';
  COMMENT ON COLUMN person.person.suffix IS 'Surname suffix. For example, Sr. or Jr.';
  COMMENT ON COLUMN person.person.email_promotion IS '0 = Contact does not wish to receive e-mail promotions, 1 = Contact does wish to receive e-mail promotions from AdventureWorks, 2 = Contact does wish to receive e-mail promotions from AdventureWorks and selected partners.';
  COMMENT ON COLUMN person.person.demographics IS 'Personal information such as hobbies, and income collected from online shoppers. Used for sales analysis.';
  COMMENT ON COLUMN person.person.additional_contact_info IS 'Additional contact information about the person stored in xml format.';

COMMENT ON TABLE sales.person_credit_card IS 'Cross-reference table mapping people to their credit card information in the credit_card table.';
  COMMENT ON COLUMN sales.person_credit_card.business_entity_id IS 'Business entity identification number. Foreign key to person.business_entity_id.';
  COMMENT ON COLUMN sales.person_credit_card.credit_card_id IS 'Credit card identification number. Foreign key to credit_card.credit_card_id.';

COMMENT ON TABLE person.person_phone IS 'Telephone number and type of a person.';
  COMMENT ON COLUMN person.person_phone.business_entity_id IS 'Business entity identification number. Foreign key to person.business_entity_id.';
  COMMENT ON COLUMN person.person_phone.phone_number IS 'Telephone number identification number.';
  COMMENT ON COLUMN person.person_phone.phone_number_type_id IS 'Kind of phone number. Foreign key to phone_number_type.phone_number_type_id.';

COMMENT ON TABLE person.phone_number_type IS 'Type of phone number of a person.';
  COMMENT ON COLUMN person.phone_number_type.phone_number_type_id IS 'Primary key for telephone number type records.';
  COMMENT ON COLUMN person.phone_number_type.name IS 'Name of the telephone number type';

COMMENT ON TABLE production.product IS 'Products sold or used in the manfacturing of sold products.';
  COMMENT ON COLUMN production.product.product_id IS 'Primary key for Product records.';
  COMMENT ON COLUMN production.product.name IS 'Name of the product.';
  COMMENT ON COLUMN production.product.product_number IS 'Unique product identification number.';
  COMMENT ON COLUMN production.product.make_flag IS '0 = Product is purchased, 1 = Product is manufactured in-house.';
  COMMENT ON COLUMN production.product.finished_goods_flag IS '0 = Product is not a salable item. 1 = Product is salable.';
  COMMENT ON COLUMN production.product.color IS 'Product color.';
  COMMENT ON COLUMN production.product.safety_stock_level IS 'Minimum inventory quantity.';
  COMMENT ON COLUMN production.product.reorder_point IS 'Inventory level that triggers a purchase order or work order.';
  COMMENT ON COLUMN production.product.standard_cost IS 'Standard cost of the product.';
  COMMENT ON COLUMN production.product.list_price IS 'Selling price.';
  COMMENT ON COLUMN production.product.size IS 'Product size.';
  COMMENT ON COLUMN production.product.size_unit_measure_code IS 'Unit of measure for Size column.';
  COMMENT ON COLUMN production.product.weight_unit_measure_code IS 'Unit of measure for Weight column.';
  COMMENT ON COLUMN production.product.weight IS 'Product weight.';
  COMMENT ON COLUMN production.product.days_to_manufacture IS 'Number of days required to manufacture the product.';
  COMMENT ON COLUMN production.product.product_line IS 'R = Road, M = Mountain, T = Touring, S = Standard';
  COMMENT ON COLUMN production.product.class IS 'H = High, M = Medium, L = Low';
  COMMENT ON COLUMN production.product.style IS 'W = Womens, M = Mens, U = Universal';
  COMMENT ON COLUMN production.product.product_subcategory_id IS 'Product is a member of this product subcategory. Foreign key to product_sub_category.product_sub_category_id.';
  COMMENT ON COLUMN production.product.product_model_id IS 'Product is a member of this product model. Foreign key to product_model.product_model_id.';
  COMMENT ON COLUMN production.product.sell_start_date IS 'Date the product was available for sale.';
  COMMENT ON COLUMN production.product.sell_end_date IS 'Date the product was no longer available for sale.';
  COMMENT ON COLUMN production.product.discontinued_date IS 'Date the product was discontinued.';

COMMENT ON TABLE production.product_category IS 'High-level product categorization.';
  COMMENT ON COLUMN production.product_category.product_category_id IS 'Primary key for product_category records.';
  COMMENT ON COLUMN production.product_category.name IS 'Category description.';

COMMENT ON TABLE production.product_cost_history IS 'Changes in the cost of a product over time.';
  COMMENT ON COLUMN production.product_cost_history.product_id IS 'Product identification number. Foreign key to product.product_id';
  COMMENT ON COLUMN production.product_cost_history.start_date IS 'Product cost start date.';
  COMMENT ON COLUMN production.product_cost_history.end_date IS 'Product cost end date.';
  COMMENT ON COLUMN production.product_cost_history.standard_cost IS 'Standard cost of the product.';

COMMENT ON TABLE production.product_description IS 'Product descriptions in several languages.';
  COMMENT ON COLUMN production.product_description.product_description_id IS 'Primary key for product_description records.';
  COMMENT ON COLUMN production.product_description.description IS 'Description of the product.';

COMMENT ON TABLE production.product_document IS 'Cross-reference table mapping products to related product documents.';
  COMMENT ON COLUMN production.product_document.product_id IS 'Product identification number. Foreign key to product.product_id.';
  COMMENT ON COLUMN production.product_document.document_node IS 'Document identification number. Foreign key to document.document_node.';

COMMENT ON TABLE production.product_inventory IS 'Product inventory information.';
  COMMENT ON COLUMN production.product_inventory.product_id IS 'Product identification number. Foreign key to product.product_id.';
  COMMENT ON COLUMN production.product_inventory.location_id IS 'Inventory location identification number. Foreign key to location.location_id.';
  COMMENT ON COLUMN production.product_inventory.shelf IS 'Storage compartment within an inventory location.';
  COMMENT ON COLUMN production.product_inventory.bin IS 'Storage container on a shelf in an inventory location.';
  COMMENT ON COLUMN production.product_inventory.quantity IS 'Quantity of products in the inventory location.';

COMMENT ON TABLE production.product_list_price_history IS 'Changes in the list price of a product over time.';
  COMMENT ON COLUMN production.product_list_price_history.product_id IS 'Product identification number. Foreign key to product.product_id';
  COMMENT ON COLUMN production.product_list_price_history.start_date IS 'List price start date.';
  COMMENT ON COLUMN production.product_list_price_history.end_date IS 'List price end date';
  COMMENT ON COLUMN production.product_list_price_history.list_price IS 'Product list price.';

COMMENT ON TABLE production.product_model IS 'Product model classification.';
  COMMENT ON COLUMN production.product_model.product_model_id IS 'Primary key for ProductModel records.';
  COMMENT ON COLUMN production.product_model.name IS 'Product model description.';
  COMMENT ON COLUMN production.product_model.catalog_description IS 'Detailed product catalog information in xml format.';
  COMMENT ON COLUMN production.product_model.instructions IS 'Manufacturing instructions in xml format.';

COMMENT ON TABLE production.product_model_illustration IS 'Cross-reference table mapping product models and illustrations.';
  COMMENT ON COLUMN production.product_model_illustration.product_model_id IS 'Primary key. Foreign key to product_model.product_model_id.';
  COMMENT ON COLUMN production.product_model_illustration.illustration_id IS 'Primary key. Foreign key to illustration.illustration_id.';

COMMENT ON TABLE production.product_model_product_description_culture IS 'Cross-reference table mapping product descriptions and the language the description is written in.';
  COMMENT ON COLUMN production.product_model_product_description_culture.product_model_id IS 'Primary key. Foreign key to product_model.product_model_id.';
  COMMENT ON COLUMN production.product_model_product_description_culture.product_description_id IS 'Primary key. Foreign key to product_description.product_description_id.';
  COMMENT ON COLUMN production.product_model_product_description_culture.culture_id IS 'Culture identification number. Foreign key to culture.culture_id.';

COMMENT ON TABLE production.product_photo IS 'Product images.';
  COMMENT ON COLUMN production.product_photo.product_photo_id IS 'Primary key for product_photo records.';
  COMMENT ON COLUMN production.product_photo.thumb_nail_photo IS 'Small image of the product.';
  COMMENT ON COLUMN production.product_photo.thumbnail_photo_file_name IS 'Small image file name.';
  COMMENT ON COLUMN production.product_photo.large_photo IS 'Large image of the product.';
  COMMENT ON COLUMN production.product_photo.large_photo_file_name IS 'Large image file name.';

COMMENT ON TABLE production.product_product_photo IS 'Cross-reference table mapping products and product photos.';
  COMMENT ON COLUMN production.product_product_photo.product_id IS 'Product identification number. Foreign key to product.product_id.';
  COMMENT ON COLUMN production.product_product_photo.product_photo_id IS 'Product photo identification number. Foreign key to product_photo.product_photo_id.';
  COMMENT ON COLUMN production.product_product_photo.primary IS '0 = Photo is not the principal image. 1 = Photo is the principal image.';

COMMENT ON TABLE production.product_review IS 'Customer reviews of products they have purchased.';
  COMMENT ON COLUMN production.product_review.product_review_id IS 'Primary key for product_review records.';
  COMMENT ON COLUMN production.product_review.product_id IS 'Product identification number. Foreign key to product.product_id.';
  COMMENT ON COLUMN production.product_review.reviewer_name IS 'Name of the reviewer.';
  COMMENT ON COLUMN production.product_review.review_date IS 'Date review was submitted.';
  COMMENT ON COLUMN production.product_review.email_address IS 'Reviewer''s e-mail address.';
  COMMENT ON COLUMN production.product_review.rating IS 'Product rating given by the reviewer. Scale is 1 to 5 with 5 as the highest rating.';
  COMMENT ON COLUMN production.product_review.comments IS 'Reviewer''s comments';

COMMENT ON TABLE production.product_subcategory IS 'Product subcategories. See product_category table.';
  COMMENT ON COLUMN production.product_subcategory.product_subcategory_id IS 'Primary key for product_subcategory records.';
  COMMENT ON COLUMN production.product_subcategory.product_category_id IS 'Product category identification number. Foreign key to product_category.product_category_id.';
  COMMENT ON COLUMN production.product_subcategory.name IS 'Subcategory description.';

COMMENT ON TABLE purchasing.product_vendor IS 'Cross-reference table mapping vendors with the products they supply.';
  COMMENT ON COLUMN purchasing.product_vendor.product_id IS 'Primary key. Foreign key to product.product_id.';
  COMMENT ON COLUMN purchasing.product_vendor.business_entity_id IS 'Primary key. Foreign key to vendor.business_entity_id.';
  COMMENT ON COLUMN purchasing.product_vendor.average_lead_time IS 'The average span of time (in days) between placing an order with the vendor and receiving the purchased product.';
  COMMENT ON COLUMN purchasing.product_vendor.standard_price IS 'The vendor''s usual selling price.';
  COMMENT ON COLUMN purchasing.product_vendor.last_receipt_cost IS 'The selling price when last purchased.';
  COMMENT ON COLUMN purchasing.product_vendor.last_receipt_date IS 'Date the product was last received by the vendor.';
  COMMENT ON COLUMN purchasing.product_vendor.min_order_qty IS 'The maximum quantity that should be ordered.';
  COMMENT ON COLUMN purchasing.product_vendor.max_order_qty IS 'The minimum quantity that should be ordered.';
  COMMENT ON COLUMN purchasing.product_vendor.on_order_qty IS 'The quantity currently on order.';
  COMMENT ON COLUMN purchasing.product_vendor.unit_measure_code IS 'The product''s unit of measure.';

COMMENT ON TABLE purchasing.purchase_order_detail IS 'Individual products associated with a specific purchase order. See purchase_order_header.';
  COMMENT ON COLUMN purchasing.purchase_order_detail.purchase_order_id IS 'Primary key. Foreign key to purchase_order_header.purchase_order_id.';
  COMMENT ON COLUMN purchasing.purchase_order_detail.purchase_order_detail_id IS 'Primary key. One line number per purchased product.';
  COMMENT ON COLUMN purchasing.purchase_order_detail.due_date IS 'Date the product is expected to be received.';
  COMMENT ON COLUMN purchasing.purchase_order_detail.order_qty IS 'Quantity ordered.';
  COMMENT ON COLUMN purchasing.purchase_order_detail.product_id IS 'Product identification number. Foreign key to product.product_id.';
  COMMENT ON COLUMN purchasing.purchase_order_detail.unit_price IS 'Vendor''s selling price of a single product.';
  COMMENT ON COLUMN purchasing.purchase_order_detail.received_qty IS 'Quantity actually received from the vendor.';
  COMMENT ON COLUMN purchasing.purchase_order_detail.rejected_qty IS 'Quantity rejected during inspection.';

COMMENT ON TABLE purchasing.purchase_order_header IS 'General purchase order information. See purchase_order_detail.';
  COMMENT ON COLUMN purchasing.purchase_order_header.purchase_order_id IS 'Primary key.';
  COMMENT ON COLUMN purchasing.purchase_order_header.revision_number IS 'Incremental number to track changes to the purchase order over time.';
  COMMENT ON COLUMN purchasing.purchase_order_header.status IS 'Order current status. 1 = Pending; 2 = Approved; 3 = Rejected; 4 = Complete';
  COMMENT ON COLUMN purchasing.purchase_order_header.employee_id IS 'Employee who created the purchase order. Foreign key to employee.business_entity_id.';
  COMMENT ON COLUMN purchasing.purchase_order_header.vendor_id IS 'Vendor with whom the purchase order is placed. Foreign key to vendor.business_entity_id.';
  COMMENT ON COLUMN purchasing.purchase_order_header.ship_method_id IS 'Shipping method. Foreign key to ship_method.ship_method_id.';
  COMMENT ON COLUMN purchasing.purchase_order_header.order_date IS 'Purchase order creation date.';
  COMMENT ON COLUMN purchasing.purchase_order_header.ship_date IS 'Estimated shipment date from the vendor.';
  COMMENT ON COLUMN purchasing.purchase_order_header.sub_total IS 'Purchase order subtotal. Computed as SUM(purchase_order_detail.LineTotal)for the appropriate purchase_order_id.';
  COMMENT ON COLUMN purchasing.purchase_order_header.tax_amt IS 'Tax amount.';
  COMMENT ON COLUMN purchasing.purchase_order_header.freight IS 'Shipping cost.';

COMMENT ON TABLE sales.sales_order_detail IS 'Individual products associated with a specific sales order. See sales_order_header.';
  COMMENT ON COLUMN sales.sales_order_detail.sales_order_id IS 'Primary key. Foreign key to sales_order_header.sales_order_id.';
  COMMENT ON COLUMN sales.sales_order_detail.sales_order_detail_id IS 'Primary key. One incremental unique number per product sold.';
  COMMENT ON COLUMN sales.sales_order_detail.carrier_tracking_number IS 'Shipment tracking number supplied by the shipper.';
  COMMENT ON COLUMN sales.sales_order_detail.order_qty IS 'Quantity ordered per product.';
  COMMENT ON COLUMN sales.sales_order_detail.product_id IS 'Product sold to customer. Foreign key to product.product_id.';
  COMMENT ON COLUMN sales.sales_order_detail.special_offer_id IS 'Promotional code. Foreign key to special_offer.special_offer_id.';
  COMMENT ON COLUMN sales.sales_order_detail.unit_price IS 'Selling price of a single product.';
  COMMENT ON COLUMN sales.sales_order_detail.unit_price_discount IS 'Discount amount.';

COMMENT ON TABLE sales.sales_order_header IS 'General sales order information.';
  COMMENT ON COLUMN sales.sales_order_header.sales_order_id IS 'Primary key.';
  COMMENT ON COLUMN sales.sales_order_header.revision_number IS 'Incremental number to track changes to the sales order over time.';
  COMMENT ON COLUMN sales.sales_order_header.order_date IS 'Dates the sales order was created.';
  COMMENT ON COLUMN sales.sales_order_header.due_date IS 'Date the order is due to the customer.';
  COMMENT ON COLUMN sales.sales_order_header.ship_date IS 'Date the order was shipped to the customer.';
  COMMENT ON COLUMN sales.sales_order_header.status IS 'Order current status. 1 = In process; 2 = Approved; 3 = Backordered; 4 = Rejected; 5 = Shipped; 6 = Cancelled';
  COMMENT ON COLUMN sales.sales_order_header.online_order_flag IS '0 = Order placed by sales person. 1 = Order placed online by customer.';
  COMMENT ON COLUMN sales.sales_order_header.purchase_order_number IS 'Customer purchase order number reference.';
  COMMENT ON COLUMN sales.sales_order_header.account_number IS 'Financial accounting number reference.';
  COMMENT ON COLUMN sales.sales_order_header.customer_id IS 'Customer identification number. Foreign key to customer.business_entity_id.';
  COMMENT ON COLUMN sales.sales_order_header.sales_person_id IS 'Sales person who created the sales order. Foreign key to sales_person.business_entity_id.';
  COMMENT ON COLUMN sales.sales_order_header.territory_id IS 'Territory in which the sale was made. Foreign key to sales_territory.sales_territory_id.';
  COMMENT ON COLUMN sales.sales_order_header.bill_to_address_id IS 'Customer billing address. Foreign key to address.address_id.';
  COMMENT ON COLUMN sales.sales_order_header.ship_to_address_id IS 'Customer shipping address. Foreign key to address.address_id.';
  COMMENT ON COLUMN sales.sales_order_header.ship_method_id IS 'Shipping method. Foreign key to ship_method.ship_method_id.';
  COMMENT ON COLUMN sales.sales_order_header.credit_card_id IS 'Credit card identification number. Foreign key to credit_card.credit_card_id.';
  COMMENT ON COLUMN sales.sales_order_header.credit_card_approval_code IS 'Approval code provided by the credit card company.';
  COMMENT ON COLUMN sales.sales_order_header.currency_rate_id IS 'Currency exchange rate used. Foreign key to currency_rate.currency_rate_id.';
  COMMENT ON COLUMN sales.sales_order_header.sub_total IS 'Sales subtotal. Computed as SUM(sales_order_detail.line_total)for the appropriate sales_order_id.';
  COMMENT ON COLUMN sales.sales_order_header.tax_amt IS 'Tax amount.';
  COMMENT ON COLUMN sales.sales_order_header.freight IS 'Shipping cost.';
  COMMENT ON COLUMN sales.sales_order_header.total_due IS 'Total due from customer. Computed as subtotal + tax_amt + freight.';
  COMMENT ON COLUMN sales.sales_order_header.comment IS 'Sales representative comments.';

COMMENT ON TABLE sales.sales_order_header_sales_reason IS 'Cross-reference table mapping sales orders to sales reason codes.';
  COMMENT ON COLUMN sales.sales_order_header_sales_reason.sales_order_id IS 'Primary key. Foreign key to sales_order_header.sales_order_id.';
  COMMENT ON COLUMN sales.sales_order_header_sales_reason.sales_reason_id IS 'Primary key. Foreign key to sales_reason.sales_reason_id.';

COMMENT ON TABLE sales.sales_person IS 'Sales representative current information.';
  COMMENT ON COLUMN sales.sales_person.business_entity_id IS 'Primary key for SalesPerson records. Foreign key to employee.business_entity_id';
  COMMENT ON COLUMN sales.sales_person.territory_id IS 'Territory currently assigned to. Foreign key to sales_territory.sales_territory_id.';
  COMMENT ON COLUMN sales.sales_person.sales_quota IS 'Projected yearly sales.';
  COMMENT ON COLUMN sales.sales_person.bonus IS 'Bonus due if quota is met.';
  COMMENT ON COLUMN sales.sales_person.commission_pct IS 'Commision percent received per sale.';
  COMMENT ON COLUMN sales.sales_person.sales_ytd IS 'Sales total year to date.';
  COMMENT ON COLUMN sales.sales_person.sales_last_year IS 'Sales total of previous year.';

COMMENT ON TABLE sales.sales_person_quota_history IS 'Sales performance tracking.';
  COMMENT ON COLUMN sales.sales_person_quota_history.business_entity_id IS 'Sales person identification number. Foreign key to sales_person.business_entity_id.';
  COMMENT ON COLUMN sales.sales_person_quota_history.quota_date IS 'Sales quota date.';
  COMMENT ON COLUMN sales.sales_person_quota_history.sales_quota IS 'Sales quota amount.';

COMMENT ON TABLE sales.sales_reason IS 'Lookup table of customer purchase reasons.';
  COMMENT ON COLUMN sales.sales_reason.sales_reason_id IS 'Primary key for sales_reason records.';
  COMMENT ON COLUMN sales.sales_reason.name IS 'Sales reason description.';
  COMMENT ON COLUMN sales.sales_reason.reason_type IS 'Category the sales reason belongs to.';

COMMENT ON TABLE sales.sales_tax_rate IS 'Tax rate lookup table.';
  COMMENT ON COLUMN sales.sales_tax_rate.sales_tax_rate_id IS 'Primary key for sales_tax_rate records.';
  COMMENT ON COLUMN sales.sales_tax_rate.state_province_id IS 'State, province, or country/region the sales tax applies to.';
  COMMENT ON COLUMN sales.sales_tax_rate.tax_type IS '1 = Tax applied to retail transactions, 2 = Tax applied to wholesale transactions, 3 = Tax applied to all sales (retail and wholesale) transactions.';
  COMMENT ON COLUMN sales.sales_tax_rate.tax_rate IS 'Tax rate amount.';
  COMMENT ON COLUMN sales.sales_tax_rate.name IS 'Tax rate description.';

COMMENT ON TABLE sales.sales_territory IS 'Sales territory lookup table.';
  COMMENT ON COLUMN sales.sales_territory.territory_id IS 'Primary key for sales_territory records.';
  COMMENT ON COLUMN sales.sales_territory.name IS 'Sales territory description';
  COMMENT ON COLUMN sales.sales_territory.country_region_code IS 'ISO standard country or region code. Foreign key to country_region.country_region_code.';
  COMMENT ON COLUMN sales.sales_territory.group IS 'Geographic area to which the sales territory belong.';
  COMMENT ON COLUMN sales.sales_territory.sales_ytd IS 'Sales in the territory year to date.';
  COMMENT ON COLUMN sales.sales_territory.sales_last_year IS 'Sales in the territory the previous year.';
  COMMENT ON COLUMN sales.sales_territory.cost_ytd IS 'Business costs in the territory year to date.';
  COMMENT ON COLUMN sales.sales_territory.cost_last_year IS 'Business costs in the territory the previous year.';

COMMENT ON TABLE sales.sales_territory_history IS 'Sales representative transfers to other sales territories.';
  COMMENT ON COLUMN sales.sales_territory_history.business_entity_id IS 'Primary key. The sales rep.  Foreign key to sales_person.business_entity_id.';
  COMMENT ON COLUMN sales.sales_territory_history.territory_id IS 'Primary key. Territory identification number. Foreign key to sales_territory.sales_territory_id.';
  COMMENT ON COLUMN sales.sales_territory_history.start_date IS 'Primary key. Date the sales representive started work in the territory.';
  COMMENT ON COLUMN sales.sales_territory_history.end_date IS 'Date the sales representative left work in the territory.';

COMMENT ON TABLE production.scrap_reason IS 'Manufacturing failure reasons lookup table.';
  COMMENT ON COLUMN production.scrap_reason.scrap_reason_id IS 'Primary key for scrap_reason records.';
  COMMENT ON COLUMN production.scrap_reason.name IS 'Failure description.';

COMMENT ON TABLE human_resources.shift IS 'Work shift lookup table.';
  COMMENT ON COLUMN human_resources.shift.shift_id IS 'Primary key for Shift records.';
  COMMENT ON COLUMN human_resources.shift.name IS 'Shift description.';
  COMMENT ON COLUMN human_resources.shift.start_time IS 'Shift start time.';
  COMMENT ON COLUMN human_resources.shift.end_time IS 'Shift end time.';

COMMENT ON TABLE purchasing.ship_method IS 'Shipping company lookup table.';
  COMMENT ON COLUMN purchasing.ship_method.ship_method_id IS 'Primary key for ship_method records.';
  COMMENT ON COLUMN purchasing.ship_method.name IS 'Shipping company name.';
  COMMENT ON COLUMN purchasing.ship_method.ship_base IS 'Minimum shipping charge.';
  COMMENT ON COLUMN purchasing.ship_method.ship_rate IS 'Shipping charge per pound.';

COMMENT ON TABLE sales.shopping_cart_item IS 'Contains online customer orders until the order is submitted or cancelled.';
  COMMENT ON COLUMN sales.shopping_cart_item.shopping_cart_item_id IS 'Primary key for shopping_cart_item records.';
  COMMENT ON COLUMN sales.shopping_cart_item.shopping_cart_id IS 'Shopping cart identification number.';
  COMMENT ON COLUMN sales.shopping_cart_item.quantity IS 'Product quantity ordered.';
  COMMENT ON COLUMN sales.shopping_cart_item.product_id IS 'Product ordered. Foreign key to product.product_id.';
  COMMENT ON COLUMN sales.shopping_cart_item.date_created IS 'Date the time the record was created.';

COMMENT ON TABLE sales.special_offer IS 'Sale discounts lookup table.';
  COMMENT ON COLUMN sales.special_offer.special_offer_id IS 'Primary key for special_offer records.';
  COMMENT ON COLUMN sales.special_offer.description IS 'Discount description.';
  COMMENT ON COLUMN sales.special_offer.discount_pct IS 'Discount precentage.';
  COMMENT ON COLUMN sales.special_offer.type IS 'Discount type category.';
  COMMENT ON COLUMN sales.special_offer.category IS 'Group the discount applies to such as Reseller or Customer.';
  COMMENT ON COLUMN sales.special_offer.start_date IS 'Discount start date.';
  COMMENT ON COLUMN sales.special_offer.end_date IS 'Discount end date.';
  COMMENT ON COLUMN sales.special_offer.min_qty IS 'Minimum discount percent allowed.';
  COMMENT ON COLUMN sales.special_offer.max_qty IS 'Maximum discount percent allowed.';

COMMENT ON TABLE sales.special_offer_product IS 'Cross-reference table mapping products to special offer discounts.';
  COMMENT ON COLUMN sales.special_offer_product.special_offer_id IS 'Primary key for special_offer_product records.';
  COMMENT ON COLUMN sales.special_offer_product.product_id IS 'Product identification number. Foreign key to product.product_id.';

COMMENT ON TABLE person.state_province IS 'State and province lookup table.';
  COMMENT ON COLUMN person.state_province.state_province_id IS 'Primary key for state_province records.';
  COMMENT ON COLUMN person.state_province.state_province_code IS 'ISO standard state or province code.';
  COMMENT ON COLUMN person.state_province.country_region_code IS 'ISO standard country or region code. Foreign key to country_region.country_region_code.';
  COMMENT ON COLUMN person.state_province.is_only_state_province_flag IS '0 = state_province_code exists. 1 = state_province_code unavailable, using country_region_code.';
  COMMENT ON COLUMN person.state_province.name IS 'State or province description.';
  COMMENT ON COLUMN person.state_province.territory_id IS 'ID of the territory in which the state or province is located. Foreign key to sales_territory.sales_territory_id.';

COMMENT ON TABLE sales.store IS 'Customers (resellers) of Adventure Works products.';
  COMMENT ON COLUMN sales.store.business_entity_id IS 'Primary key. Foreign key to customer.business_entity_id.';
  COMMENT ON COLUMN sales.store.name IS 'Name of the store.';
  COMMENT ON COLUMN sales.store.sales_person_id IS 'ID of the sales person assigned to the customer. Foreign key to sales_person.business_entity_id.';
  COMMENT ON COLUMN sales.store.demographics IS 'Demographic informationg about the store such as the number of employees, annual sales and store type.';


COMMENT ON TABLE production.transaction_history IS 'Record of each purchase order, sales order, or work order transaction year to date.';
  COMMENT ON COLUMN production.transaction_history.transaction_id IS 'Primary key for TransactionHistory records.';
  COMMENT ON COLUMN production.transaction_history.product_id IS 'Product identification number. Foreign key to product.product_id.';
  COMMENT ON COLUMN production.transaction_history.reference_order_id IS 'Purchase order, sales order, or work order identification number.';
  COMMENT ON COLUMN production.transaction_history.reference_order_line_id IS 'Line number associated with the purchase order, sales order, or work order.';
  COMMENT ON COLUMN production.transaction_history.transaction_date IS 'Date and time of the transaction.';
  COMMENT ON COLUMN production.transaction_history.transaction_type IS 'W = Work Order, S = Sales Order, P = Purchase Order';
  COMMENT ON COLUMN production.transaction_history.quantity IS 'Product quantity.';
  COMMENT ON COLUMN production.transaction_history.actual_cost IS 'Product cost.';

COMMENT ON TABLE production.transaction_history_archive IS 'Transactions for previous years.';
  COMMENT ON COLUMN production.transaction_history_archive.transaction_id IS 'Primary key for transaction_history_archive records.';
  COMMENT ON COLUMN production.transaction_history_archive.product_id IS 'Product identification number. Foreign key to product.product_id.';
  COMMENT ON COLUMN production.transaction_history_archive.reference_order_id IS 'Purchase order, sales order, or work order identification number.';
  COMMENT ON COLUMN production.transaction_history_archive.reference_order_line_id IS 'Line number associated with the purchase order, sales order, or work order.';
  COMMENT ON COLUMN production.transaction_history_archive.transaction_date IS 'Date and time of the transaction.';
  COMMENT ON COLUMN production.transaction_history_archive.transaction_type IS 'W = Work Order, S = Sales Order, P = Purchase Order';
  COMMENT ON COLUMN production.transaction_history_archive.quantity IS 'Product quantity.';
  COMMENT ON COLUMN production.transaction_history_archive.actual_cost IS 'Product cost.';

COMMENT ON TABLE production.unit_measure IS 'Unit of measure lookup table.';
  COMMENT ON COLUMN production.unit_measure.unit_measure_code IS 'Primary key.';
  COMMENT ON COLUMN production.unit_measure.name IS 'Unit of measure description.';

COMMENT ON TABLE purchasing.vendor IS 'Companies from whom Adventure Works Cycles purchases parts or other goods.';
  COMMENT ON COLUMN purchasing.vendor.business_entity_id IS 'Primary key for Vendor records.  Foreign key to business_entity.business_entity_id';
  COMMENT ON COLUMN purchasing.vendor.account_number IS 'Vendor account (identification) number.';
  COMMENT ON COLUMN purchasing.vendor.name IS 'Company name.';
  COMMENT ON COLUMN purchasing.vendor.credit_rating IS '1 = Superior, 2 = Excellent, 3 = Above average, 4 = Average, 5 = Below average';
  COMMENT ON COLUMN purchasing.vendor.preferred_vendor_status IS '0 = Do not use if another vendor is available. 1 = Preferred over other vendors supplying the same product.';
  COMMENT ON COLUMN purchasing.vendor.active_flag IS '0 = Vendor no longer used. 1 = Vendor is actively used.';
  COMMENT ON COLUMN purchasing.vendor.purchasing_web_service_url IS 'Vendor URL.';

COMMENT ON TABLE production.work_order IS 'Manufacturing work orders.';
  COMMENT ON COLUMN production.work_order.work_order_id IS 'Primary key for work_order records.';
  COMMENT ON COLUMN production.work_order.product_id IS 'Product identification number. Foreign key to product.product_id.';
  COMMENT ON COLUMN production.work_order.order_qty IS 'Product quantity to build.';
  COMMENT ON COLUMN production.work_order.scrapped_qty IS 'Quantity that failed inspection.';
  COMMENT ON COLUMN production.work_order.start_date IS 'Work order start date.';
  COMMENT ON COLUMN production.work_order.end_date IS 'Work order end date.';
  COMMENT ON COLUMN production.work_order.due_date IS 'Work order due date.';
  COMMENT ON COLUMN production.work_order.scrap_reason_id IS 'Reason for inspection failure.';

COMMENT ON TABLE production.work_order_routing IS 'Work order details.';
  COMMENT ON COLUMN production.work_order_routing.work_order_id IS 'Primary key. Foreign key to work_order.work_order_id.';
  COMMENT ON COLUMN production.work_order_routing.product_id IS 'Primary key. Foreign key to product.product_id.';
  COMMENT ON COLUMN production.work_order_routing.operation_sequence IS 'Primary key. Indicates the manufacturing process sequence.';
  COMMENT ON COLUMN production.work_order_routing.location_id IS 'Manufacturing location where the part is processed. Foreign key to location.location_id.';
  COMMENT ON COLUMN production.work_order_routing.scheduled_start_date IS 'Planned manufacturing start date.';
  COMMENT ON COLUMN production.work_order_routing.scheduled_end_date IS 'Planned manufacturing end date.';
  COMMENT ON COLUMN production.work_order_routing.actual_start_date IS 'Actual start date.';
  COMMENT ON COLUMN production.work_order_routing.actual_end_date IS 'Actual end date.';
  COMMENT ON COLUMN production.work_order_routing.actual_resource_hrs IS 'Number of manufacturing hours used.';
  COMMENT ON COLUMN production.work_order_routing.planned_cost IS 'Estimated manufacturing cost.';
  COMMENT ON COLUMN production.work_order_routing.actual_cost IS 'Actual manufacturing cost.';



-------------------------------------
-- PRIMARY KEYS
-------------------------------------

ALTER TABLE person.address ADD CONSTRAINT "pk_address_address_id" PRIMARY KEY (address_id);
CLUSTER person.address USING "pk_address_address_id";

ALTER TABLE person.address_type ADD CONSTRAINT "pk_address_type_address_type_id" PRIMARY KEY (address_type_id);
CLUSTER person.address_type USING "pk_address_type_address_type_id";

ALTER TABLE production.bill_of_materials ADD CONSTRAINT "pk_bill_of_materials_bill_of_materials_id" PRIMARY KEY (bill_of_materials_id);

ALTER TABLE person.business_entity ADD CONSTRAINT "pk_business_entity_business_entity_id" PRIMARY KEY (business_entity_id);
CLUSTER person.business_entity USING "pk_business_entity_business_entity_id";

ALTER TABLE person.business_entity_address ADD CONSTRAINT "pk_business_entity_address_business_entity_id_address_id_address_type" PRIMARY KEY (business_entity_id, address_id, address_type_id);
CLUSTER person.business_entity_address USING "pk_business_entity_address_business_entity_id_address_id_address_type";

ALTER TABLE person.business_entity_contact ADD CONSTRAINT "pk_business_entity_contact_business_entity_id_person_id_contact_type_i" PRIMARY KEY (business_entity_id, person_id, contact_type_id);
CLUSTER person.business_entity_contact USING "pk_business_entity_contact_business_entity_id_person_id_contact_type_i";

ALTER TABLE person.contact_type ADD CONSTRAINT "pk_contact_type_contact_type_id" PRIMARY KEY (contact_type_id);
CLUSTER person.contact_type USING "pk_contact_type_contact_type_id";

ALTER TABLE sales.country_region_currency ADD CONSTRAINT "pk_country_region_currency_country_region_code_currency_code" PRIMARY KEY (country_region_code, currency_code);
CLUSTER sales.country_region_currency USING "pk_country_region_currency_country_region_code_currency_code";

ALTER TABLE person.country_region ADD CONSTRAINT "pk_country_region_country_region_code" PRIMARY KEY (country_region_code);
CLUSTER person.country_region USING "pk_country_region_country_region_code";

ALTER TABLE sales.credit_card ADD CONSTRAINT "pk_credit_card_credit_card_id" PRIMARY KEY (credit_card_id);
CLUSTER sales.credit_card USING "pk_credit_card_credit_card_id";

ALTER TABLE sales.currency ADD CONSTRAINT "pk_currency_currency_code" PRIMARY KEY (currency_code);
CLUSTER sales.Currency USING "pk_currency_currency_code";

ALTER TABLE sales.currency_rate ADD CONSTRAINT "pk_currency_rate_currency_rate_id" PRIMARY KEY (currency_rate_id);
CLUSTER sales.currency_rate USING "pk_currency_rate_currency_rate_id";

ALTER TABLE sales.customer ADD CONSTRAINT "pk_customer_customer_id" PRIMARY KEY (customer_id);
CLUSTER sales.customer USING "pk_customer_customer_id";

ALTER TABLE production.culture ADD CONSTRAINT "pk_culture_culture_id" PRIMARY KEY (culture_id);
CLUSTER production.Culture USING "pk_culture_culture_id";

ALTER TABLE production.document ADD CONSTRAINT "pk_document_document_node" PRIMARY KEY (document_node);
CLUSTER production.document USING "pk_document_document_node";

ALTER TABLE person.email_address ADD CONSTRAINT "pk_email_address_business_entity_id_email_address_id" PRIMARY KEY (business_entity_id, email_address_id);
CLUSTER person.email_address USING "pk_email_address_business_entity_id_email_address_id";

ALTER TABLE human_resources.department ADD CONSTRAINT "pk_department_department_id" PRIMARY KEY (department_id);
CLUSTER human_resources.department USING "pk_department_department_id";

ALTER TABLE human_resources.employee ADD CONSTRAINT "pk_employee_business_entity_id" PRIMARY KEY (business_entity_id);
CLUSTER human_resources.employee USING "pk_employee_business_entity_id";

ALTER TABLE human_resources.employee_department_history ADD CONSTRAINT "pk_employee_department_history_business_entity_id_start_date_departm" PRIMARY KEY (business_entity_id, start_date, department_id, shift_id);
CLUSTER human_resources.employee_department_history USING "pk_employee_department_history_business_entity_id_start_date_departm";

ALTER TABLE human_resources.employee_pay_history ADD CONSTRAINT "pk_employee_pay_history_business_entity_id_rate_change_date" PRIMARY KEY (business_entity_id, rate_change_date);
CLUSTER human_resources.employee_pay_history USING "pk_employee_pay_history_business_entity_id_rate_change_date";

ALTER TABLE human_resources.job_candidate ADD CONSTRAINT "pk_job_candidate_job_candidate_id" PRIMARY KEY (job_candidate_id);
CLUSTER human_resources.job_candidate USING "pk_job_candidate_job_candidate_id";

ALTER TABLE production.illustration ADD CONSTRAINT "pk_illustration_illustration_id" PRIMARY KEY (illustration_id);
CLUSTER production.illustration USING "pk_illustration_illustration_id";

ALTER TABLE production.location ADD CONSTRAINT "pk_location_location_id" PRIMARY KEY (location_id);
CLUSTER production.location USING "pk_location_location_id";

ALTER TABLE person.password ADD CONSTRAINT "pk_password_business_entity_id" PRIMARY KEY (business_entity_id);
CLUSTER person.password USING "pk_password_business_entity_id";

ALTER TABLE person.person ADD CONSTRAINT "pk_person_business_entity_id" PRIMARY KEY (business_entity_id);
CLUSTER person.person USING "pk_person_business_entity_id";

ALTER TABLE person.person_phone ADD CONSTRAINT "pk_person_phone_business_entity_id_phone_number_phone_number_type_id" PRIMARY KEY (business_entity_id, phone_number, phone_number_type_id);
CLUSTER person.person_phone USING "pk_person_phone_business_entity_id_phone_number_phone_number_type_id";

ALTER TABLE person.phone_number_type ADD CONSTRAINT "pk_phone_number_type_phone_number_type_id" PRIMARY KEY (phone_number_type_id);
CLUSTER person.phone_number_type USING "pk_phone_number_type_phone_number_type_id";

ALTER TABLE production.product ADD CONSTRAINT "pk_product_product_id" PRIMARY KEY (product_id);
CLUSTER production.product USING "pk_product_product_id";

ALTER TABLE production.product_category ADD CONSTRAINT "pk_product_category_product_category_id" PRIMARY KEY (product_category_id);
CLUSTER production.product_category USING "pk_product_category_product_category_id";

ALTER TABLE production.product_cost_history ADD CONSTRAINT "pk_product_cost_history_product_id_start_date" PRIMARY KEY (product_id, start_date);
CLUSTER production.product_cost_history USING "pk_product_cost_history_product_id_start_date";

ALTER TABLE production.product_description ADD CONSTRAINT "pk_product_description_product_description_id" PRIMARY KEY (product_description_id);
CLUSTER production.product_description USING "pk_product_description_product_description_id";

ALTER TABLE production.product_document ADD CONSTRAINT "pk_product_document_product_id_document_node" PRIMARY KEY (product_id, document_node);
CLUSTER production.product_document USING "pk_product_document_product_id_document_node";

ALTER TABLE production.product_inventory ADD CONSTRAINT "pk_product_inventory_product_id_location_id" PRIMARY KEY (product_id, location_id);
CLUSTER production.product_inventory USING "pk_product_inventory_product_id_location_id";

ALTER TABLE production.product_list_price_history ADD CONSTRAINT "pk_product_list_price_history_product_id_start_date" PRIMARY KEY (product_id, start_date);
CLUSTER production.product_list_price_history USING "pk_product_list_price_history_product_id_start_date";

ALTER TABLE production.product_model ADD CONSTRAINT "pk_product_model_product_model_id" PRIMARY KEY (product_model_id);
CLUSTER production.product_model USING "pk_product_model_product_model_id";

ALTER TABLE production.product_model_illustration ADD CONSTRAINT "pk_product_model_illustration_product_model_id_illustration_id" PRIMARY KEY (product_model_id, illustration_id);
CLUSTER production.product_model_illustration USING "pk_product_model_illustration_product_model_id_illustration_id";

ALTER TABLE production.product_model_product_description_culture ADD CONSTRAINT "pk_product_model_product_description_culture_product_model_id_product" PRIMARY KEY (product_model_id, product_description_id, culture_id);
CLUSTER production.product_model_product_description_culture USING "pk_product_model_product_description_culture_product_model_id_product";

ALTER TABLE production.product_photo ADD CONSTRAINT "pk_product_photo_product_photo_id" PRIMARY KEY (product_photo_id);
CLUSTER production.product_photo USING "pk_product_photo_product_photo_id";

ALTER TABLE production.product_product_photo ADD CONSTRAINT "pk_product_product_photo_product_id_product_photo_id" PRIMARY KEY (product_id, product_photo_id);

ALTER TABLE production.product_review ADD CONSTRAINT "pk_product_review_product_review_id" PRIMARY KEY (product_review_id);
CLUSTER production.product_review USING "pk_product_review_product_review_id";

ALTER TABLE production.product_subcategory ADD CONSTRAINT "pk_product_subcategory_product_subcategory_id" PRIMARY KEY (product_subcategory_id);
CLUSTER production.product_subcategory USING "pk_product_subcategory_product_subcategory_id";

ALTER TABLE purchasing.product_vendor ADD CONSTRAINT "pk_product_vendor_product_id_business_entity_id" PRIMARY KEY (product_id, business_entity_id);
CLUSTER purchasing.product_vendor USING "pk_product_vendor_product_id_business_entity_id";

ALTER TABLE purchasing.purchase_order_detail ADD CONSTRAINT "pk_purchase_order_detail_purchase_order_id_purchase_order_detail_id" PRIMARY KEY (purchase_order_id, purchase_order_detail_id);
CLUSTER purchasing.purchase_order_detail USING "pk_purchase_order_detail_purchase_order_id_purchase_order_detail_id";

ALTER TABLE purchasing.purchase_order_header ADD CONSTRAINT "pk_purchase_order_header_purchase_order_id" PRIMARY KEY (purchase_order_id);
CLUSTER purchasing.purchase_order_header USING "pk_purchase_order_header_purchase_order_id";

ALTER TABLE sales.person_credit_card ADD CONSTRAINT "pk_person_credit_card_business_entity_id_credit_card_id" PRIMARY KEY (business_entity_id, credit_card_id);
CLUSTER sales.person_credit_card USING "pk_person_credit_card_business_entity_id_credit_card_id";

ALTER TABLE sales.sales_order_detail ADD CONSTRAINT "pk_sales_order_detail_sales_order_id_sales_order_detail_id" PRIMARY KEY (sales_order_id, sales_order_detail_id);
CLUSTER sales.sales_order_detail USING "pk_sales_order_detail_sales_order_id_sales_order_detail_id";

ALTER TABLE sales.sales_order_header ADD CONSTRAINT "pk_sales_order_header_sales_order_id" PRIMARY KEY (sales_order_id);
CLUSTER sales.sales_order_header USING "pk_sales_order_header_sales_order_id";

ALTER TABLE sales.sales_order_header_sales_reason ADD CONSTRAINT "pk_sales_order_header_sales_reason_sales_order_id_sales_reason_id" PRIMARY KEY (sales_order_id, sales_reason_id);
CLUSTER sales.sales_order_header_sales_reason USING "pk_sales_order_header_sales_reason_sales_order_id_sales_reason_id";

ALTER TABLE sales.sales_person ADD CONSTRAINT "pk_sales_person_business_entity_id" PRIMARY KEY (business_entity_id);
CLUSTER sales.sales_person USING "pk_sales_person_business_entity_id";

ALTER TABLE sales.sales_person_quota_history ADD CONSTRAINT "pk_sales_person_quota_history_business_entity_id_quota_date" PRIMARY KEY (business_entity_id, quota_date); -- product_category_id);
CLUSTER sales.sales_person_quota_history USING "pk_sales_person_quota_history_business_entity_id_quota_date";

ALTER TABLE sales.sales_reason ADD CONSTRAINT "pk_sales_reason_sales_reason_id" PRIMARY KEY (sales_reason_id);
CLUSTER sales.sales_reason USING "pk_sales_reason_sales_reason_id";

ALTER TABLE sales.sales_tax_rate ADD CONSTRAINT "pk_sales_tax_rate_sales_tax_rate_id" PRIMARY KEY (sales_tax_rate_id);
CLUSTER sales.sales_tax_rate USING "pk_sales_tax_rate_sales_tax_rate_id";

ALTER TABLE sales.sales_territory ADD CONSTRAINT "pk_sales_territory_territory_id" PRIMARY KEY (territory_id);
CLUSTER sales.sales_territory USING "pk_sales_territory_territory_id";

ALTER TABLE sales.sales_territory_history ADD CONSTRAINT "pk_sales_territory_history_business_entity_id_start_date_territory_id" PRIMARY KEY (business_entity_id, start_date, territory_id);
CLUSTER sales.sales_territory_history USING "pk_sales_territory_history_business_entity_id_start_date_territory_id";

ALTER TABLE production.scrap_reason ADD CONSTRAINT "pk_scrap_reason_scrap_reason_id" PRIMARY KEY (scrap_reason_id);
CLUSTER production.scrap_reason USING "pk_scrap_reason_scrap_reason_id";

ALTER TABLE human_resources.shift ADD CONSTRAINT "pk_shift_shift_id" PRIMARY KEY (shift_id);
CLUSTER human_resources.shift USING "pk_shift_shift_id";

ALTER TABLE purchasing.ship_method ADD CONSTRAINT "pk_ship_method_ship_method_id" PRIMARY KEY (ship_method_id);
CLUSTER purchasing.ship_method USING "pk_ship_method_ship_method_id";

ALTER TABLE sales.shopping_cart_item ADD CONSTRAINT "pk_shopping_cart_item_shopping_cart_item_id" PRIMARY KEY (shopping_cart_item_id);
CLUSTER sales.shopping_cart_item USING "pk_shopping_cart_item_shopping_cart_item_id";

ALTER TABLE sales.special_offer ADD CONSTRAINT "pk_special_offer_special_offer_id" PRIMARY KEY (special_offer_id);
CLUSTER sales.special_offer USING "pk_special_offer_special_offer_id";

ALTER TABLE sales.special_offer_product ADD CONSTRAINT "pk_special_offer_product_special_offer_id_product_id" PRIMARY KEY (special_offer_id, product_id);
CLUSTER sales.special_offer_product USING "pk_special_offer_product_special_offer_id_product_id";

ALTER TABLE person.state_province ADD CONSTRAINT "pk_state_province_state_province_id" PRIMARY KEY (state_province_id);
CLUSTER person.state_province USING "pk_state_province_state_province_id";

ALTER TABLE sales.store ADD CONSTRAINT "pk_store_business_entity_id" PRIMARY KEY (business_entity_id);
CLUSTER sales.store USING "pk_store_business_entity_id";

ALTER TABLE production.transaction_history ADD CONSTRAINT "pk_transaction_history_transaction_id" PRIMARY KEY (transaction_id);
CLUSTER production.transaction_history USING "pk_transaction_history_transaction_id";

ALTER TABLE production.transaction_history_archive ADD CONSTRAINT "pk_transaction_history_archive_transaction_id" PRIMARY KEY (transaction_id);
CLUSTER production.transaction_history_archive USING "pk_transaction_history_archive_transaction_id";

ALTER TABLE production.unit_measure ADD CONSTRAINT "pk_unit_measure_unit_measure_code" PRIMARY KEY (unit_measure_code);
CLUSTER production.unit_measure USING "pk_unit_measure_unit_measure_code";

ALTER TABLE purchasing.vendor ADD CONSTRAINT "pk_vendor_business_entity_id" PRIMARY KEY (business_entity_id);
CLUSTER purchasing.vendor USING "pk_vendor_business_entity_id";

ALTER TABLE production.work_order ADD CONSTRAINT "pk_work_order_work_order_id" PRIMARY KEY (work_order_id);
CLUSTER production.work_order USING "pk_work_order_work_order_id";

ALTER TABLE production.work_order_routing ADD CONSTRAINT "pk_work_order_routing_work_order_id_product_id_operation_sequence" PRIMARY KEY (work_order_id, product_id, operation_sequence);
CLUSTER production.work_order_routing USING "pk_work_order_routing_work_order_id_product_id_operation_sequence";



-------------------------------------
-- FOREIGN KEYS
-------------------------------------

ALTER TABLE person.address ADD CONSTRAINT "fk_address_state_province_state_province_id" FOREIGN KEY (state_province_id) REFERENCES person.state_province(state_province_id);

ALTER TABLE production.bill_of_materials ADD CONSTRAINT "fk_bill_of_materials_product_product_assembly_id" FOREIGN KEY (product_assembly_id) REFERENCES production.product(product_id);
ALTER TABLE production.bill_of_materials ADD CONSTRAINT "fk_bill_of_materials_product_component_id" FOREIGN KEY (component_id) REFERENCES production.product(product_id);
ALTER TABLE production.bill_of_materials ADD CONSTRAINT "fk_bill_of_materials_unit_measure_unit_measure_code" FOREIGN KEY (unit_measure_code) REFERENCES production.unit_measure(unit_measure_code);

ALTER TABLE person.business_entity_address ADD CONSTRAINT "fk_business_entity_address_address_type_address_type_id" FOREIGN KEY (address_type_id) REFERENCES person.address_type(address_type_id);
ALTER TABLE person.business_entity_address ADD CONSTRAINT "fk_business_entity_address_business_entity_business_entity_id" FOREIGN KEY (business_entity_id) REFERENCES person.business_entity(business_entity_id);

ALTER TABLE person.business_entity_contact ADD CONSTRAINT "fk_business_entity_contact_person_person_id" FOREIGN KEY (person_id) REFERENCES person.person(business_entity_id);
ALTER TABLE person.business_entity_contact ADD CONSTRAINT "fk_business_entity_contact_contact_type_contact_type_id" FOREIGN KEY (contact_type_id) REFERENCES person.contact_type(contact_type_id);
ALTER TABLE person.business_entity_contact ADD CONSTRAINT "fk_business_entity_contact_business_entity_business_entity_id" FOREIGN KEY (business_entity_id) REFERENCES person.business_entity(business_entity_id);

ALTER TABLE sales.country_region_currency ADD CONSTRAINT "fk_country_region_currency_country_region_country_region_code" FOREIGN KEY (country_region_code) REFERENCES person.country_region(country_region_code);
ALTER TABLE sales.country_region_currency ADD CONSTRAINT "fk_country_region_currency_currency_currency_code" FOREIGN KEY (currency_code) REFERENCES sales.currency(currency_code);

ALTER TABLE sales.currency_rate ADD CONSTRAINT "fk_currency_rate_currency_from_currency_code" FOREIGN KEY (from_currency_code) REFERENCES sales.currency(currency_code);
ALTER TABLE sales.currency_rate ADD CONSTRAINT "fk_currency_rate_currency_to_currency_code" FOREIGN KEY (to_currency_code) REFERENCES sales.currency(currency_code);

ALTER TABLE sales.customer ADD CONSTRAINT "fk_customer_person_person_id" FOREIGN KEY (person_id) REFERENCES person.Person(business_entity_id);
ALTER TABLE sales.customer ADD CONSTRAINT "fk_customer_store_store_id" FOREIGN KEY (store_id) REFERENCES sales.store(business_entity_id);
ALTER TABLE sales.customer ADD CONSTRAINT "fk_customer_sales_territory_territory_id" FOREIGN KEY (territory_id) REFERENCES sales.sales_territory(territory_id);

ALTER TABLE production.document ADD CONSTRAINT "fk_document_employee_owner" FOREIGN KEY (owner) REFERENCES human_resources.employee(business_entity_id);

ALTER TABLE person.email_address ADD CONSTRAINT "fk_email_address_person_business_entity_id" FOREIGN KEY (business_entity_id) REFERENCES person.person(business_entity_id);

ALTER TABLE human_resources.employee ADD CONSTRAINT "fk_employee_person_business_entity_id" FOREIGN KEY (business_entity_id) REFERENCES person.person(business_entity_id);

ALTER TABLE human_resources.employee_department_history ADD CONSTRAINT "fk_employee_department_history_department_department_id" FOREIGN KEY (department_id) REFERENCES human_resources.department(department_id);
ALTER TABLE human_resources.employee_department_history ADD CONSTRAINT "fk_employee_department_history_employee_business_entity_id" FOREIGN KEY (business_entity_id) REFERENCES human_resources.employee(business_entity_id);
ALTER TABLE human_resources.employee_department_history ADD CONSTRAINT "fk_employee_department_history_shift_shift_id" FOREIGN KEY (shift_id) REFERENCES human_resources.shift(shift_id);

ALTER TABLE human_resources.employee_pay_history ADD CONSTRAINT "fk_employee_pay_history_employee_business_entity_id" FOREIGN KEY (business_entity_id) REFERENCES human_resources.employee(business_entity_id);

ALTER TABLE human_resources.job_candidate ADD CONSTRAINT "fk_job_candidate_employee_business_entity_id" FOREIGN KEY (business_entity_id) REFERENCES human_resources.employee(business_entity_id);

ALTER TABLE person.password ADD CONSTRAINT "fk_password_person_business_entity_id" FOREIGN KEY (business_entity_id) REFERENCES person.person(business_entity_id);

ALTER TABLE person.person ADD CONSTRAINT "fk_person_business_entity_business_entity_id" FOREIGN KEY (business_entity_id) REFERENCES person.business_entity(business_entity_id);

ALTER TABLE sales.person_credit_card ADD CONSTRAINT "fk_person_credit_card_person_business_entity_id" FOREIGN KEY (business_entity_id) REFERENCES person.person(business_entity_id);
ALTER TABLE sales.person_credit_card ADD CONSTRAINT "fk_person_credit_card_credit_card_credit_card_id" FOREIGN KEY (credit_card_id) REFERENCES sales.credit_card(credit_card_id);

ALTER TABLE person.person_phone ADD CONSTRAINT "fk_person_phone_person_business_entity_id" FOREIGN KEY (business_entity_id) REFERENCES person.Person(business_entity_id);
ALTER TABLE person.person_phone ADD CONSTRAINT "fk_person_phone_phone_number_type_phone_number_type_id" FOREIGN KEY (phone_number_type_id) REFERENCES person.phone_number_type(phone_number_type_id);

ALTER TABLE production.product ADD CONSTRAINT "fk_product_unit_measure_size_unit_measure_code" FOREIGN KEY (size_unit_measure_code) REFERENCES production.unit_measure(unit_measure_code);
ALTER TABLE production.product ADD CONSTRAINT "fk_product_unit_measure_weight_unit_measure_code" FOREIGN KEY (weight_unit_measure_code) REFERENCES production.unit_measure(unit_measure_code);
ALTER TABLE production.product ADD CONSTRAINT "fk_product_product_subcategory_product_subcategory_id" FOREIGN KEY (product_subcategory_id) REFERENCES production.product_subcategory(product_subcategory_id);

ALTER TABLE production.product_cost_history ADD CONSTRAINT "fk_product_cost_history_product_product_id" FOREIGN KEY (product_id) REFERENCES production.product(product_id);

ALTER TABLE production.product_document ADD CONSTRAINT "fk_product_document_product_product_id" FOREIGN KEY (product_id) REFERENCES production.product(product_id);
ALTER TABLE production.product_inventory ADD CONSTRAINT "fk_product_inventory_location_location_id" FOREIGN KEY (location_id) REFERENCES production.location(location_id);
ALTER TABLE production.product_inventory ADD CONSTRAINT "fk_product_inventory_product_product_id" FOREIGN KEY (product_id) REFERENCES production.product(product_id);

ALTER TABLE production.product_list_price_history ADD CONSTRAINT "fk_product_list_price_history_product_product_id" FOREIGN KEY (product_id) REFERENCES production.product(product_id);

ALTER TABLE production.product_model_illustration ADD CONSTRAINT "fk_product_model_illustration_illustration_illustration_id" FOREIGN KEY (illustration_id) REFERENCES production.illustration(illustration_id);

ALTER TABLE production.product_model_product_description_culture ADD CONSTRAINT "fk_product_model_product_description_culture_product_description_pro" FOREIGN KEY (product_description_id) REFERENCES production.product_description(product_description_id);
ALTER TABLE production.product_model_product_description_culture ADD CONSTRAINT "fk_product_model_product_description_culture_culture_culture_id" FOREIGN KEY (culture_id) REFERENCES production.culture(culture_id);

ALTER TABLE production.product_product_photo ADD CONSTRAINT "fk_product_product_photo_product_product_id" FOREIGN KEY (product_id) REFERENCES production.product(product_id);
ALTER TABLE production.product_product_photo ADD CONSTRAINT "fk_product_product_photo_product_photo_product_photo_id" FOREIGN KEY (product_photo_id) REFERENCES production.product_photo(product_photo_id);

ALTER TABLE production.product_review ADD CONSTRAINT "fk_product_review_product_product_id" FOREIGN KEY (product_id) REFERENCES production.product(product_id);

ALTER TABLE production.product_subcategory ADD CONSTRAINT "fk_product_subcategory_product_category_product_category_id" FOREIGN KEY (product_category_id) REFERENCES production.product_category(product_category_id);
ALTER TABLE purchasing.product_vendor ADD CONSTRAINT "fk_product_vendor_product_product_id" FOREIGN KEY (product_id) REFERENCES production.product(product_id);
ALTER TABLE purchasing.product_vendor ADD CONSTRAINT "fk_product_vendor_unit_measure_unit_measure_code" FOREIGN KEY (unit_measure_code) REFERENCES production.unit_measure(unit_measure_code);
ALTER TABLE purchasing.product_vendor ADD CONSTRAINT "fk_product_vendor_vendor_business_entity_id" FOREIGN KEY (business_entity_id) REFERENCES purchasing.vendor(business_entity_id);

ALTER TABLE purchasing.purchase_order_detail ADD CONSTRAINT "fk_purchase_order_detail_product_product_id" FOREIGN KEY (product_id) REFERENCES production.product(product_id);
ALTER TABLE purchasing.purchase_order_detail ADD CONSTRAINT "fk_purchase_order_detail_purchase_order_header_purchase_order_id" FOREIGN KEY (purchase_order_id) REFERENCES purchasing.purchase_order_header(purchase_order_id);

ALTER TABLE purchasing.purchase_order_header ADD CONSTRAINT "fk_purchase_order_header_employee_employee_id" FOREIGN KEY (employee_id) REFERENCES human_resources.employee(business_entity_id);
ALTER TABLE purchasing.purchase_order_header ADD CONSTRAINT "fk_purchase_order_header_vendor_vendor_id" FOREIGN KEY (vendor_id) REFERENCES purchasing.vendor(business_entity_id);
ALTER TABLE purchasing.purchase_order_header ADD CONSTRAINT "fk_purchase_order_header_ship_method_ship_method_id" FOREIGN KEY (ship_method_id) REFERENCES purchasing.ship_method(ship_method_id);

ALTER TABLE sales.sales_order_detail ADD CONSTRAINT "fk_sales_order_detail_sales_order_header_sales_order_id" FOREIGN KEY (sales_order_id) REFERENCES sales.sales_order_header(sales_order_id) ON DELETE CASCADE;
ALTER TABLE sales.sales_order_detail ADD CONSTRAINT "fk_sales_order_detail_special_offer_product_special_offer_id_product_id" FOREIGN KEY (special_offer_id, product_id) REFERENCES sales.special_offer_product(special_offer_id, product_id);

ALTER TABLE sales.sales_order_header ADD CONSTRAINT "fk_sales_order_header_address_bill_to_address_id" FOREIGN KEY (bill_to_address_id) REFERENCES person.Address(address_id);
ALTER TABLE sales.sales_order_header ADD CONSTRAINT "fk_sales_order_header_address_ship_to_address_id" FOREIGN KEY (ship_to_address_id) REFERENCES person.address(address_id);
ALTER TABLE sales.sales_order_header ADD CONSTRAINT "fk_sales_order_header_credit_card_credit_card_id" FOREIGN KEY (credit_card_id) REFERENCES sales.credit_card(credit_card_id);
ALTER TABLE sales.sales_order_header ADD CONSTRAINT "fk_sales_order_header_currency_rate_currency_rate_id" FOREIGN KEY (currency_rate_id) REFERENCES sales.currency_rate(currency_rate_id);
ALTER TABLE sales.sales_order_header ADD CONSTRAINT "fk_sales_order_header_customer_customer_id" FOREIGN KEY (customer_id) REFERENCES sales.customer(customer_id);
ALTER TABLE sales.sales_order_header ADD CONSTRAINT "fk_sales_order_header_sales_person_sales_person_id" FOREIGN KEY (sales_person_id) REFERENCES sales.sales_person(business_entity_id);
ALTER TABLE sales.sales_order_header ADD CONSTRAINT "fk_sales_order_header_ship_method_ship_method_id" FOREIGN KEY (ship_method_id) REFERENCES purchasing.ship_method(ship_method_id);
ALTER TABLE sales.sales_order_header ADD CONSTRAINT "fk_sales_order_header_sales_territory_territory_id" FOREIGN KEY (territory_id) REFERENCES sales.sales_territory(territory_id);

ALTER TABLE sales.sales_order_header_sales_reason ADD CONSTRAINT "fk_sales_order_header_sales_reason_sales_reason_sales_reason_id" FOREIGN KEY (sales_reason_id) REFERENCES sales.sales_reason(sales_reason_id);
ALTER TABLE sales.sales_order_header_sales_reason ADD CONSTRAINT "fk_sales_order_header_sales_reason_sales_order_header_sales_order_id" FOREIGN KEY (sales_order_id) REFERENCES sales.sales_order_header(sales_order_id) ON DELETE CASCADE;

ALTER TABLE sales.sales_person ADD CONSTRAINT "fk_sales_person_employee_business_entity_id" FOREIGN KEY (business_entity_id) REFERENCES human_resources.employee(business_entity_id);
ALTER TABLE sales.sales_person ADD CONSTRAINT "fk_sales_person_sales_territory_territory_id" FOREIGN KEY (territory_id) REFERENCES sales.sales_territory(territory_id);

ALTER TABLE sales.sales_person_quota_history ADD CONSTRAINT "fk_sales_person_quota_history_sales_person_business_entity_id" FOREIGN KEY (business_entity_id) REFERENCES sales.sales_person(business_entity_id);

ALTER TABLE sales.sales_tax_rate ADD CONSTRAINT "fk_sales_tax_rate_state_province_state_province_id" FOREIGN KEY (state_province_id) REFERENCES person.state_province(state_province_id);

ALTER TABLE sales.sales_territory ADD CONSTRAINT "fk_sales_territory_country_region_country_region_code" FOREIGN KEY (country_region_code) REFERENCES person.country_region(country_region_code);

ALTER TABLE sales.sales_territory_history ADD CONSTRAINT "fk_sales_territory_history_sales_person_business_entity_id" FOREIGN KEY (business_entity_id) REFERENCES sales.sales_person(business_entity_id);
ALTER TABLE sales.sales_territory_history ADD CONSTRAINT "fk_sales_territory_history_sales_territory_territory_id" FOREIGN KEY (territory_id) REFERENCES sales.sales_territory(territory_id);

ALTER TABLE sales.shopping_cart_item ADD CONSTRAINT "fk_shopping_cart_item_product_product_id" FOREIGN KEY (product_id) REFERENCES production.product(product_id);

ALTER TABLE sales.special_offer_product ADD CONSTRAINT "fk_special_offer_product_product_product_id" FOREIGN KEY (product_id) REFERENCES production.product(product_id);
ALTER TABLE sales.special_offer_product ADD CONSTRAINT "fk_special_offer_product_special_offer_special_offer_id" FOREIGN KEY (special_offer_id) REFERENCES sales.special_offer(special_offer_id);

ALTER TABLE person.state_province ADD CONSTRAINT "fk_state_province_country_region_country_region_code" FOREIGN KEY (country_region_code) REFERENCES person.country_region(country_region_code);
ALTER TABLE person.state_province ADD CONSTRAINT "fk_state_province_sales_territory_territory_id" FOREIGN KEY (territory_id) REFERENCES sales.sales_territory(territory_id);

ALTER TABLE sales.store ADD CONSTRAINT "fk_store_business_entity_business_entity_id" FOREIGN KEY (business_entity_id) REFERENCES person.business_entity(business_entity_id);
ALTER TABLE sales.store ADD CONSTRAINT "fk_store_sales_person_sales_person_id" FOREIGN KEY (sales_person_id) REFERENCES sales.sales_person(business_entity_id);

ALTER TABLE production.transaction_history ADD CONSTRAINT "fk_transaction_history_product_product_id" FOREIGN KEY (product_id) REFERENCES production.product(product_id);

ALTER TABLE purchasing.vendor ADD CONSTRAINT "fk_vendor_business_entity_business_entity_id" FOREIGN KEY (business_entity_id) REFERENCES person.business_entity(business_entity_id);

ALTER TABLE production.work_order ADD CONSTRAINT "fk_work_order_product_product_id" FOREIGN KEY (product_id) REFERENCES production.product(product_id);
ALTER TABLE production.work_order ADD CONSTRAINT "fk_work_order_scrap_reason_scrap_reason_id" FOREIGN KEY (scrap_reason_id) REFERENCES production.scrap_reason(scrap_reason_id);

ALTER TABLE production.work_order_routing ADD CONSTRAINT "fk_work_order_routing_location_location_id" FOREIGN KEY (location_id) REFERENCES production.Location(location_id);
ALTER TABLE production.work_order_routing ADD CONSTRAINT "fk_work_order_routing_work_order_work_order_id" FOREIGN KEY (work_order_id) REFERENCES production.work_order(work_order_id);


-------------------------------------
-- VIEWS
-------------------------------------

CREATE VIEW human_resources.v_employee
AS
SELECT
    e.business_entity_id
    ,p.title
    ,p.first_name
    ,p.middle_name
    ,p.last_name
    ,p.suffix
    ,e.job_title
    ,pp.phone_number
    ,pnt.name AS phone_number_type
    ,ea.email_address
    ,p.email_promotion
    ,a.address_line1
    ,a.address_line2
    ,a.city
    ,sp.name AS state_province_name
    ,a.postal_code
    ,cr.name AS country_region_name
    ,p.additional_contact_info
FROM human_resources.employee e
  INNER JOIN person.person p
  ON p.business_entity_id = e.business_entity_id
    INNER JOIN person.business_entity_address bea
    ON bea.business_entity_id = e.business_entity_id
    INNER JOIN person.address a
    ON a.address_id = bea.address_id
    INNER JOIN person.state_province sp
    ON sp.state_province_id = a.state_province_id
    INNER JOIN person.country_region cr
    ON cr.country_region_code = sp.country_region_code
    LEFT OUTER JOIN person.person_phone pp
    ON pp.business_entity_id = p.business_entity_id
    LEFT OUTER JOIN person.phone_number_type pnt
    ON pp.phone_number_type_id = pnt.phone_number_type_id
    LEFT OUTER JOIN person.email_address ea
    ON p.business_entity_id = ea.business_entity_id;

CREATE VIEW human_resources.v_employee_department
AS
SELECT
    e.business_entity_id
    ,p.title
    ,p.first_name
    ,p.middle_name
    ,p.last_name
    ,p.suffix
    ,e.job_title
    ,d.name AS department
    ,d.group_name
    ,edh.start_date
FROM human_resources.employee e
  INNER JOIN person.person p
  ON p.business_entity_id = e.business_entity_id
    INNER JOIN human_resources.employee_department_history edh
    ON e.business_entity_id = edh.business_entity_id
    INNER JOIN human_resources.department d
    ON edh.department_id = d.department_id
WHERE edh.end_date IS NULL;

CREATE VIEW human_resources.v_employee_department_history
AS
SELECT
    e.business_entity_id
    ,p.title
    ,p.first_name
    ,p.middle_name
    ,p.last_name
    ,p.suffix
    ,s.name AS shift
    ,d.name AS department
    ,d.group_name
    ,edh.start_date
    ,edh.end_date
FROM human_resources.employee e
  INNER JOIN person.person p
  ON p.business_entity_id = e.business_entity_id
    INNER JOIN human_resources.employee_department_history edh
    ON e.business_entity_id = edh.business_entity_id
    INNER JOIN human_resources.department d
    ON edh.department_id = d.department_id
    INNER JOIN human_resources.shift s
    ON s.shift_id = edh.shift_id;

CREATE VIEW sales.v_individual_customer
AS
SELECT
    p.business_entity_id
    ,p.title
    ,p.first_name
    ,p.middle_name
    ,p.last_name
    ,p.suffix
    ,pp.phone_number
    ,pnt.name AS phone_number_type
    ,ea.email_address
    ,p.email_promotion
    ,at.name AS address_type
    ,a.address_line1
    ,a.address_line2
    ,a.city
    ,sp.name AS state_province_name
    ,a.postal_code
    ,cr.name AS country_region_name
    ,p.demographics
FROM person.person p
    INNER JOIN person.business_entity_address bea
    ON bea.business_entity_id = p.business_entity_id
    INNER JOIN person.address a
    ON a.address_id = bea.address_id
    INNER JOIN person.state_province sp
    ON sp.state_province_id = a.state_province_id
    INNER JOIN person.country_region cr
    ON cr.country_region_code = sp.country_region_code
    INNER JOIN person.address_type at
    ON at.address_type_id = bea.address_type_id
  INNER JOIN sales.customer c
  ON c.person_id = p.business_entity_id
  LEFT OUTER JOIN person.email_address ea
  ON ea.business_entity_id = p.business_entity_id
  LEFT OUTER JOIN person.person_phone pp
  ON pp.business_entity_id = p.business_entity_id
  LEFT OUTER JOIN person.phone_number_type pnt
  ON pnt.phone_number_type_id = pp.phone_number_type_id
WHERE c.store_id IS NULL;

CREATE VIEW production.v_product_and_description
AS
SELECT
    p.product_id
    ,p.name
    ,pm.name AS product_model
    ,pmx.culture_id
    ,pd.description
FROM production.product p
    INNER JOIN production.product_model pm
    ON p.product_model_id = pm.product_model_id
    INNER JOIN production.product_model_product_description_culture pmx
    ON pm.product_model_id = pmx.product_model_id
    INNER JOIN production.product_description pd
    ON pmx.product_description_id = pd.product_description_id;

CREATE VIEW sales.v_sales_person
AS
SELECT
    s.business_entity_id
    ,p.title
    ,p.first_name
    ,p.middle_name
    ,p.last_name
    ,p.suffix
    ,e.job_title
    ,pp.phone_number
    ,pnt.name AS phone_number_type
    ,ea.email_address
    ,p.email_promotion
    ,a.address_line1
    ,a.address_line2
    ,a.city
    ,sp.name AS state_province_name
    ,a.postal_code
    ,cr.name AS country_region_name
    ,st.name AS territory_name
    ,st.group AS territory_group
    ,s.sales_quota
    ,s.sales_ytd
    ,s.sales_last_year
FROM sales.sales_person s
    INNER JOIN human_resources.employee e
    ON e.business_entity_id = s.business_entity_id
  INNER JOIN person.Person p
  ON p.business_entity_id = s.business_entity_id
    INNER JOIN person.business_entity_address bea
    ON bea.business_entity_id = s.business_entity_id
    INNER JOIN person.address a
    ON a.address_id = bea.address_id
    INNER JOIN person.state_province sp
    ON sp.state_province_id = a.state_province_id
    INNER JOIN person.country_region cr
    ON cr.country_region_code = sp.country_region_code
    LEFT OUTER JOIN sales.sales_territory st
    ON st.territory_id = s.territory_id
  LEFT OUTER JOIN person.email_address ea
  ON ea.business_entity_id = p.business_entity_id
  LEFT OUTER JOIN person.person_phone pp
  ON pp.business_entity_id = p.business_entity_id
  LEFT OUTER JOIN person.phone_number_type pnt
  ON pnt.phone_number_type_id = pp.phone_number_type_id;

CREATE VIEW person.v_state_province_country_region
AS
SELECT
    sp.state_province_id
    ,sp.state_province_code
    ,sp.is_only_state_province_flag
    ,sp.name AS state_province_name
    ,sp.territory_id
    ,cr.country_region_code
    ,cr.name AS country_region_name
FROM person.state_province sp
    INNER JOIN person.country_region cr
    ON sp.country_region_code = cr.country_region_code;

CREATE VIEW sales.v_store_with_contacts AS
SELECT
    s.business_entity_id
    ,s.name
    ,ct.name AS contact_type
    ,p.title
    ,p.first_name
    ,p.middle_name
    ,p.last_name
    ,p.suffix
    ,pp.phone_number
    ,pnt.name AS phone_number_type
    ,ea.email_address
    ,p.email_promotion
FROM sales.store s
    INNER JOIN person.business_entity_contact bec
    ON bec.business_entity_id = s.business_entity_id
  INNER JOIN person.contact_type ct
  ON ct.contact_type_id = bec.contact_type_id
  INNER JOIN person.person p
  ON p.business_entity_id = bec.person_id
  LEFT OUTER JOIN person.email_address ea
  ON ea.business_entity_id = p.business_entity_id
  LEFT OUTER JOIN person.person_phone pp
  ON pp.business_entity_id = p.business_entity_id
  LEFT OUTER JOIN person.phone_number_type pnt
  ON pnt.phone_number_type_id = pp.phone_number_type_id;

CREATE VIEW sales.v_store_with_addresses AS
SELECT
    s.business_entity_id
    ,s.name
    ,at.name AS address_type
    ,a.address_line1
    ,a.address_line2
    ,a.city
    ,sp.name AS state_province_name
    ,a.postal_code
    ,cr.name AS country_region_name
FROM sales.store s
    INNER JOIN person.business_entity_address bea
    ON bea.business_entity_id = s.business_entity_id
    INNER JOIN person.address a
    ON a.address_id = bea.address_id
    INNER JOIN person.state_province sp
    ON sp.state_province_id = a.state_province_id
    INNER JOIN person.country_region cr
    ON cr.country_region_code = sp.country_region_code
    INNER JOIN person.address_type at
    ON at.address_type_id = bea.address_type_id;

CREATE VIEW purchasing.v_vendor_with_contacts AS
SELECT
    v.business_entity_id
    ,v.name
    ,ct.name AS contact_type
    ,p.title
    ,p.first_name
    ,p.middle_name
    ,p.last_name
    ,p.suffix
    ,pp.phone_number
    ,pnt.name AS phone_number_type
    ,ea.email_address
    ,p.email_promotion
FROM purchasing.vendor v
    INNER JOIN person.business_entity_contact bec
    ON bec.business_entity_id = v.business_entity_id
  INNER JOIN person.contact_type ct
  ON ct.contact_type_id = bec.contact_type_id
  INNER JOIN person.person p
  ON p.business_entity_id = bec.person_id
  LEFT OUTER JOIN person.email_address ea
  ON ea.business_entity_id = p.business_entity_id
  LEFT OUTER JOIN person.person_phone pp
  ON pp.business_entity_id = p.business_entity_id
  LEFT OUTER JOIN person.phone_number_type pnt
  ON pnt.phone_number_type_id = pp.phone_number_type_id;

CREATE VIEW purchasing.v_vendor_with_addresses AS
SELECT
    v.business_entity_id
    ,v.name
    ,at.name AS address_type
    ,a.address_line1
    ,a.address_line2
    ,a.city
    ,sp.name AS state_province_name
    ,a.postal_code
    ,cr.name AS country_region_name
FROM purchasing.vendor v
    INNER JOIN person.business_entity_address bea
    ON bea.business_entity_id = v.business_entity_id
    INNER JOIN person.address a
    ON a.address_id = bea.address_id
    INNER JOIN person.state_province sp
    ON sp.state_province_id = a.state_province_id
    INNER JOIN person.country_region cr
    ON cr.country_region_code = sp.country_region_code
    INNER JOIN person.address_type at
    ON at.address_type_id = bea.address_type_id;

-- For indices order test
CREATE INDEX "idx_country_region_currency_currency_code" ON sales.country_region_currency (currency_code);
