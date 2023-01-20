CREATE TABLE public.users (
  id bigint NOT NULL,
  name character varying,
  password character varying,
  created_at timestamp(6) without time zone NOT NULL,
  updated_at timestamp(6) without time zone NOT NULL,
  email character varying DEFAULT ''::character varying NOT NULL
);
