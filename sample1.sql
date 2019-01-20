--
-- PostgreSQL database dump
--

-- Dumped from database version 10.6 (Ubuntu 10.6-1.pgdg18.04+1)
-- Dumped by pg_dump version 11.1 (Ubuntu 11.1-1.pgdg18.04+1)

-- Started on 2019-01-20 12:06:32 EET

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 4 (class 2615 OID 16386)
-- Name: test; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA test;


ALTER SCHEMA test OWNER TO postgres;

--
-- TOC entry 206 (class 1255 OID 16419)
-- Name: comment_del(integer); Type: FUNCTION; Schema: test; Owner: postgres
--

CREATE FUNCTION test.comment_del(_id integer) RETURNS json
    LANGUAGE plpgsql
    AS $$
begin
  delete from test.comments where id = _id;

  return json_build_object('id', _id);

  exception when others then

  raise notice 'Illegal operation: %', SQLERRM;

  return json_build_object('error', SQLERRM);
end
$$;


ALTER FUNCTION test.comment_del(_id integer) OWNER TO postgres;

--
-- TOC entry 219 (class 1255 OID 16416)
-- Name: comment_get(integer); Type: FUNCTION; Schema: test; Owner: postgres
--

CREATE FUNCTION test.comment_get(_id integer) RETURNS json
    LANGUAGE plpgsql
    AS $$
declare
  _ret json;
begin
  if _id = 0 then
    select array_to_json(array(
      select row_to_json(r)
      from (
        select c.id, c.id_user, c.txt
        from test.comments c
      ) r
    )) into _ret;
  else
    select row_to_json(r) into _ret
    from (
      select c.id, c.id_user, c.txt
      from test.comments c
      where id = _id
    ) r;
  end if;

  return _ret;

  exception when others then

  return json_build_object('error', SQLERRM);
end
$$;


ALTER FUNCTION test.comment_get(_id integer) OWNER TO postgres;

--
-- TOC entry 205 (class 1255 OID 16418)
-- Name: comment_upd(integer, json); Type: FUNCTION; Schema: test; Owner: postgres
--

CREATE FUNCTION test.comment_upd(_id integer, _params json) RETURNS json
    LANGUAGE plpgsql
    AS $$
begin
  update test.comments set
    txt = _params->>'txt'
  where id = _id;

  return json_build_object('id', _id);

  exception when others then

  return json_build_object('error', SQLERRM);
end
$$;


ALTER FUNCTION test.comment_upd(_id integer, _params json) OWNER TO postgres;

--
-- TOC entry 221 (class 1255 OID 16438)
-- Name: user_comment_get(integer, integer); Type: FUNCTION; Schema: test; Owner: postgres
--

CREATE FUNCTION test.user_comment_get(_id_user integer, _id integer) RETURNS json
    LANGUAGE plpgsql
    AS $$
declare
  _ret json;
begin
if _id = 0 and _id_user=0 then
    select array_to_json(array(
      select row_to_json(r)
      from (
        select c.id, c.id_user, c.txt
        from test.comments c
      ) r
    )) into _ret; 
elsif _id = 0 then
    select array_to_json(array(
      select row_to_json(r)
      from (
        select c.id, c.id_user, c.txt
        from test.comments c
 where id_user=_id_user
      ) r
    )) into _ret;
elsif _id_user=0 then
    select row_to_json(r) into _ret
    from (
      select c.id, c.id_user, c.txt
        from test.comments c
      where id = _id
    ) r;
else 
    select row_to_json(r) into _ret
    from (
      select c.id, c.id_user, c.txt
        from test.comments c
      where id = _id and id_user=_id_user
    ) r;
  end if;

  return _ret;

  exception when others then

  return json_build_object('error', SQLERRM);
end
$$;


ALTER FUNCTION test.user_comment_get(_id_user integer, _id integer) OWNER TO postgres;

--
-- TOC entry 220 (class 1255 OID 16428)
-- Name: user_comment_ins(integer, json); Type: FUNCTION; Schema: test; Owner: postgres
--

CREATE FUNCTION test.user_comment_ins(_id_user integer, _params json) RETURNS json
    LANGUAGE plpgsql
    AS $$
declare
  _newid integer;
begin
  _newid = 0;

  insert into test.comments (id_user,txt)
  values (_id_user,'initializating comment');
  update test.comments set
  txt=_params->>'txt' 
  where id=(select max(id) from test.comments)
  returning id into _newid;

  return json_build_object('id', _newid);

  exception when others then

  return json_build_object('error', SQLERRM);
end
$$;


ALTER FUNCTION test.user_comment_ins(_id_user integer, _params json) OWNER TO postgres;

--
-- TOC entry 204 (class 1255 OID 16415)
-- Name: user_del(integer); Type: FUNCTION; Schema: test; Owner: postgres
--

CREATE FUNCTION test.user_del(_id integer) RETURNS json
    LANGUAGE plpgsql
    AS $$
begin
  delete from test.users where id = _id;

  return json_build_object('id', _id);

  exception when others then

  raise notice 'Illegal operation: %', SQLERRM;

  return json_build_object('error', SQLERRM);
end
$$;


ALTER FUNCTION test.user_del(_id integer) OWNER TO postgres;

--
-- TOC entry 201 (class 1255 OID 16412)
-- Name: user_get(integer); Type: FUNCTION; Schema: test; Owner: postgres
--

CREATE FUNCTION test.user_get(_id integer) RETURNS json
    LANGUAGE plpgsql
    AS $$
declare
  _ret json;
begin
  if _id = 0 then
    select array_to_json(array(
      select row_to_json(r)
      from (
        select u.id, u.name, u.email
        from test.users u
      ) r
    )) into _ret;
  else
    select row_to_json(r) into _ret
    from (
      select u.id, u.name, u.email
      from test.users u
      where id = _id
    ) r;
  end if;

  return _ret;

  exception when others then

  return json_build_object('error', SQLERRM);
