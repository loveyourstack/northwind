
CREATE TABLE hr.employee
(	
  id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  country_fk bigint NOT NULL REFERENCES common.country(id),
  reports_to_fk bigint REFERENCES hr.employee(id), -- self-reference - is made mandatory after data is inserted
  date_of_birth date NOT NULL CHECK (date_of_birth > '1920-01-01'),
  hire_date date NOT NULL CHECK (hire_date > '2010-01-01'),
  entry_at tracking_at,
  last_modified_at tracking_at,
  address text NOT NULL,
  city text NOT NULL,
  entry_by tracking_by,
  first_name text NOT NULL,
  job_title text NOT NULL,
  home_phone text NOT NULL,
  last_modified_by tracking_by,
  last_name text NOT NULL,
  name text NOT NULL UNIQUE, -- not generated, in order to allow user to specify uniqueness however he chooses, e.g. using department or some other criterion
  notes text NOT NULL,
  postal_code text NOT NULL,
  state text NOT NULL,
  title text NOT NULL
);
COMMENT ON TABLE hr.employee IS 'shortname: hr_e';
