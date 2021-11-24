--
-- PostgreSQL database dump
--

-- Dumped from database version 14.1 (Debian 14.1-1.pgdg110+1)
-- Dumped by pg_dump version 14.1 (Debian 14.1-1.pgdg110+1)

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

INSERT INTO public.organizations (id, name, key, timezone, created_at, updated_at) VALUES ('bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Sportify', '\x522b7e5f4047552528743979792f3b445d4b5b77755e767e227e25394e3e5173', 'Europe/Berlin', '2020-06-26 11:41:24.202', '2021-02-05 15:04:26.23');


--
-- Data for Name: actions; Type: TABLE DATA; Schema: public; Owner: attractify
--

INSERT INTO public.actions (id, organization_id, name, type, tags, state, properties, targeting, capping, hooks, test_users, created_at, updated_at) VALUES ('178af873-67e9-43ec-9e53-b1e831475b31', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Not Available', 'not_available', '["improvement"]', 'active', '[{"name": "title", "type": "text", "value": "Benachrichtige mich bei Verfügbarkeit", "channels": ["web", "store", "app"], "sourceKey": "", "sourceType": ""}, {"name": "body", "type": "text", "value": "Leider ist der ausgewählte Artikel nicht mehr verfügbar, aber wir bekommen diesen sicher wieder rein. Gerne informieren wir Dich dann per E-Mail.", "channels": ["web", "store", "app"], "sourceKey": "", "sourceType": ""}, {"name": "cta", "type": "text", "value": "Benachrichtigung erhalten", "channels": ["web", "store", "app"], "sourceKey": "", "sourceType": ""}]', '{"end": {"date": null, "time": null}, "start": {"date": null, "time": null}, "channels": ["web", "store", "app"], "audiences": null, "traitConditions": [], "contextConditions": []}', '[{"count": 2, "event": "hidden", "group": "user", "within": 0, "channels": ["web", "store", "app"]}, {"count": 1, "event": "accepted", "group": "user", "within": 0, "channels": ["store", "web", "app"]}]', '[]', '[]', '2020-07-13 07:37:51.1', '2020-07-13 13:23:00.734');
INSERT INTO public.actions (id, organization_id, name, type, tags, state, properties, targeting, capping, hooks, test_users, created_at, updated_at) VALUES ('1b05b65c-cc97-43b5-9a2b-875e7f7d6da4', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Incentivize Abandoned Cart Users', 'coupon_overlay', '["upselling"]', 'active', '[{"name": "title", "type": "text", "value": "Schließe jetzt Deinen Kauf ab!!!!", "channels": ["web"], "sourceKey": "", "sourceType": ""}, {"name": "body", "type": "text", "value": "Wir dachten uns, dass Dir vielleicht dieser 10% Gutschein hilft, Deinen Kauf abzuschließen?", "channels": ["web"], "sourceKey": "", "sourceType": ""}, {"name": "cta", "type": "text", "value": "Gutschein jetzt einlösen", "channels": ["web"], "sourceKey": "", "sourceType": ""}, {"name": "couponCode", "type": "text", "value": "B4ZXQZVX", "channels": ["web"], "sourceKey": "", "sourceType": ""}]', '{"end": {"date": null, "time": null}, "start": {"date": null, "time": null}, "channels": ["web"], "audiences": ["80a45fce-2136-4954-a403-6894957a5ad4"], "traitConditions": null, "contextConditions": []}', '[{"count": 1, "event": "accepted", "group": "user", "within": 0, "channels": ["web"]}, {"count": 2, "event": "hidden", "group": "user", "within": 0, "channels": ["web"]}]', '[]', '[]', '2020-07-02 08:42:05.258', '2021-04-08 08:17:53.349');
INSERT INTO public.actions (id, organization_id, name, type, tags, state, properties, targeting, capping, hooks, test_users, created_at, updated_at) VALUES ('528cab3c-416c-48d1-9b59-614aec2df897', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Preselect Size', 'size_preset', '["improvement"]', 'active', '[{"name": "size", "type": "computed_trait", "value": "", "channels": ["web"], "sourceKey": "favoriteSize", "sourceType": ""}]', '{"end": {"date": null, "time": null}, "start": {"date": null, "time": null}, "channels": ["web"], "audiences": null, "traitConditions": [{"key": "favoriteSize", "type": "string", "value": null, "source": "computed", "operator": "exists"}], "contextConditions": [{"key": "page.path", "type": "string", "value": "/?product=", "operator": "starts_with"}]}', '[]', '[]', '[]', '2020-07-07 07:10:58.013', '2020-07-07 08:11:09.677');
INSERT INTO public.actions (id, organization_id, name, type, tags, state, properties, targeting, capping, hooks, test_users, created_at, updated_at) VALUES ('6cfc5ada-568f-4dc8-8d2e-ea06604c3ce8', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Preselect category', 'favorite_category', '["improvement"]', 'active', '[{"name": "category", "type": "computed_trait", "value": "", "channels": ["web"], "sourceKey": "favoriteCategory", "sourceType": ""}]', '{"end": {"date": null, "time": null}, "start": {"date": null, "time": null}, "channels": ["web"], "audiences": null, "traitConditions": [{"key": "favoriteCategory", "type": "string", "value": null, "source": "computed", "operator": "exists"}], "contextConditions": [{"key": "page.path", "type": "string", "value": "/?post_type=product", "operator": "starts_with"}]}', '[]', '[]', '[]', '2020-07-07 07:20:37.433', '2020-07-07 08:11:24.153');


--
-- Data for Name: audiences; Type: TABLE DATA; Schema: public; Owner: attractify
--

INSERT INTO public.audiences (id, organization_id, name, description, include_anonymous, events, traits, current_set_id, profile_count, created_at, updated_at, refreshed_at) VALUES ('7049b736-d2f2-4343-a1e4-caf6c6a325c8', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Users With Product In Basket', 'Users that have a product in their basket.', true, '[{"id": "872fafc4-29e1-4b05-8f59-72357843362b", "count": 1, "operator": "more_or_exactly", "internalId": "2e8e1663-8064-4847-ad77-53974d021f14", "properties": [], "timeWindowOperator": "any_time"}]', '[]', '045a20a2-4593-4536-810c-936429d293ac', 0, '2020-06-30 12:25:19.389', '2021-02-11 11:02:28.029', '2021-02-11 11:02:28.029');
INSERT INTO public.audiences (id, organization_id, name, description, include_anonymous, events, traits, current_set_id, profile_count, created_at, updated_at, refreshed_at) VALUES ('80a45fce-2136-4954-a403-6894957a5ad4', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Checkout Not Completed', 'The user has not completed the checkout process.', true, '[{"id": "5b762e7f-600f-4f6a-abab-fdc22f80f9ac", "count": 0, "operator": "more_than", "internalId": "fe8354e0-ef2d-4875-9fe0-31ab26b1dd32", "properties": [], "timeWindowOperator": "any_time"}, {"id": "627f9036-dabd-4148-835e-fe0e284a7b23", "count": 0, "exclude": false, "operator": "more_than", "parentId": "fe8354e0-ef2d-4875-9fe0-31ab26b1dd32", "internalId": "c0391e1e-2bf7-427b-b77e-d11002324071", "properties": [], "timeWindowOperator": "any_time"}, {"id": "7eeec417-7f1d-47b5-8548-af25ba50143d", "count": 0, "exclude": true, "operator": "more_than", "parentId": "c0391e1e-2bf7-427b-b77e-d11002324071", "internalId": "d1d5c7b1-4d3d-4027-9df8-9c5b7b04cba4", "properties": [], "timeWindowOperator": "any_time"}]', '[]', 'c6b8cdf9-cacd-4005-8fe8-12523ad54979', 1, '2020-07-01 12:18:47.049', '2021-11-05 09:20:52.351', '2021-11-05 09:20:52.351');


--
-- Data for Name: profiles; Type: TABLE DATA; Schema: public; Owner: attractify
--



--
-- Data for Name: audience_profiles; Type: TABLE DATA; Schema: public; Owner: attractify
--



--
-- Data for Name: auth_tokens; Type: TABLE DATA; Schema: public; Owner: attractify
--

INSERT INTO public.auth_tokens (id, organization_id, token, channel, created_at) VALUES ('385686a1-5b19-476e-9e08-4b28ea5c7cb0', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'web-0IzIpVVc8oIU17MObozfrx.UWy02c7SPhFbzFU7mRZ.ZwDDRIJx.mu1bwOn3s-oa', 'web', '2021-11-23 14:59:29.096593');


--
-- Data for Name: channels; Type: TABLE DATA; Schema: public; Owner: attractify
--

INSERT INTO public.channels (id, organization_id, name, key, created_at, updated_at) VALUES ('4e3f9f5b-f63f-415c-853a-cdd633309455', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Web', 'web', '2020-06-26 11:46:49.417', '2020-06-26 11:46:49.417');
INSERT INTO public.channels (id, organization_id, name, key, created_at, updated_at) VALUES ('9044e340-4cb5-45d7-80f2-782083db9888', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Store', 'store', '2020-07-03 06:46:44.995', '2020-07-03 06:46:44.995');
INSERT INTO public.channels (id, organization_id, name, key, created_at, updated_at) VALUES ('9d430b8a-f7d6-4918-b560-203b15ff1945', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'App', 'app', '2020-07-03 06:46:23.368', '2020-07-03 06:46:23.368');


--
-- Data for Name: computed_traits; Type: TABLE DATA; Schema: public; Owner: attractify
--

INSERT INTO public.computed_traits (id, organization_id, name, key, type, event_id, conditions, properties, created_at, updated_at, refreshed_at) VALUES ('0e9b18b3-b7ed-459c-99cd-f6aab5f4d6ee', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Last Activity', 'lastActivity', 'last_event', '5d9e3c6f-3fe5-4148-a8f1-77db058eb7be', '[]', '{"type": "dateTime", "useTimestamp": true}', '2020-07-01 08:41:39.438', '2020-07-06 20:43:18.272', '2020-07-01 08:41:39.438');
INSERT INTO public.computed_traits (id, organization_id, name, key, type, event_id, conditions, properties, created_at, updated_at, refreshed_at) VALUES ('6ff4980c-e973-4521-bd7f-e147358ef6e8', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Average Product Price', 'averageProductPrice', 'aggregation', '77c2382f-e9d8-4a2a-b0d0-d1dc3a45bd0a', '[]', '{"type": "float", "property": "price", "aggregationType": "avg"}', '2020-06-30 21:04:47.685', '2020-07-06 20:43:28.062', '2020-06-30 21:04:47.685');
INSERT INTO public.computed_traits (id, organization_id, name, key, type, event_id, conditions, properties, created_at, updated_at, refreshed_at) VALUES ('87fecab9-b050-4013-8821-d9b26bd62fae', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Favorite Category', 'favoriteCategory', 'most_frequent', '872fafc4-29e1-4b05-8f59-72357843362b', '[]', '{"type": "string", "property": "category", "minFrequency": 1}', '2020-06-30 09:00:27.24', '2021-02-05 13:57:20.98', '2020-06-30 09:00:27.24');
INSERT INTO public.computed_traits (id, organization_id, name, key, type, event_id, conditions, properties, created_at, updated_at, refreshed_at) VALUES ('912bb4ea-c551-4787-9876-71258456ad53', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Total Product Value Viewed', 'totalValue', 'aggregation', '77c2382f-e9d8-4a2a-b0d0-d1dc3a45bd0a', '[{"type": "float", "value": "1", "operator": "greater_than", "property": "price"}]', '{"type": "float", "property": "price", "aggregationType": "sum"}', '2021-01-25 13:31:31.232', '2021-01-25 13:31:31.232', '2021-01-25 13:31:31.232');
INSERT INTO public.computed_traits (id, organization_id, name, key, type, event_id, conditions, properties, created_at, updated_at, refreshed_at) VALUES ('9ad88d9f-c198-496a-a6ed-3a0f26e5aa4f', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Favorite Size', 'favoriteSize', 'most_frequent', '3c11343c-b51a-4aef-9ac9-d8f2972fcb0c', '[{"type": "string", "value": "groesse", "operator": "equals", "property": "attribute"}]', '{"type": "string", "property": "variant", "minFrequency": 1}', '2020-07-02 08:07:04.1', '2021-02-05 14:32:05.233', '2020-07-02 08:07:04.1');
INSERT INTO public.computed_traits (id, organization_id, name, key, type, event_id, conditions, properties, created_at, updated_at, refreshed_at) VALUES ('9c97cd3f-0d96-4612-bc18-d97e87c2abe6', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Current Cart Value', 'currentCartValue', 'aggregation', '5b762e7f-600f-4f6a-abab-fdc22f80f9ac', '[]', '{"type": "float", "property": "total", "aggregationType": "sum"}', '2020-07-01 09:01:25.464', '2020-07-06 20:43:54.406', '2020-07-01 09:01:25.464');
INSERT INTO public.computed_traits (id, organization_id, name, key, type, event_id, conditions, properties, created_at, updated_at, refreshed_at) VALUES ('b9d7d130-d817-40a5-9cc4-b265d5c6c021', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Favorite Shoe Size', 'favoriteShoeSize', 'most_frequent', '3c11343c-b51a-4aef-9ac9-d8f2972fcb0c', '[{"type": "string", "value": "schuhgroesse", "operator": "equals", "property": "attribute"}]', '{"type": "string", "property": "variant", "minFrequency": 1}', '2020-07-02 08:24:26.812', '2021-02-05 13:56:32.277', '2020-07-02 08:24:26.812');
INSERT INTO public.computed_traits (id, organization_id, name, key, type, event_id, conditions, properties, created_at, updated_at, refreshed_at) VALUES ('c496f648-724b-455e-957f-8c6c77a7994c', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'First Page Viewed', 'firstPageVisited', 'first_event', '5d9e3c6f-3fe5-4148-a8f1-77db058eb7be', '[]', '{"type": "string", "property": "page", "useTimestamp": false}', '2020-07-01 07:38:57.974', '2020-07-06 20:44:14.404', '2020-07-01 07:38:57.974');
INSERT INTO public.computed_traits (id, organization_id, name, key, type, event_id, conditions, properties, created_at, updated_at, refreshed_at) VALUES ('fb82775c-2f1a-4a19-8815-b81ec6de0f07', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Last Page Viewed', 'lastPageViewed', 'last_event', '5d9e3c6f-3fe5-4148-a8f1-77db058eb7be', '[]', '{"type": "string", "property": "page"}', '2020-07-01 07:39:28.358', '2020-07-06 20:44:26.447', '2020-07-01 07:39:28.358');


--
-- Data for Name: contexts; Type: TABLE DATA; Schema: public; Owner: attractify
--

INSERT INTO public.contexts (id, organization_id, channel, structure, json_schema, properties, created_at, updated_at) VALUES ('963eb52a-0ecc-4b09-8d34-9ec8b9c4553e', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'web', '[{"id": "page", "name": "page", "children": [{"id": "path", "name": "path", "properties": {"key": "path", "type": "string", "isRequired": true}}], "properties": {"key": "page", "type": "object", "isRequired": true}}, {"id": "ua", "name": "ua", "properties": {"key": "ua", "type": "string", "isRequired": true}}, {"id": "lang", "name": "lang", "properties": {"key": "lang", "type": "string", "isRequired": true}}, {"id": "conn", "name": "conn", "properties": {"key": "conn", "type": "string"}}, {"id": "screen", "name": "screen", "children": [{"id": "w", "name": "w", "properties": {"key": "w", "type": "integer", "isRequired": true}}, {"id": "h", "name": "h", "properties": {"key": "h", "type": "integer", "isRequired": true}}], "properties": {"key": "screen", "type": "object", "isRequired": true}}, {"id": "window", "name": "window", "children": [{"id": "w", "name": "w", "properties": {"key": "w", "type": "integer", "isRequired": true}}, {"id": "h", "name": "h", "properties": {"key": "h", "type": "integer", "isRequired": true}}], "properties": {"key": "window", "type": "object", "isRequired": true}}]', '{"type": "object", "required": ["page", "ua", "lang", "screen", "window"], "properties": {"ua": {"type": "string", "additionalItems": false, "additionalProperties": false}, "conn": {"type": "string", "additionalItems": false, "additionalProperties": false}, "lang": {"type": "string", "additionalItems": false, "additionalProperties": false}, "page": {"type": "object", "required": ["path"], "properties": {"path": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}, "screen": {"type": "object", "required": ["w", "h"], "properties": {"h": {"type": "integer", "additionalItems": false, "additionalProperties": false}, "w": {"type": "integer", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}, "window": {"type": "object", "required": ["w", "h"], "properties": {"h": {"type": "integer", "additionalItems": false, "additionalProperties": false}, "w": {"type": "integer", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "page.path", "type": "string"}, {"key": "ua", "type": "string"}, {"key": "lang", "type": "string"}, {"key": "conn", "type": "string"}, {"key": "screen.w", "type": "integer"}, {"key": "screen.h", "type": "integer"}, {"key": "window.w", "type": "integer"}, {"key": "window.h", "type": "integer"}]', '2020-06-30 08:51:08.706', '2020-06-30 08:51:08.706');


--
-- Data for Name: custom_traits; Type: TABLE DATA; Schema: public; Owner: attractify
--

INSERT INTO public.custom_traits (organization_id, structure, json_schema, properties, created_at, updated_at) VALUES ('bc70b33d-c77f-4fe3-813d-a2605c0915cb', '[{"id": "age", "name": "age", "properties": {"key": "age", "type": "integer", "isRequired": true}}, {"id": "gender", "name": "gender", "properties": {"key": "gender", "type": "string", "isRequired": true}}, {"id": "customerSince", "name": "customerSince", "properties": {"key": "customerSince", "type": "dateTime", "isRequired": true}}]', '{"type": "object", "required": ["age", "gender", "customerSince"], "properties": {"age": {"type": "integer", "additionalItems": false, "additionalProperties": false}, "gender": {"type": "string", "additionalItems": false, "additionalProperties": false}, "customerSince": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "age", "type": "integer"}, {"key": "gender", "type": "string"}, {"key": "customerSince", "type": "dateTime"}]', '2020-07-05 19:24:53.302', '2020-07-15 07:43:56.114');


--
-- Data for Name: events; Type: TABLE DATA; Schema: public; Owner: attractify
--

INSERT INTO public.events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('3c11343c-b51a-4aef-9ac9-d8f2972fcb0c', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Variant Selected', 'The user has selected a product variant.', 0, '[{"id": "product", "name": "product", "properties": {"key": "product", "type": "string", "isRequired": true}}, {"id": "variant", "name": "variant", "properties": {"key": "variant", "type": "string", "isRequired": true}}, {"id": "attribute", "name": "attribute", "properties": {"key": "attribute", "type": "string", "isRequired": true}}]', '{"type": "object", "required": ["product", "variant", "attribute"], "properties": {"product": {"type": "string", "additionalItems": false, "additionalProperties": false}, "variant": {"type": "string", "additionalItems": false, "additionalProperties": false}, "attribute": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "product", "type": "string"}, {"key": "variant", "type": "string"}, {"key": "attribute", "type": "string"}]', '2020-07-02 08:04:18.014', '2021-02-05 14:31:46.916');
INSERT INTO public.events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('5b762e7f-600f-4f6a-abab-fdc22f80f9ac', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Updated Cart', 'A product has been added to or removed from the shopping cart.', 0, '[{"id": "product", "name": "product", "properties": {"key": "product", "type": "string", "isRequired": true}}, {"id": "price", "name": "price", "properties": {"key": "price", "type": "float", "isRequired": true}}, {"id": "quantity", "name": "quantity", "properties": {"key": "quantity", "type": "integer", "isRequired": true}}, {"id": "variant", "name": "variant", "children": [{"id": "schuhgroesse", "name": "schuhgroesse", "properties": {"key": "schuhgroesse", "type": "string", "isRequired": false}}, {"id": "farbe", "name": "farbe", "properties": {"key": "farbe", "type": "string"}}, {"id": "groesse", "name": "groesse", "properties": {"key": "groesse", "type": "string"}}], "properties": {"key": "variant", "type": "object"}}, {"id": "total", "name": "total", "properties": {"key": "total", "type": "float", "isRequired": true}}, {"id": "operation", "name": "operation", "properties": {"key": "operation", "type": "string", "isRequired": true}}]', '{"type": "object", "required": ["product", "price", "quantity", "total", "operation"], "properties": {"price": {"type": "number", "additionalItems": false, "additionalProperties": false}, "total": {"type": "number", "additionalItems": false, "additionalProperties": false}, "product": {"type": "string", "additionalItems": false, "additionalProperties": false}, "variant": {"type": "object", "properties": {"farbe": {"type": "string", "additionalItems": false, "additionalProperties": false}, "groesse": {"type": "string", "additionalItems": false, "additionalProperties": false}, "schuhgroesse": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}, "quantity": {"type": "integer", "additionalItems": false, "additionalProperties": false}, "operation": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "product", "type": "string"}, {"key": "price", "type": "float"}, {"key": "quantity", "type": "integer"}, {"key": "variant.schuhgroesse", "type": "string"}, {"key": "variant.farbe", "type": "string"}, {"key": "variant.groesse", "type": "string"}, {"key": "total", "type": "float"}, {"key": "operation", "type": "string"}]', '2020-06-30 08:15:47.753', '2020-07-01 10:39:34.436');
INSERT INTO public.events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('5d9e3c6f-3fe5-4148-a8f1-77db058eb7be', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Page Viewed', 'The user has views a page.', 0, '[{"id": "page", "name": "page", "properties": {"key": "page", "type": "string", "isRequired": true}}]', '{"type": "object", "required": ["page"], "properties": {"page": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "page", "type": "string"}]', '2020-07-01 07:38:22.204', '2020-07-01 07:38:22.204');
INSERT INTO public.events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('627f9036-dabd-4148-835e-fe0e284a7b23', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Checkout Started', 'The user has started the checkout process.', 0, '[{"id": "total", "name": "total", "properties": {"key": "total", "type": "float", "isRequired": true}}, {"id": "items", "name": "items", "properties": {"key": "items", "type": "integer", "isRequired": true}}]', '{"type": "object", "required": ["total", "items"], "properties": {"items": {"type": "integer", "additionalItems": false, "additionalProperties": false}, "total": {"type": "number", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "total", "type": "float"}, {"key": "items", "type": "integer"}]', '2020-07-01 12:09:56.082', '2020-07-01 12:09:56.082');
INSERT INTO public.events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('77c2382f-e9d8-4a2a-b0d0-d1dc3a45bd0a', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Product Viewed', 'The user has viewed a product.', 0, '[{"id": "product", "name": "product", "children": [], "properties": {"key": "product", "type": "string", "isRequired": true}}, {"id": "price", "name": "price", "properties": {"key": "price", "type": "float", "isRequired": true}}]', '{"type": "object", "required": ["product", "price"], "properties": {"price": {"type": "number", "additionalItems": false, "additionalProperties": false}, "product": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "product", "type": "string"}, {"key": "price", "type": "float"}]', '2020-06-30 08:54:43.643', '2020-06-30 08:59:15.896');
INSERT INTO public.events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('7adb3c3a-a1eb-4123-9d15-a6383cf6e8a0', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Searched', 'The user has search for a product.', 0, '[{"id": "query", "name": "query", "properties": {"key": "query", "type": "string", "isRequired": true}}]', '{"type": "object", "required": ["query"], "properties": {"query": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "query", "type": "string"}]', '2020-07-01 12:10:57.651', '2020-07-01 12:10:57.651');
INSERT INTO public.events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('7eeec417-7f1d-47b5-8548-af25ba50143d', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Checkout Completed', 'The user has completed the checkout.', 0, '[{"id": "total", "name": "total", "properties": {"key": "total", "type": "float", "isRequired": true}}, {"id": "items", "name": "items", "properties": {"key": "items", "type": "integer", "isRequired": true}}]', '{"type": "object", "required": ["total", "items"], "properties": {"items": {"type": "integer", "additionalItems": false, "additionalProperties": false}, "total": {"type": "number", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "total", "type": "float"}, {"key": "items", "type": "integer"}]', '2020-07-01 12:05:38.673', '2020-07-01 12:05:38.673');
INSERT INTO public.events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('872fafc4-29e1-4b05-8f59-72357843362b', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Category Viewed', 'The user has viewed a category page.', 0, '[{"id": "category", "name": "category", "properties": {"key": "category", "type": "string", "isRequired": true}}]', '{"type": "object", "required": ["category"], "properties": {"category": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "category", "type": "string"}]', '2020-06-30 08:54:14.623', '2020-06-30 08:54:14.623');
INSERT INTO public.events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('df6e1ca4-4c37-4ac0-8091-1cd010201233', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Field Focussed', 'The user has focussed a field in the checkout form.', 0, '[{"id": "field", "name": "field", "properties": {"key": "field", "type": "string", "isRequired": true}}, {"id": "empty", "name": "empty", "properties": {"key": "empty", "type": "boolean", "isRequired": true}}]', '{"type": "object", "required": ["field", "empty"], "properties": {"empty": {"type": "boolean", "additionalItems": false, "additionalProperties": false}, "field": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "field", "type": "string"}, {"key": "empty", "type": "boolean"}]', '2020-07-01 12:12:36.371', '2020-07-01 12:13:44.82');
INSERT INTO public.events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('fa29c00b-7c5a-4b1e-94ab-2e57d56d9ff2', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Product Displayed', 'A new product has been displayed', 0, '[{"id": "product_id", "name": "product_id", "properties": {"key": "product_id", "type": "string", "isRequired": true}}]', '{"type": "object", "required": ["product_id"], "properties": {"product_id": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "product_id", "type": "string"}]', '2021-01-14 08:47:30.453', '2021-01-14 08:47:30.453');


--
-- Data for Name: locked_profile_identities; Type: TABLE DATA; Schema: public; Owner: attractify
--



--
-- Data for Name: profile_identities; Type: TABLE DATA; Schema: public; Owner: attractify
--



--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: attractify
--

INSERT INTO public.users (id, organization_id, email, password, salt, name, role, logged_out_at, created_at, updated_at) VALUES ('dde601b1-3fab-466b-846d-5a0d0a298e91', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'demo@example.com', '\xd1c6b4d329d65630356b233e1cce57ff3513b4b3288cfa6cb46512e69dada89eecee2d84d982bc5cf59e945d827b7f51d613cc690205c3414135bc09789ff40b', '\x7e16f8b545a3e06082e511c245870794', 'Demo', 'admin', NULL, '2021-11-23 14:35:19.359234', '2021-11-23 14:35:19.359234');


--
-- PostgreSQL database dump complete
--

