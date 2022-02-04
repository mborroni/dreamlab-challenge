CREATE TABLE  IF NOT EXISTS ip2location_px7
(
    ip_from      bigint                 NOT NULL,
    ip_to        bigint                 NOT NULL,
    proxy_type   character varying(3)   NOT NULL,
    country_code character(2)           NOT NULL,
    country_name character varying(64)  NOT NULL,
    region_name  character varying(128) NOT NULL,
    city_name    character varying(128) NOT NULL,
    isp          character varying(256) NOT NULL,
    domain       character varying(128) NOT NULL,
    usage_type   character varying(11)  NOT NULL,
    asn          character varying(10)  NOT NULL,
    "as"         character varying(256) NOT NULL,
    CONSTRAINT ip2location_db1_pkey PRIMARY KEY (ip_from, ip_to)
);