end
$$;


ALTER FUNCTION test.user_get(_id integer) OWNER TO postgres;

--
-- TOC entry 202 (class 1255 OID 16413)
-- Name: user_ins(json); Type: FUNCTION; Schema: test; Owner: postgres
--

CREATE FUNCTION test.user_ins(_params json) RETURNS json
    LANGUAGE plpgsql
    AS $$
declare
  _newid integer;
begin
  _newid = 0;

  insert into test.users (name, email)
  select name, email
  from json_populate_record(null::test.users, _params)
  returning id into _newid;

  return json_build_object('id', _newid);

  exception when others then

  return json_build_object('error', SQLERRM);
end
$$;


ALTER FUNCTION test.user_ins(_params json) OWNER TO postgres;

--
-- TOC entry 203 (class 1255 OID 16414)
-- Name: user_upd(integer, json); Type: FUNCTION; Schema: test; Owner: postgres
--

CREATE FUNCTION test.user_upd(_id integer, _params json) RETURNS json
    LANGUAGE plpgsql
    AS $$
begin
  update test.users set
    name = _params->>'name',
    email = _params->>'email'
  where id = _id;

  return json_build_object('id', _id);

  exception when others then

  return json_build_object('error', SQLERRM);
end
$$;


ALTER FUNCTION test.user_upd(_id integer, _params json) OWNER TO postgres;

--
-- TOC entry 198 (class 1259 OID 16389)
-- Name: seq_comments; Type: SEQUENCE; Schema: test; Owner: postgres
--

CREATE SEQUENCE test.seq_comments
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE test.seq_comments OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- TOC entry 200 (class 1259 OID 16403)
-- Name: comments; Type: TABLE; Schema: test; Owner: postgres
--

CREATE TABLE test.comments (
    id integer DEFAULT nextval('test.seq_comments'::regclass) NOT NULL,
    id_user integer NOT NULL,
    txt character varying NOT NULL
);


ALTER TABLE test.comments OWNER TO postgres;

--
-- TOC entry 197 (class 1259 OID 16387)
-- Name: seq_users; Type: SEQUENCE; Schema: test; Owner: postgres
--

CREATE SEQUENCE test.seq_users
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE test.seq_users OWNER TO postgres;

--
-- TOC entry 199 (class 1259 OID 16391)
-- Name: users; Type: TABLE; Schema: test; Owner: postgres
--

CREATE TABLE test.users (
    id integer DEFAULT nextval('test.seq_users'::regclass) NOT NULL,
    name character varying NOT NULL,
    email character varying NOT NULL,
    CONSTRAINT "CHK_users_email" CHECK (((email)::text ~~ '%@%'::text))
);


ALTER TABLE test.users OWNER TO postgres;

--
-- TOC entry 2939 (class 0 OID 16403)
-- Dependencies: 200
-- Data for Name: comments; Type: TABLE DATA; Schema: test; Owner: postgres
--

COPY test.comments (id, id_user, txt) FROM stdin;
4	1	my 1 comment
5	2	my  best comment
8	3	IMHO
9	4	nice pic
10	1	nice info
11	1	my 1 comment
13	4	i am the goat
14	5	i am the gomy favorite flowers are roses
16	5	i am the goshit
17	3	+1000
19	3	beast
20	6	loshad
21	6	puncher
22	7	give me  stars
23	2	what about FAQ?
24	5	new comment from the GOAT
27	5	cheers
28	5	new comment
29	5	cheers
30	1	my comment will be the last
15	5	deadlock?
31	4	111
\.


--
-- TOC entry 2938 (class 0 OID 16391)
-- Dependencies: 199
-- Data for Name: users; Type: TABLE DATA; Schema: test; Owner: postgres
--

COPY test.users (id, name, email) FROM stdin;
2	Potap	potya@mail.ru
3	John	johny@mail.ru
4	Huanito	Huan@mail.es
5	Hommer78	Hom78@mail.es
1	Vlad	vladimir999@mail.ru
7	Avery	fdd@mail.ru
9	Noname	misterx@mail.ru
10	Irina	missis@mail.ru
11	Zevs	Pharaon@mail.ru
\.


--
-- TOC entry 2945 (class 0 OID 0)
-- Dependencies: 198
-- Name: seq_comments; Type: SEQUENCE SET; Schema: test; Owner: postgres
--

SELECT pg_catalog.setval('test.seq_comments', 31, true);


--
-- TOC entry 2946 (class 0 OID 0)
-- Dependencies: 197
-- Name: seq_users; Type: SEQUENCE SET; Schema: test; Owner: postgres
--

SELECT pg_catalog.setval('test.seq_users', 11, true);


--
-- TOC entry 2814 (class 2606 OID 16411)
-- Name: comments PK_comments; Type: CONSTRAINT; Schema: test; Owner: postgres
--

ALTER TABLE ONLY test.comments
    ADD CONSTRAINT "PK_comments" PRIMARY KEY (id);


--
-- TOC entry 2810 (class 2606 OID 16400)
-- Name: users PK_users; Type: CONSTRAINT; Schema: test; Owner: postgres
--

ALTER TABLE ONLY test.users
    ADD CONSTRAINT "PK_users" PRIMARY KEY (id);


--
-- TOC entry 2812 (class 2606 OID 16402)
-- Name: users UQ_users_email; Type: CONSTRAINT; Schema: test; Owner: postgres
--

ALTER TABLE ONLY test.users
    ADD CONSTRAINT "UQ_users_email" UNIQUE (email);


-- Completed on 2019-01-20 12:06:32 EET

--
-- PostgreSQL database dump complete
--

