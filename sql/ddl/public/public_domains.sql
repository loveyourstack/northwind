
-- general use
--------------------------------------------------

CREATE DOMAIN tracking_at AS timestamp with time zone NOT NULL DEFAULT now();
CREATE DOMAIN tracking_by AS text NOT NULL DEFAULT 'Unknown';

-- stati, codes, shortnames
CREATE DOMAIN text_short AS varchar(64) NOT NULL; 
CREATE DOMAIN text_short_mandatory AS varchar(64) NOT NULL CHECK (value != '');

-- person names, address lines
CREATE DOMAIN text_medium AS varchar(255) NOT NULL;
CREATE DOMAIN text_medium_mandatory AS varchar(255) NOT NULL CHECK (value != '');

CREATE DOMAIN int_gte0 AS integer NOT NULL CHECK (value >= 0);
CREATE DOMAIN int_positive AS integer NOT NULL CHECK (value > 0);

-- nw specific
--------------------------------------------------

CREATE DOMAIN price_gte0 AS numeric(12,2) NOT NULL CHECK (value >= 0.0);