CREATE TABLE public.user_activites (
    id bigint NOT NULL,
    activity_id bigint,
    type character varying NOT NULL,
    metadata json DEFAULT '"{}"'::json,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


CREATE TABLE public.agent_api (
    id bigint NOT NULL,
    base_url character varying NOT NULL,
    "position" character varying NOT NULL,
    impact public.business_impact_enum,
    status public.api_status_enum DEFAULT 'not_configured'::public.api_status_enum NOT NULL,
    errored character varying DEFAULT ''::character varying NOT NULL,
    api_ids integer[] DEFAULT '{}'::integer[],
    custom_compliance_standards character varying[] DEFAULT '{}'::character varying[]
);

CREATE TABLE public.api_scanned_repo (
    id bigint NOT NULL,
    api_set_id integer NOT NULL,
    repository_id integer NOT NULL,
    comment text,
    created_at timestamp(6) without time zone NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL
);


CREATE TABLE IF NOT EXISTS customers
(
    id int auto_increment primary key,
    update_time timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    is_active tinyint(1) unsigned default 1 not null,
    create_time timestamp default CURRENT_TIMESTAMP not null,
    login varchar(128) default '' not null,
    `key` varchar(50) not NULL
);

CREATE TABLE IF NOT EXISTS `tblsample` (
  `id` int(11) NOT NULL auto_increment,
  `recid` int(11) NOT NULL default '0',
  `cvfilename` varchar(250)  NOT NULL default '',
  `data` varchar(100) NOT NULL default '',
   PRIMARY KEY  (`id`)
);