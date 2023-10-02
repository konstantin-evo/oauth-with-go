--
-- Name: oauth_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--
CREATE SEQUENCE public.oauth_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE CACHE 1;

ALTER TABLE public.oauth_id_seq OWNER TO postgres;

SET
default_tablespace = '';
SET
default_table_access_method = heap;

--
-- Name: oauth; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.oauth
(
    id            integer DEFAULT nextval('public.oauth_id_seq'::regclass) NOT NULL,
    access_token  text,
    token_type    text,
    expires_in    integer,
    refresh_token text,
    scope         text
);

ALTER TABLE public.oauth OWNER TO postgres;

--
-- Name: oauth_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.oauth_id_seq', 1, true);

--
-- Name: oauth_oauth_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.oauth
    ADD CONSTRAINT oauth_pkey PRIMARY KEY (id);
