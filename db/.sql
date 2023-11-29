-- public.tb01 definition

-- Drop table

-- DROP TABLE public.tb01;

CREATE TABLE public.tb01 (
	id serial4 NOT NULL,
	col_texto text NOT NULL,
	col_dt timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT tb01_pkey PRIMARY KEY (id)
);