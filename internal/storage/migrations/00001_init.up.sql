BEGIN TRANSACTION;

CREATE TABLE public.metric (
	id varchar NOT NULL,
	mtype int2 NOT NULL,
	delta int8 NULL,
	value float8 NULL
);


END TRANSACTION;