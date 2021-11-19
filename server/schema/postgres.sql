CREATE TABLE IF NOT EXISTS organizations (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	name text NOT NULL,
	key bytea NOT NULL,
	timezone text NOT NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE ON UPDATE CASCADE,
	email text NOT NULL,
	password bytea NOT NULL,
	salt bytea NOT NULL,
	name text NOT NULL,
	role text NOT NULL,
	logged_out_at timestamp,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now(),
	UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS actions (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE ON UPDATE CASCADE,
	name text NOT NULL,
	type text NOT NULL,
	tags jsonb NOT NULL DEFAULT '[]'::jsonb,
	state text NOT NULL DEFAULT 'inactive',
	properties jsonb NOT NULL DEFAULT '[]'::jsonb,
	targeting jsonb NOT NULL DEFAULT '{}'::jsonb,
	capping jsonb NOT NULL DEFAULT '{}'::jsonb,
	hooks jsonb NOT NULL DEFAULT '[]'::jsonb,
	test_users jsonb NOT NULL DEFAULT '[]'::jsonb,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS profiles (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE ON UPDATE CASCADE,
	custom_traits jsonb NOT NULL DEFAULT '{}'::jsonb,
	computed_traits jsonb NOT NULL DEFAULT '{}'::jsonb,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS profilex_x ON profiles (organization_id, created_at DESC);

CREATE TABLE IF NOT EXISTS profile_identities (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE ON UPDATE CASCADE,
	profile_id uuid NOT NULL REFERENCES profiles(id) ON DELETE CASCADE ON UPDATE CASCADE,
	channel text NOT NULL,
	type text NOT NULL,
	user_id text NOT NULL,
	is_anonymous bool NOT NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now(),
	UNIQUE (organization_id, channel, user_id)
);

CREATE TABLE IF NOT EXISTS locked_profile_identities (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE ON UPDATE CASCADE,
	user_id text NOT NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	UNIQUE (organization_id, user_id)
);

CREATE TABLE IF NOT EXISTS events (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE ON UPDATE CASCADE,
	name text NOT NULL,
	description text NOT NULL,
	version int NOT NULL DEFAULT 0,
	structure jsonb NOT NULL DEFAULT '{}'::jsonb,
	json_schema jsonb NOT NULL DEFAULT '{}'::jsonb,
	properties jsonb NOT NULL DEFAULT '[]'::jsonb,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS audiences (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE ON UPDATE CASCADE,
	name text NOT NULL,
	description text NOT NULL,
	include_anonymous bool NOT NULL DEFAULT true,
	events jsonb NOT NULL DEFAULT '{}'::jsonb,
	traits jsonb NOT NULL DEFAULT '{}'::jsonb,
	current_set_id uuid,
	profile_count int NOT NULL DEFAULT 0,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now(),
	refreshed_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS audience_profiles (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE ON UPDATE CASCADE,
	audience_id uuid NOT NULL REFERENCES audiences(id) ON DELETE CASCADE ON UPDATE CASCADE,
	profile_id uuid NOT NULL REFERENCES profiles(id) ON DELETE CASCADE ON UPDATE CASCADE,
	set_id uuid NOT NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	UNIQUE (organization_id, audience_id, profile_id, set_id)
);
CREATE INDEX IF NOT EXISTS audience_profiles_x ON audience_profiles (organization_id, profile_id);

CREATE TABLE IF NOT EXISTS channels (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE ON UPDATE CASCADE,
	name text NOT NULL,
	key  text NOT NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now(),
	UNIQUE (organization_id, key)
);

CREATE TABLE IF NOT EXISTS auth_tokens (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE ON UPDATE CASCADE,
	token text NOT NULL,
	channel text NOT NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	UNIQUE (token)
);

CREATE TABLE IF NOT EXISTS computed_traits (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE ON UPDATE CASCADE,
	name text NOT NULL,
	key text NOT NULL,
	type text NOT NULL,
	event_id uuid NOT NULL,
	conditions jsonb NOT NULL DEFAULT '{}'::jsonb,
	properties jsonb NOT NULL DEFAULT '{}'::jsonb,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now(),
	refreshed_at timestamp NOT NULL DEFAULT now(),
	UNIQUE (organization_id, key)
);
CREATE INDEX IF NOT EXISTS organization_id_event_id ON computed_traits (organization_id, event_id);

CREATE TABLE IF NOT EXISTS custom_traits (
	organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE ON UPDATE CASCADE,
	structure jsonb NOT NULL DEFAULT '{}'::jsonb,
	json_schema jsonb NOT NULL DEFAULT '{}'::jsonb,
	properties jsonb NOT NULL DEFAULT '[]'::jsonb,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now(),
	UNIQUE (organization_id)
);

CREATE TABLE IF NOT EXISTS contexts (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE ON UPDATE CASCADE,
	channel text NOT NULL,
	structure jsonb NOT NULL DEFAULT '{}'::jsonb,
	json_schema jsonb NOT NULL DEFAULT '{}'::jsonb,
	properties jsonb NOT NULL DEFAULT '[]'::jsonb,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now(),
	UNIQUE (organization_id, channel)
);

-- Views

CREATE VIEW full_identities (id, profile_id, organization_id, channel, type, user_id, is_anonymous, custom_traits, computed_traits, created_at) AS
SELECT pi.id, pi.profile_id, pi.organization_id, pi.channel, pi.type, pi.user_id, pi.is_anonymous, p.custom_traits, p.computed_traits, pi.created_at
FROM profile_identities pi
INNER JOIN profiles p
ON p.id = pi.profile_id


--
-- PostgreSQL database dump
--

-- Dumped from database version 14.0 (Debian 14.0-1.pgdg110+1)
-- Dumped by pg_dump version 14.0 (Debian 14.0-1.pgdg110+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: organizations; Type: TABLE DATA; Schema: public; Owner: attractify
--

INSERT INTO public.organizations (id, name, key, timezone, created_at, updated_at) VALUES ('5647e66f-ef2c-4f29-9584-6f788f573e94', '', '\x802f027b6e811a8881c3ccbf8067bd5ce7aeefc5b23fcab7114d5e705c9ed322', 'Europe/Berlin', '2021-10-25 06:50:05.110673', '2021-10-25 06:50:05.110673');


--
-- Data for Name: actions; Type: TABLE DATA; Schema: public; Owner: attractify
--

INSERT INTO public.actions (id, organization_id, name, type, tags, state, properties, targeting, capping, hooks, test_users, created_at, updated_at) VALUES ('a1c70527-dc68-4072-a231-7bdf257847ee', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Not Available', 'not_available', '["improvement"]', 'active', '[{"name": "title", "type": "text", "value": "Notify me when available", "channels": ["web", "store", "app"], "sourceKey": "", "sourceType": ""}, {"name": "body", "type": "text", "value": "Sorry, the selected product is currently unavailable, but it sure will be back in stock soon. We''ll then notify you via e-mail.", "channels": ["web", "store", "app"], "sourceKey": "", "sourceType": ""}, {"name": "cta", "type": "text", "value": "Notification received", "channels": ["web", "store", "app"], "sourceKey": "", "sourceType": ""}]', '{"end": {"date": null, "time": null}, "start": {"date": null, "time": null}, "channels": ["web", "store", "app"], "audiences": null, "traitConditions": [], "contextConditions": []}', '[{"count": 2, "event": "hidden", "group": "user", "within": 0, "channels": ["web", "store", "app"]}, {"count": 1, "event": "accepted", "group": "user", "within": 0, "channels": ["store", "web", "app"]}]', '[]', '[]', '2021-10-27 09:04:28.45338', '2021-10-27 09:04:28.45338');
INSERT INTO public.actions (id, organization_id, name, type, tags, state, properties, targeting, capping, hooks, test_users, created_at, updated_at) VALUES ('b8184291-6177-4158-8f6f-b8e6d54482d6', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Incentivize Abandoned Cart Users', 'coupon_overlay', '["upselling"]', 'active', '[{"name": "title", "type": "text", "value": "Complete Your order now !!!", "channels": ["web"], "sourceKey": "", "sourceType": ""}, {"name": "body", "type": "text", "value": "We thought that maybe this 10% discount coupon will help You, to complete Your order? ", "channels": ["web"], "sourceKey": "", "sourceType": ""}, {"name": "cta", "type": "text", "value": "Redeem coupon now", "channels": ["web"], "sourceKey": "", "sourceType": ""}, {"name": "couponCode", "type": "text", "value": "B4ZXQZVX", "channels": ["web"], "sourceKey": "", "sourceType": ""}]', '{"end": {"date": null, "time": null}, "start": {"date": null, "time": null}, "channels": ["web"], "audiences": ["f443213f-cc5a-45b8-9662-49e692c745a3"], "traitConditions": [], "contextConditions": []}', '[{"count": 1, "event": "accepted", "group": "user", "within": 0, "channels": ["web"]}, {"count": 2, "event": "hidden", "group": "user", "within": 0, "channels": ["web"]}]', '[]', '[]', '2021-10-27 09:11:36.505495', '2021-10-27 09:11:36.505495');
INSERT INTO public.actions (id, organization_id, name, type, tags, state, properties, targeting, capping, hooks, test_users, created_at, updated_at) VALUES ('6ced6975-defe-4661-9634-017cc9564e83', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Preselect Size', 'size_preset', '["improvement"]', 'active', '[{"name": "size", "type": "computed_trait", "value": "", "channels": ["web"], "sourceKey": "\tfavoriteSize", "sourceType": ""}]', '{"end": {"date": null, "time": null}, "start": {"date": null, "time": null}, "channels": ["web"], "audiences": null, "traitConditions": [{"key": "\tfavoriteSize", "type": "string", "value": null, "source": "computed", "operator": "exists"}], "contextConditions": [{"key": "ua", "type": "string", "value": "123", "channel": "store", "operator": "contains"}]}', '[]', '[]', '[]', '2021-10-27 09:16:00.725879', '2021-10-27 09:16:00.725879');
INSERT INTO public.actions (id, organization_id, name, type, tags, state, properties, targeting, capping, hooks, test_users, created_at, updated_at) VALUES ('869d9397-43b5-4251-9dd3-c4b4cb66ecf6', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Preselect category', 'favorite_category', '["improvement"]', 'active', '[{"name": "category", "type": "computed_trait", "value": "", "channels": ["web"], "sourceKey": "favoriteCategory", "sourceType": ""}]', '{"end": {"date": null, "time": null}, "start": {"date": null, "time": null}, "channels": ["web"], "audiences": null, "traitConditions": [{"key": "favoriteCategory", "type": "string", "value": null, "source": "computed", "operator": "exists"}], "contextConditions": [{"key": "screen.w", "type": "integer", "value": 111, "channel": "store", "operator": "less_than"}]}', '[]', '[]', '[]', '2021-10-27 09:30:25.323187', '2021-10-27 09:30:43.559237');


--
-- Data for Name: audiences; Type: TABLE DATA; Schema: public; Owner: attractify
--

INSERT INTO public.audiences (id, organization_id, name, description, include_anonymous, events, traits, current_set_id, profile_count, created_at, updated_at, refreshed_at) VALUES ('564ee35e-5958-42bc-b822-98781285aaea', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Users With Product In Basket', 'Users that have a product in their basket.', true, '[{"id": "203aa5e3-dd69-4de7-b340-216e61dde639", "count": 1, "operator": "more_or_exactly", "internalId": "8fcd71fd-4d6b-4ebc-b282-b471a4087ef9", "properties": [], "timeWindowOperator": "any_time"}]', '[]', NULL, 0, '2021-10-27 08:45:25.026945', '2021-10-27 08:45:25.026945', '2021-10-27 08:45:25.026945');
INSERT INTO public.audiences (id, organization_id, name, description, include_anonymous, events, traits, current_set_id, profile_count, created_at, updated_at, refreshed_at) VALUES ('f443213f-cc5a-45b8-9662-49e692c745a3', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Checkout Not Completed', 'The user has not completed the checkout process.', true, '[{"id": "e9a352c3-c071-4a4d-b460-98f03b54ec44", "count": 0, "operator": "more_than", "internalId": "a1cf73ab-1251-4ada-bd6b-af042545b8de", "properties": [], "timeWindowOperator": "any_time"}, {"id": "846143c0-3b66-474a-9f3e-6884b4ccefee", "count": 0, "exclude": false, "operator": "more_than", "parentId": "a1cf73ab-1251-4ada-bd6b-af042545b8de", "internalId": "3f1df9b4-85bd-4ae5-8161-e7d6f2524326", "properties": [], "timeWindowOperator": "any_time"}, {"id": "478736a9-90db-43f5-a253-20874572a6c3", "count": 0, "exclude": true, "operator": "more_than", "parentId": "3f1df9b4-85bd-4ae5-8161-e7d6f2524326", "internalId": "e6a71dc5-f814-4dec-b774-459c548e43c1", "properties": [], "timeWindowOperator": "any_time"}]', '[]', NULL, 0, '2021-10-27 08:49:19.946065', '2021-10-27 08:49:19.946065', '2021-10-27 08:49:19.946065');


--
-- Data for Name: profiles; Type: TABLE DATA; Schema: public; Owner: attractify
--



--
-- Data for Name: audience_profiles; Type: TABLE DATA; Schema: public; Owner: attractify
--



--
-- Data for Name: auth_tokens; Type: TABLE DATA; Schema: public; Owner: attractify
--

INSERT INTO public.auth_tokens (id, organization_id, token, channel, created_at) VALUES ('26990bdf-c28f-4f0a-88b6-f7dfd5cfe641', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'demo-pj_r1k3qEmk3Fa7hcUdUO-rAW_6KrGqNJIFrDlDioBScg.qGw3jqw24lEANxpq-G', 'demo', '2021-10-25 06:51:09.593693');


--
-- Data for Name: channels; Type: TABLE DATA; Schema: public; Owner: attractify
--

INSERT INTO public.channels (id, organization_id, name, key, created_at, updated_at) VALUES ('81694542-0e95-4148-89be-4d12bdaab39f', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'demo', 'demo', '2021-10-25 06:50:31.656247', '2021-10-25 06:50:31.656247');
INSERT INTO public.channels (id, organization_id, name, key, created_at, updated_at) VALUES ('26029e2f-be0d-4896-ad19-f88d09776755', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Web', 'web', '2021-10-27 07:41:59.775928', '2021-10-27 07:41:59.775928');
INSERT INTO public.channels (id, organization_id, name, key, created_at, updated_at) VALUES ('7f72e26e-b776-4d5f-bf96-5c0f90d5fe31', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Store', 'store', '2021-10-27 07:42:15.146715', '2021-10-27 07:42:15.146715');
INSERT INTO public.channels (id, organization_id, name, key, created_at, updated_at) VALUES ('dade9bb3-7c43-4452-b38c-4bf5488ddad7', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'App', 'app', '2021-10-27 07:42:24.973623', '2021-10-27 07:42:24.973623');


--
-- Data for Name: computed_traits; Type: TABLE DATA; Schema: public; Owner: attractify
--

INSERT INTO public.computed_traits (id, organization_id, name, key, type, event_id, conditions, properties, created_at, updated_at, refreshed_at) VALUES ('c46cb15d-cc53-4e42-b9c2-33f2ec89ad71', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Last Activity', 'lastActivity', 'last_event', '77a01b0f-1cee-4e62-b728-f3f2d998ff5a', '[]', '{"type": "dateTime", "useTimestamp": true}', '2021-10-27 08:25:44.690557', '2021-10-27 08:25:44.690557', '2021-10-27 08:25:44.690557');
INSERT INTO public.computed_traits (id, organization_id, name, key, type, event_id, conditions, properties, created_at, updated_at, refreshed_at) VALUES ('0829a211-42fb-48bb-9671-10e5fb943342', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Average Product Price', 'averageProductPrice', 'aggregation', '78edfe75-beb1-46f2-a01a-2ca1e4106974', '[]', '{"type": "float", "property": "price", "aggregationType": "avg"}', '2021-10-27 08:27:17.785716', '2021-10-27 08:27:17.785716', '2021-10-27 08:27:17.785716');
INSERT INTO public.computed_traits (id, organization_id, name, key, type, event_id, conditions, properties, created_at, updated_at, refreshed_at) VALUES ('1b5a3cec-e537-4159-b0d3-486456de7017', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Favorite Category', 'favoriteCategory', 'most_frequent', '203aa5e3-dd69-4de7-b340-216e61dde639', '[]', '{"type": "string", "property": "category", "minFrequency": 1}', '2021-10-27 08:28:40.773572', '2021-10-27 08:28:40.773572', '2021-10-27 08:28:40.773572');
INSERT INTO public.computed_traits (id, organization_id, name, key, type, event_id, conditions, properties, created_at, updated_at, refreshed_at) VALUES ('fa60cfd4-c211-4cd6-823e-bd53304770e2', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Total Product Value Viewed', 'totalValue', 'aggregation', '78edfe75-beb1-46f2-a01a-2ca1e4106974', '[{"type": "float", "value": "1", "operator": "greater_than", "property": "price"}]', '{"type": "float", "property": "price", "aggregationType": "sum"}', '2021-10-27 08:30:39.464992', '2021-10-27 08:30:39.464992', '2021-10-27 08:30:39.464992');
INSERT INTO public.computed_traits (id, organization_id, name, key, type, event_id, conditions, properties, created_at, updated_at, refreshed_at) VALUES ('863df0b4-cef2-43d7-b268-f6a492328bb7', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Favorite Size', '	favoriteSize', 'most_frequent', '77fcc4dc-0fab-4511-ac67-b81ee80c7597', '[{"type": "string", "value": "size", "operator": "equals", "property": "attribute"}]', '{"type": "string", "property": "variant", "minFrequency": 1}', '2021-10-27 08:33:00.066331', '2021-10-27 08:33:12.875335', '2021-10-27 08:33:00.066331');
INSERT INTO public.computed_traits (id, organization_id, name, key, type, event_id, conditions, properties, created_at, updated_at, refreshed_at) VALUES ('e003e6b5-ea91-4f72-a3c4-cba8f42f12aa', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Current Cart Value', 'currentCartValue', 'aggregation', 'e9a352c3-c071-4a4d-b460-98f03b54ec44', '[]', '{"type": "float", "property": "total", "aggregationType": "sum"}', '2021-10-27 08:34:13.126717', '2021-10-27 08:34:13.126717', '2021-10-27 08:34:13.126717');
INSERT INTO public.computed_traits (id, organization_id, name, key, type, event_id, conditions, properties, created_at, updated_at, refreshed_at) VALUES ('3e2f6f42-0812-40b7-9734-f80c99c8cacb', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Favorite Shoe Size', 'favoriteShoeSize', 'most_frequent', '77fcc4dc-0fab-4511-ac67-b81ee80c7597', '[{"type": "string", "value": "shoesize", "operator": "equals", "property": "attribute"}]', '{"type": "string", "property": "variant", "minFrequency": 1}', '2021-10-27 08:40:20.834193', '2021-10-27 08:40:20.834193', '2021-10-27 08:40:20.834193');
INSERT INTO public.computed_traits (id, organization_id, name, key, type, event_id, conditions, properties, created_at, updated_at, refreshed_at) VALUES ('cb11b28d-b2a6-4ee7-ac3a-06d1e3f57e14', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'First Page Viewed', 'firstPageVisited', 'first_event', '77a01b0f-1cee-4e62-b728-f3f2d998ff5a', '[]', '{"type": "string", "property": "page"}', '2021-10-27 08:41:41.824015', '2021-10-27 08:41:41.824015', '2021-10-27 08:41:41.824015');
INSERT INTO public.computed_traits (id, organization_id, name, key, type, event_id, conditions, properties, created_at, updated_at, refreshed_at) VALUES ('9fa894ae-b1a8-43dc-879a-0c661702df5b', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Last Page Viewed', 'lastPageViewed', 'last_event', '77a01b0f-1cee-4e62-b728-f3f2d998ff5a', '[]', '{"type": "string", "property": "page"}', '2021-10-27 08:42:52.156455', '2021-10-27 08:42:52.156455', '2021-10-27 08:42:52.156455');


--
-- Data for Name: contexts; Type: TABLE DATA; Schema: public; Owner: attractify
--

INSERT INTO public.contexts (id, organization_id, channel, structure, json_schema, properties, created_at, updated_at) VALUES ('e45ecdb6-330e-4af0-8943-c2635defccff', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'web', '[{"id": "page", "name": "page", "children": [{"id": "path", "name": "path", "properties": {"key": "path", "type": "string", "isRequired": true}}], "properties": {"key": "page", "type": "object", "isRequired": true}}, {"id": "ua", "name": "ua", "properties": {"key": "ua", "type": "string", "isRequired": true}}, {"id": "lang", "name": "lang", "properties": {"key": "lang", "type": "string", "isRequired": true}}, {"id": "conn", "name": "conn", "properties": {"key": "conn", "type": "string", "isRequired": false}}, {"id": "screen", "name": "screen", "children": [{"id": "w", "name": "w", "properties": {"key": "w", "type": "integer", "isRequired": true}}, {"id": "h", "name": "h", "properties": {"key": "h", "type": "integer", "isRequired": true}}], "properties": {"key": "screen", "type": "object", "isRequired": true}}, {"id": "window", "name": "window", "children": [{"id": "w", "name": "w", "properties": {"key": "w", "type": "integer", "isRequired": true}}, {"id": "h", "name": "h", "properties": {"key": "h", "type": "integer", "isRequired": true}}], "properties": {"key": "window", "type": "object", "isRequired": true}}]', '{"type": "object", "required": ["page", "ua", "lang", "screen", "window"], "properties": {"ua": {"type": "string", "additionalItems": false, "additionalProperties": false}, "conn": {"type": "string", "additionalItems": false, "additionalProperties": false}, "lang": {"type": "string", "additionalItems": false, "additionalProperties": false}, "page": {"type": "object", "required": ["path"], "properties": {"path": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}, "screen": {"type": "object", "required": ["w", "h"], "properties": {"h": {"type": "integer", "additionalItems": false, "additionalProperties": false}, "w": {"type": "integer", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}, "window": {"type": "object", "required": ["w", "h"], "properties": {"h": {"type": "integer", "additionalItems": false, "additionalProperties": false}, "w": {"type": "integer", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "page.path", "type": "string"}, {"key": "ua", "type": "string"}, {"key": "lang", "type": "string"}, {"key": "conn", "type": "string"}, {"key": "screen.w", "type": "integer"}, {"key": "screen.h", "type": "integer"}, {"key": "window.w", "type": "integer"}, {"key": "window.h", "type": "integer"}]', '2021-10-27 08:14:36.227648', '2021-10-27 08:14:36.227648');


--
-- Data for Name: custom_traits; Type: TABLE DATA; Schema: public; Owner: attractify
--

INSERT INTO public.custom_traits (organization_id, structure, json_schema, properties, created_at, updated_at) VALUES ('5647e66f-ef2c-4f29-9584-6f788f573e94', '[{"id": "age", "name": "age", "properties": {"key": "age", "type": "integer", "isRequired": true}}, {"id": "gender", "name": "gender", "properties": {"key": "gender", "type": "string", "isRequired": true}}, {"id": "customerSince", "name": "customerSince", "properties": {"key": "customerSince", "type": "dateTime", "isRequired": true}}]', '{"type": "object", "required": ["age", "gender", "customerSince"], "properties": {"age": {"type": "integer", "additionalItems": false, "additionalProperties": false}, "gender": {"type": "string", "additionalItems": false, "additionalProperties": false}, "customerSince": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "age", "type": "integer"}, {"key": "gender", "type": "string"}, {"key": "customerSince", "type": "dateTime"}]', '2021-10-27 08:16:39.924184', '2021-10-27 08:16:39.924184');


--
-- Data for Name: events; Type: TABLE DATA; Schema: public; Owner: attractify
--

INSERT INTO public.events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('203aa5e3-dd69-4de7-b340-216e61dde639', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Category Viewed', 'The user has viewed a category page.', 0, '[{"id": "category", "name": "category", "properties": {"key": "category", "type": "string", "isRequired": true}}]', '{"type": "object", "required": ["category"], "properties": {"category": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "category", "type": "string"}]', '2021-10-27 07:43:44.777799', '2021-10-27 07:43:44.777799');
INSERT INTO public.events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('478736a9-90db-43f5-a253-20874572a6c3', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Checkout Completed', 'The user has completed the checkout.', 0, '[{"id": "total", "name": "total", "properties": {"key": "total", "type": "float", "isRequired": true}}, {"id": "items", "name": "items", "properties": {"key": "items", "type": "integer", "isRequired": true}}]', '{"type": "object", "required": ["total", "items"], "properties": {"items": {"type": "integer", "additionalItems": false, "additionalProperties": false}, "total": {"type": "number", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "total", "type": "float"}, {"key": "items", "type": "integer"}]', '2021-10-27 07:45:24.205949', '2021-10-27 07:45:24.205949');
INSERT INTO public.events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('846143c0-3b66-474a-9f3e-6884b4ccefee', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Checkout Started', 'The user has started the checkout process.', 0, '[{"id": "total", "name": "total", "properties": {"key": "total", "type": "float", "isRequired": true}}, {"id": "items", "name": "items", "properties": {"key": "items", "type": "integer", "isRequired": true}}]', '{"type": "object", "required": ["total", "items"], "properties": {"items": {"type": "integer", "additionalItems": false, "additionalProperties": false}, "total": {"type": "number", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "total", "type": "float"}, {"key": "items", "type": "integer"}]', '2021-10-27 07:46:49.493265', '2021-10-27 07:46:49.493265');
INSERT INTO public.events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('7eba1e2d-48a4-49ba-8588-ecaaecac0b7e', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Field Focussed', 'The user has focussed a field in the checkout form.', 0, '[{"id": "field", "name": "field", "properties": {"key": "field", "type": "string", "isRequired": true}}, {"id": "empty", "name": "empty", "properties": {"key": "empty", "type": "boolean", "isRequired": true}}]', '{"type": "object", "required": ["field", "empty"], "properties": {"empty": {"type": "boolean", "additionalItems": false, "additionalProperties": false}, "field": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "field", "type": "string"}, {"key": "empty", "type": "boolean"}]', '2021-10-27 07:55:12.367147', '2021-10-27 07:55:12.367147');
INSERT INTO public.events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('77a01b0f-1cee-4e62-b728-f3f2d998ff5a', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Page Viewed', 'The user has views a page.', 0, '[{"id": "page", "name": "page", "properties": {"key": "page", "type": "string", "isRequired": true}}]', '{"type": "object", "required": ["page"], "properties": {"page": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "page", "type": "string"}]', '2021-10-27 07:56:23.526705', '2021-10-27 07:56:23.526705');
INSERT INTO public.events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('0abea2ca-9846-4412-9117-adb706aacc56', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Product Displayed', 'A new product has been displayed.', 0, '[{"id": "product_id", "name": "product_id", "properties": {"key": "product_id", "type": "string", "isRequired": true}}]', '{"type": "object", "required": ["product_id"], "properties": {"product_id": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "product_id", "type": "string"}]', '2021-10-27 07:57:41.88693', '2021-10-27 07:57:41.88693');
INSERT INTO public.events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('78edfe75-beb1-46f2-a01a-2ca1e4106974', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Product Viewed', 'The user has viewed a product.', 0, '[{"id": "product", "name": "product", "properties": {"key": "product", "type": "string", "isRequired": true}}, {"id": "price", "name": "price", "properties": {"key": "price", "type": "float", "isRequired": true}}]', '{"type": "object", "required": ["product", "price"], "properties": {"price": {"type": "number", "additionalItems": false, "additionalProperties": false}, "product": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "product", "type": "string"}, {"key": "price", "type": "float"}]', '2021-10-27 07:59:16.011025', '2021-10-27 07:59:16.011025');
INSERT INTO public.events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('44ecac83-dea6-4536-ac2e-1810dac4a164', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Searched', 'The user has search for a product.', 0, '[{"id": "query", "name": "query", "properties": {"key": "query", "type": "string", "isRequired": true}}]', '{"type": "object", "required": ["query"], "properties": {"query": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "query", "type": "string"}]', '2021-10-27 08:00:34.822401', '2021-10-27 08:00:34.822401');
INSERT INTO public.events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('e9a352c3-c071-4a4d-b460-98f03b54ec44', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Updated Cart', 'A product has been added to or removed from the shopping cart.', 0, '[{"id": "product", "name": "product", "properties": {"key": "product", "type": "string", "isRequired": true}}, {"id": "price", "name": "price", "properties": {"key": "price", "type": "float", "isRequired": true}}, {"id": "quantity", "name": "quantity", "properties": {"key": "quantity", "type": "integer", "isRequired": true}}, {"id": "variant", "name": "variant", "children": [{"id": "shoesize", "name": "shoesize", "properties": {"key": "shoesize", "type": "string"}}, {"id": "color", "name": "color", "properties": {"key": "color", "type": "string"}}, {"id": "size", "name": "size", "properties": {"key": "size", "type": "string"}}], "properties": {"key": "variant", "type": "object"}}, {"id": "total", "name": "total", "properties": {"key": "total", "type": "float", "isRequired": true}}, {"id": "operation", "name": "operation", "properties": {"key": "operation", "type": "string", "isRequired": true}}]', '{"type": "object", "required": ["product", "price", "quantity", "total", "operation"], "properties": {"price": {"type": "number", "additionalItems": false, "additionalProperties": false}, "total": {"type": "number", "additionalItems": false, "additionalProperties": false}, "product": {"type": "string", "additionalItems": false, "additionalProperties": false}, "variant": {"type": "object", "properties": {"size": {"type": "string", "additionalItems": false, "additionalProperties": false}, "color": {"type": "string", "additionalItems": false, "additionalProperties": false}, "shoesize": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}, "quantity": {"type": "integer", "additionalItems": false, "additionalProperties": false}, "operation": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "product", "type": "string"}, {"key": "price", "type": "float"}, {"key": "quantity", "type": "integer"}, {"key": "variant.shoesize", "type": "string"}, {"key": "variant.color", "type": "string"}, {"key": "variant.size", "type": "string"}, {"key": "total", "type": "float"}, {"key": "operation", "type": "string"}]', '2021-10-27 08:07:13.849731', '2021-10-27 08:07:13.849731');
INSERT INTO public.events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('77fcc4dc-0fab-4511-ac67-b81ee80c7597', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'Variant Selected', 'The user has selected a product variant.', 0, '[{"id": "product", "name": "product", "properties": {"key": "product", "type": "string", "isRequired": true}}, {"id": "variant", "name": "variant", "properties": {"key": "variant", "type": "string", "isRequired": true}}, {"id": "attribute", "name": "attribute", "properties": {"key": "attribute", "type": "string", "isRequired": true}}]', '{"type": "object", "required": ["product", "variant", "attribute"], "properties": {"product": {"type": "string", "additionalItems": false, "additionalProperties": false}, "variant": {"type": "string", "additionalItems": false, "additionalProperties": false}, "attribute": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "product", "type": "string"}, {"key": "variant", "type": "string"}, {"key": "attribute", "type": "string"}]', '2021-10-27 08:09:01.750382', '2021-10-27 08:09:01.750382');


--
-- Data for Name: locked_profile_identities; Type: TABLE DATA; Schema: public; Owner: attractify
--



--
-- Data for Name: profile_identities; Type: TABLE DATA; Schema: public; Owner: attractify
--



--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: attractify
--

INSERT INTO public.users (id, organization_id, email, password, salt, name, role, logged_out_at, created_at, updated_at) VALUES ('7b63afee-ac87-4573-8c38-b90d659c7471', '5647e66f-ef2c-4f29-9584-6f788f573e94', 'demo@example.com', '\x96d00a74060e546acd4701fbab2da37d175442a0fcd90df06a50bfa2c3f121ec1dc68461719c2875968312c0b85ad9bc7e2a44e683b33cabde53c2cb5963bb6c', '\xc14ac4309af48f30a7cdb7fa569f2bcc', 'demo', 'admin', NULL, '2021-10-25 06:50:06.307233', '2021-10-25 06:50:06.307233');


--
-- PostgreSQL database dump complete
--

