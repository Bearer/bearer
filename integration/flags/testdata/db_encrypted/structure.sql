CREATE TABLE public.Users (
  id bigint NOT NULL,
  first_name character varying NOT NULL,
  last_name character varying NOT NULL,
  email character varying NOT NULL,
  tanker_encrypted_date_of_birth character varying NOT NULL,
  city character varying NOT NULL,
  country character varying NOT NULL,
  encrypted_gender character varying,
);
