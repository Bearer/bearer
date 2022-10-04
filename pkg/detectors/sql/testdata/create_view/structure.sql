

CREATE VIEW s3_export_data.users AS
 SELECT users.id,
    users.created_at,
    users.updated_at,
    users.name,
    array_to_json(users.origin) AS users,
    users.country
   FROM public.users users;



CREATE VIEW s3_export_data.user_accounts AS
 SELECT doctor_accounts.id,
    user_accounts.email,
    user_accounts.last_name,
    user_accounts.first_name,
    user_accounts.phone_number,
    md5((user_accounts.email)::text) AS hashed_email,
    "substring"((user_accounts.email)::text, ("position"((user_accounts.email)::text, '@'::text) + 1), (length((user_accounts.email)::text) - "position"((user_accounts.email)::text, '@'::text))) AS email_domain
   FROM public.user_accounts user_accounts;



CREATE VIEW s3_export_data.users_av_version AS
 SELECT ua.id,
    ua.item_type,
    ua.item_id,
    aa.created_at,
    aa.user_agent,
   FROM (public.users_accounts ua
     JOIN public.admin_accounts aa ON ((ua.item_id = aa.id)));


CREATE VIEW s3_export_data.master_users AS
SELECT
    NULL::boolean AS has_own_phone_number,
    NULL::boolean AS has_main_with_different_phone_number,
    NULL::boolean AS has_own_email,


CREATE VIEW s3_export_data.notifications AS
 SELECT notifications.id,
    notifications.account_id,
    notifications.anonymized_at,
    ((public.yaml_to_json(notifications.meta_data) ->> 'current_start_date'::text))::timestamp without time zone AS current_start_date,
   FROM public.notifications;


CREATE VIEW s3_export_data.users_contact AS
 SELECT users.id,
    users.created_at,
        CASE
            WHEN ((COALESCE(users.email, ''::character varying))::text <> ''::text) THEN true
            ELSE false
        END AS has_email,
        CASE
            WHEN ((COALESCE(users.phone_number, ''::character varying))::text <> ''::text) THEN true
            ELSE false
        END AS has_phone_number,
        CASE
            WHEN ((COALESCE(users.secondary_phone_number, ''::character varying))::text <> ''::text) THEN true
            ELSE false
        END AS has_secondary_phone_number,
   FROM ((public.users
     LEFT JOIN public.users_security ON ((users.public_id = users_security.base_id)))
     LEFT JOIN public.admins ON ((admins.user_id = users.id)));


CREATE VIEW s3_export_data.admins_tables AS
 SELECT now() AS inserted_at,
    a.table_schema,
    a.table_name,
    (a.row_estimate)::bigint AS row_estimate,
    a.toast_bytes,
    ((a.total_bytes - a.index_bytes) - COALESCE(a.toast_bytes, (0)::bigint)) AS table_bytes
   FROM ( SELECT c.oid AS pg_oid,
            n.nspname AS table_schema,
            c.relname AS table_name,
           FROM (pg_class c
             LEFT JOIN pg_namespace n ON ((n.oid = c.relnamespace)))
          WHERE (c.relkind = 'r'::"char")) a;

