
CREATE TABLE public.users (
    id bigint NOT NULL,
    class_id bigint,
    type character varying NOT NULL,
    metadata json DEFAULT '"{}"'::json,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);