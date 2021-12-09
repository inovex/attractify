INSERT INTO organizations(id, "name", "key", timezone, created_at, updated_at) VALUES('bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, 'Sportify', decode('522B7E5F4047552528743979792F3B445D4B5B77755E767E227E25394E3E5173','hex'), 'Europe/Berlin', '2020-06-26 11:41:24.202', '2021-02-05 15:04:26.230');

INSERT INTO users (id, organization_id, email, password, salt, name, role, logged_out_at, created_at, updated_at) VALUES ('dde601b1-3fab-466b-846d-5a0d0a298e91', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, 'demo@example.com', '\xd1c6b4d329d65630356b233e1cce57ff3513b4b3288cfa6cb46512e69dada89eecee2d84d982bc5cf59e945d827b7f51d613cc690205c3414135bc09789ff40b', '\x7e16f8b545a3e06082e511c245870794', 'Demo', 'admin', NULL, '2021-11-23 14:35:19.359234', '2021-11-23 14:35:19.359234');

INSERT INTO channels
(id, organization_id, "name", "key", created_at, updated_at)
VALUES('4e3f9f5b-f63f-415c-853a-cdd633309455'::uuid, 'bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, 'Web', 'web', '2020-06-26 11:46:49.417', '2020-06-26 11:46:49.417');
INSERT INTO channels
(id, organization_id, "name", "key", created_at, updated_at)
VALUES('9044e340-4cb5-45d7-80f2-782083db9888'::uuid, 'bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, 'Store', 'store', '2020-07-03 06:46:44.995', '2020-07-03 06:46:44.995');
INSERT INTO channels
(id, organization_id, "name", "key", created_at, updated_at)
VALUES('9d430b8a-f7d6-4918-b560-203b15ff1945'::uuid, 'bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, 'App', 'app', '2020-07-03 06:46:23.368', '2020-07-03 06:46:23.368');

INSERT INTO auth_tokens (id, organization_id, token, channel, created_at) VALUES ('385686a1-5b19-476e-9e08-4b28ea5c7cb0', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'web-0IzIpVVc8oIU17MObozfrx.UWy02c7SPhFbzFU7mRZ.ZwDDRIJx.mu1bwOn3s-oa', 'web', '2021-11-23 14:59:29.096593');

INSERT INTO events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('3c11343c-b51a-4aef-9ac9-d8f2972fcb0c', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Variant Selected', 'The user has selected a product variant.', 0, '[{"id": "product", "name": "product", "properties": {"key": "product", "type": "string", "isRequired": true}}, {"id": "variant", "name": "variant", "properties": {"key": "variant", "type": "string", "isRequired": true}}, {"id": "attribute", "name": "attribute", "properties": {"key": "attribute", "type": "string", "isRequired": true}}]', '{"type": "object", "required": ["product", "variant", "attribute"], "properties": {"product": {"type": "string", "additionalItems": false, "additionalProperties": false}, "variant": {"type": "string", "additionalItems": false, "additionalProperties": false}, "attribute": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "product", "type": "string"}, {"key": "variant", "type": "string"}, {"key": "attribute", "type": "string"}]', '2020-07-02 08:04:18.014', '2021-02-05 14:31:46.916');
INSERT INTO events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('5d9e3c6f-3fe5-4148-a8f1-77db058eb7be', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Page Viewed', 'The user has views a page.', 0, '[{"id": "page", "name": "page", "properties": {"key": "page", "type": "string", "isRequired": true}}]', '{"type": "object", "required": ["page"], "properties": {"page": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "page", "type": "string"}]', '2020-07-01 07:38:22.204', '2020-07-01 07:38:22.204');
INSERT INTO events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('627f9036-dabd-4148-835e-fe0e284a7b23', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Checkout Started', 'The user has started the checkout process.', 0, '[{"id": "total", "name": "total", "properties": {"key": "total", "type": "float", "isRequired": true}}, {"id": "items", "name": "items", "properties": {"key": "items", "type": "integer", "isRequired": true}}]', '{"type": "object", "required": ["total", "items"], "properties": {"items": {"type": "integer", "additionalItems": false, "additionalProperties": false}, "total": {"type": "number", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "total", "type": "float"}, {"key": "items", "type": "integer"}]', '2020-07-01 12:09:56.082', '2020-07-01 12:09:56.082');
INSERT INTO events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('77c2382f-e9d8-4a2a-b0d0-d1dc3a45bd0a', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Product Viewed', 'The user has viewed a product.', 0, '[{"id": "product", "name": "product", "children": [], "properties": {"key": "product", "type": "string", "isRequired": true}}, {"id": "price", "name": "price", "properties": {"key": "price", "type": "float", "isRequired": true}}]', '{"type": "object", "required": ["product", "price"], "properties": {"price": {"type": "number", "additionalItems": false, "additionalProperties": false}, "product": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "product", "type": "string"}, {"key": "price", "type": "float"}]', '2020-06-30 08:54:43.643', '2020-06-30 08:59:15.896');
INSERT INTO events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('7adb3c3a-a1eb-4123-9d15-a6383cf6e8a0', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Searched', 'The user has search for a product.', 0, '[{"id": "query", "name": "query", "properties": {"key": "query", "type": "string", "isRequired": true}}]', '{"type": "object", "required": ["query"], "properties": {"query": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "query", "type": "string"}]', '2020-07-01 12:10:57.651', '2020-07-01 12:10:57.651');
INSERT INTO events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('7eeec417-7f1d-47b5-8548-af25ba50143d', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Checkout Completed', 'The user has completed the checkout.', 0, '[{"id": "total", "name": "total", "properties": {"key": "total", "type": "float", "isRequired": true}}, {"id": "items", "name": "items", "properties": {"key": "items", "type": "integer", "isRequired": true}}]', '{"type": "object", "required": ["total", "items"], "properties": {"items": {"type": "integer", "additionalItems": false, "additionalProperties": false}, "total": {"type": "number", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "total", "type": "float"}, {"key": "items", "type": "integer"}]', '2020-07-01 12:05:38.673', '2020-07-01 12:05:38.673');
INSERT INTO events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('872fafc4-29e1-4b05-8f59-72357843362b', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Category Viewed', 'The user has viewed a category page.', 0, '[{"id": "category", "name": "category", "properties": {"key": "category", "type": "string", "isRequired": true}}]', '{"type": "object", "required": ["category"], "properties": {"category": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "category", "type": "string"}]', '2020-06-30 08:54:14.623', '2020-06-30 08:54:14.623');
INSERT INTO events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('df6e1ca4-4c37-4ac0-8091-1cd010201233', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Field Focussed', 'The user has focussed a field in the checkout form.', 0, '[{"id": "field", "name": "field", "properties": {"key": "field", "type": "string", "isRequired": true}}, {"id": "empty", "name": "empty", "properties": {"key": "empty", "type": "boolean", "isRequired": true}}]', '{"type": "object", "required": ["field", "empty"], "properties": {"empty": {"type": "boolean", "additionalItems": false, "additionalProperties": false}, "field": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "field", "type": "string"}, {"key": "empty", "type": "boolean"}]', '2020-07-01 12:12:36.371', '2020-07-01 12:13:44.82');
INSERT INTO events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('fa29c00b-7c5a-4b1e-94ab-2e57d56d9ff2', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Product Displayed', 'A new product has been displayed', 0, '[{"id": "product_id", "name": "product_id", "properties": {"key": "product_id", "type": "string", "isRequired": true}}]', '{"type": "object", "required": ["product_id"], "properties": {"product_id": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "product_id", "type": "string"}]', '2021-01-14 08:47:30.453', '2021-01-14 08:47:30.453');
INSERT INTO events (id, organization_id, name, description, version, structure, json_schema, properties, created_at, updated_at) VALUES ('5b762e7f-600f-4f6a-abab-fdc22f80f9ac', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Updated Cart', 'A product has been added to or removed from the shopping cart.', 0, '[{"id": "product", "name": "product", "properties": {"key": "product", "type": "string", "isRequired": true}}, {"id": "price", "name": "price", "properties": {"key": "price", "type": "float", "isRequired": true}}, {"id": "quantity", "name": "quantity", "properties": {"key": "quantity", "type": "integer", "isRequired": true}}, {"id": "variant", "name": "variant", "children": [{"id": "size", "name": "size", "properties": {"key": "size", "type": "string", "isRequired": false}}, {"id": "color", "name": "color", "properties": {"key": "color", "type": "string"}}, {"id": "size", "name": "size", "properties": {"key": "size", "type": "string"}}], "properties": {"key": "variant", "type": "object"}}, {"id": "total", "name": "total", "properties": {"key": "total", "type": "float", "isRequired": true}}, {"id": "operation", "name": "operation", "properties": {"key": "operation", "type": "string", "isRequired": true}}]', '{"type": "object", "required": ["product", "price", "quantity", "total", "operation"], "properties": {"price": {"type": "number", "additionalItems": false, "additionalProperties": false}, "total": {"type": "number", "additionalItems": false, "additionalProperties": false}, "product": {"type": "string", "additionalItems": false, "additionalProperties": false}, "variant": {"type": "object", "properties": {"size": {"type": "string", "additionalItems": false, "additionalProperties": false}, "color": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}, "quantity": {"type": "integer", "additionalItems": false, "additionalProperties": false}, "operation": {"type": "string", "additionalItems": false, "additionalProperties": false}}, "additionalItems": false, "additionalProperties": false}', '[{"key": "product", "type": "string"}, {"key": "price", "type": "float"}, {"key": "quantity", "type": "integer"}, {"key": "variant.size", "type": "string"}, {"key": "variant.color", "type": "string"}, {"key": "variant.size", "type": "string"}, {"key": "total", "type": "float"}, {"key": "operation", "type": "string"}]', '2020-06-30 08:15:47.753', '2021-12-09 10:07:41.557362');

INSERT INTO actions
(id, organization_id, "name", "type", tags, state, properties, targeting, capping, hooks, test_users, created_at, updated_at)
VALUES('178af873-67e9-43ec-9e53-b1e831475b31'::uuid, 'bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, 'Not Available', 'not_available', '["improvement"]'::jsonb::jsonb, 'active', '[{"channels": ["web", "store", "app"], "name": "title", "sourceKey": "", "sourceType": "", "type": "text", "value": "Benachrichtige mich bei Verfügbarkeit"}, {"channels": ["web", "store", "app"], "name": "body", "sourceKey": "", "sourceType": "", "type": "text", "value": "Leider ist der ausgewählte Artikel nicht mehr verfügbar, aber wir bekommen diesen sicher wieder rein. Gerne informieren wir Dich dann per E-Mail."}, {"channels": ["web", "store", "app"], "name": "cta", "sourceKey": "", "sourceType": "", "type": "text", "value": "Benachrichtigung erhalten"}]'::jsonb::jsonb, '{"audiences": null, "channels": ["web", "store", "app"], "contextConditions": [], "end": {"date": null, "time": null}, "start": {"date": null, "time": null}, "traitConditions": []}'::jsonb::jsonb, '[{"channels": ["web", "store", "app"], "count": 2, "event": "hidden", "group": "user", "within": 0}, {"channels": ["store", "web", "app"], "count": 1, "event": "accepted", "group": "user", "within": 0}]'::jsonb::jsonb, '[]'::jsonb::jsonb, '[]'::jsonb::jsonb, '2020-07-13 07:37:51.100', '2020-07-13 13:23:00.734');
INSERT INTO actions
(id, organization_id, "name", "type", tags, state, properties, targeting, capping, hooks, test_users, created_at, updated_at)
VALUES('1b05b65c-cc97-43b5-9a2b-875e7f7d6da4'::uuid, 'bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, 'Incentivize Abandoned Cart Users', 'coupon_overlay', '["upselling"]'::jsonb::jsonb, 'active', '[{"channels": ["web"], "name": "title", "sourceKey": "", "sourceType": "", "type": "text", "value": "Schließe jetzt Deinen Kauf ab!!!!"}, {"channels": ["web"], "name": "body", "sourceKey": "", "sourceType": "", "type": "text", "value": "Wir dachten uns, dass Dir vielleicht dieser 10% Gutschein hilft, Deinen Kauf abzuschließen?"}, {"channels": ["web"], "name": "cta", "sourceKey": "", "sourceType": "", "type": "text", "value": "Gutschein jetzt einlösen"}, {"channels": ["web"], "name": "couponCode", "sourceKey": "", "sourceType": "", "type": "text", "value": "B4ZXQZVX"}]'::jsonb::jsonb, '{"audiences": ["80a45fce-2136-4954-a403-6894957a5ad4"], "channels": ["web"], "contextConditions": [], "end": {"date": null, "time": null}, "start": {"date": null, "time": null}, "traitConditions": null}'::jsonb::jsonb, '[{"channels": ["web"], "count": 1, "event": "accepted", "group": "user", "within": 0}, {"channels": ["web"], "count": 2, "event": "hidden", "group": "user", "within": 0}]'::jsonb::jsonb, '[]'::jsonb::jsonb, '[]'::jsonb::jsonb, '2020-07-02 08:42:05.258', '2021-04-08 08:17:53.349');
INSERT INTO actions
(id, organization_id, "name", "type", tags, state, properties, targeting, capping, hooks, test_users, created_at, updated_at)
VALUES('528cab3c-416c-48d1-9b59-614aec2df897'::uuid, 'bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, 'Preselect Size', 'size_preset', '["improvement"]'::jsonb::jsonb, 'active', '[{"channels": ["web"], "name": "size", "sourceKey": "favoriteSize", "sourceType": "", "type": "computed_trait", "value": ""}]'::jsonb::jsonb, '{"audiences": null, "channels": ["web"], "contextConditions": [{"key": "page.path", "operator": "starts_with", "type": "string", "value": "/?product="}], "end": {"date": null, "time": null}, "start": {"date": null, "time": null}, "traitConditions": [{"key": "favoriteSize", "operator": "exists", "source": "computed", "type": "string", "value": null}]}'::jsonb::jsonb, '[]'::jsonb::jsonb, '[]'::jsonb::jsonb, '[]'::jsonb::jsonb, '2020-07-07 07:10:58.013', '2020-07-07 08:11:09.677');
INSERT INTO actions
(id, organization_id, "name", "type", tags, state, properties, targeting, capping, hooks, test_users, created_at, updated_at)
VALUES('6cfc5ada-568f-4dc8-8d2e-ea06604c3ce8'::uuid, 'bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, 'Preselect category', 'favorite_category', '["improvement"]'::jsonb::jsonb, 'active', '[{"channels": ["web"], "name": "category", "sourceKey": "favoriteCategory", "sourceType": "", "type": "computed_trait", "value": ""}]'::jsonb::jsonb, '{"audiences": null, "channels": ["web"], "contextConditions": [{"key": "page.path", "operator": "starts_with", "type": "string", "value": "/?post_type=product"}], "end": {"date": null, "time": null}, "start": {"date": null, "time": null}, "traitConditions": [{"key": "favoriteCategory", "operator": "exists", "source": "computed", "type": "string", "value": null}]}'::jsonb::jsonb, '[]'::jsonb::jsonb, '[]'::jsonb::jsonb, '[]'::jsonb::jsonb, '2020-07-07 07:20:37.433', '2020-07-07 08:11:24.153');

INSERT INTO computed_traits
(id, organization_id, "name", "key", "type", event_id, conditions, properties, created_at, updated_at, refreshed_at)
VALUES('0e9b18b3-b7ed-459c-99cd-f6aab5f4d6ee'::uuid, 'bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, 'Last Activity', 'lastActivity', 'last_event', '5d9e3c6f-3fe5-4148-a8f1-77db058eb7be'::uuid, '[]'::jsonb::jsonb, '{"type": "dateTime", "useTimestamp": true}'::jsonb::jsonb, '2020-07-01 08:41:39.438', '2020-07-06 20:43:18.272', '2020-07-01 08:41:39.438');
INSERT INTO computed_traits
(id, organization_id, "name", "key", "type", event_id, conditions, properties, created_at, updated_at, refreshed_at)
VALUES('6ff4980c-e973-4521-bd7f-e147358ef6e8'::uuid, 'bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, 'Average Product Price', 'averageProductPrice', 'aggregation', '77c2382f-e9d8-4a2a-b0d0-d1dc3a45bd0a'::uuid, '[]'::jsonb::jsonb, '{"aggregationType": "avg", "property": "price", "type": "float"}'::jsonb::jsonb, '2020-06-30 21:04:47.685', '2020-07-06 20:43:28.062', '2020-06-30 21:04:47.685');
INSERT INTO computed_traits
(id, organization_id, "name", "key", "type", event_id, conditions, properties, created_at, updated_at, refreshed_at)
VALUES('87fecab9-b050-4013-8821-d9b26bd62fae'::uuid, 'bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, 'Favorite Category', 'favoriteCategory', 'most_frequent', '872fafc4-29e1-4b05-8f59-72357843362b'::uuid, '[]'::jsonb::jsonb, '{"minFrequency": 1, "property": "category", "type": "string"}'::jsonb::jsonb, '2020-06-30 09:00:27.240', '2021-02-05 13:57:20.980', '2020-06-30 09:00:27.240');
INSERT INTO computed_traits
(id, organization_id, "name", "key", "type", event_id, conditions, properties, created_at, updated_at, refreshed_at)
VALUES('912bb4ea-c551-4787-9876-71258456ad53'::uuid, 'bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, 'Total Product Value Viewed', 'totalValue', 'aggregation', '77c2382f-e9d8-4a2a-b0d0-d1dc3a45bd0a'::uuid, '[{"operator": "greater_than", "property": "price", "type": "float", "value": "1"}]'::jsonb::jsonb, '{"aggregationType": "sum", "property": "price", "type": "float"}'::jsonb::jsonb, '2021-01-25 13:31:31.232', '2021-01-25 13:31:31.232', '2021-01-25 13:31:31.232');
INSERT INTO computed_traits
(id, organization_id, name, key, type, event_id, conditions, properties, created_at, updated_at, refreshed_at)
VALUES ('9ad88d9f-c198-496a-a6ed-3a0f26e5aa4f', 'bc70b33d-c77f-4fe3-813d-a2605c0915cb', 'Favorite Size', 'favoriteSize', 'most_frequent', '3c11343c-b51a-4aef-9ac9-d8f2972fcb0c', '[{"type": "string", "value": "size", "operator": "equals", "property": "attribute"}]', '{"type": "string", "property": "variant", "minFrequency": 1}', '2020-07-02 08:07:04.1', '2021-12-09 17:56:51.090873', '2020-07-02 08:07:04.1');
INSERT INTO computed_traits
(id, organization_id, "name", "key", "type", event_id, conditions, properties, created_at, updated_at, refreshed_at)
VALUES('9c97cd3f-0d96-4612-bc18-d97e87c2abe6'::uuid, 'bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, 'Current Cart Value', 'currentCartValue', 'aggregation', '5b762e7f-600f-4f6a-abab-fdc22f80f9ac'::uuid, '[]'::jsonb::jsonb, '{"aggregationType": "sum", "property": "total", "type": "float"}'::jsonb::jsonb, '2020-07-01 09:01:25.464', '2020-07-06 20:43:54.406', '2020-07-01 09:01:25.464');
INSERT INTO computed_traits
(id, organization_id, "name", "key", "type", event_id, conditions, properties, created_at, updated_at, refreshed_at)
VALUES('b9d7d130-d817-40a5-9cc4-b265d5c6c021'::uuid, 'bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, 'Favorite Shoe Size', 'favoriteShoeSize', 'most_frequent', '3c11343c-b51a-4aef-9ac9-d8f2972fcb0c'::uuid, '[{"operator": "equals", "property": "attribute", "type": "string", "value": "schuhgroesse"}]'::jsonb::jsonb, '{"minFrequency": 1, "property": "variant", "type": "string"}'::jsonb::jsonb, '2020-07-02 08:24:26.812', '2021-02-05 13:56:32.277', '2020-07-02 08:24:26.812');
INSERT INTO computed_traits
(id, organization_id, "name", "key", "type", event_id, conditions, properties, created_at, updated_at, refreshed_at)
VALUES('c496f648-724b-455e-957f-8c6c77a7994c'::uuid, 'bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, 'First Page Viewed', 'firstPageVisited', 'first_event', '5d9e3c6f-3fe5-4148-a8f1-77db058eb7be'::uuid, '[]'::jsonb::jsonb, '{"property": "page", "type": "string", "useTimestamp": false}'::jsonb::jsonb, '2020-07-01 07:38:57.974', '2020-07-06 20:44:14.404', '2020-07-01 07:38:57.974');
INSERT INTO computed_traits
(id, organization_id, "name", "key", "type", event_id, conditions, properties, created_at, updated_at, refreshed_at)
VALUES('fb82775c-2f1a-4a19-8815-b81ec6de0f07'::uuid, 'bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, 'Last Page Viewed', 'lastPageViewed', 'last_event', '5d9e3c6f-3fe5-4148-a8f1-77db058eb7be'::uuid, '[]'::jsonb::jsonb, '{"property": "page", "type": "string"}'::jsonb::jsonb, '2020-07-01 07:39:28.358', '2020-07-06 20:44:26.447', '2020-07-01 07:39:28.358');

INSERT INTO custom_traits
(organization_id, "structure", json_schema, properties, created_at, updated_at)
VALUES('bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, '[{"id": "age", "name": "age", "properties": {"isRequired": true, "key": "age", "type": "integer"}}, {"id": "gender", "name": "gender", "properties": {"isRequired": true, "key": "gender", "type": "string"}}, {"id": "customerSince", "name": "customerSince", "properties": {"isRequired": true, "key": "customerSince", "type": "dateTime"}}]'::jsonb::jsonb, '{"additionalItems": false, "additionalProperties": false, "properties": {"age": {"additionalItems": false, "additionalProperties": false, "type": "integer"}, "customerSince": {"additionalItems": false, "additionalProperties": false, "type": "string"}, "gender": {"additionalItems": false, "additionalProperties": false, "type": "string"}}, "required": ["age", "gender", "customerSince"], "type": "object"}'::jsonb::jsonb, '[{"key": "age", "type": "integer"}, {"key": "gender", "type": "string"}, {"key": "customerSince", "type": "dateTime"}]'::jsonb::jsonb, '2020-07-05 19:24:53.302', '2020-07-15 07:43:56.114');

INSERT INTO audiences
(id, organization_id, "name", description, include_anonymous, events, traits, current_set_id, profile_count, created_at, updated_at, refreshed_at)
VALUES('7049b736-d2f2-4343-a1e4-caf6c6a325c8'::uuid, 'bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, 'Users With Product In Basket', 'Users that have a product in their basket.', true, '[{"count": 1, "id": "872fafc4-29e1-4b05-8f59-72357843362b", "internalId": "2e8e1663-8064-4847-ad77-53974d021f14", "operator": "more_or_exactly", "properties": [], "timeWindowOperator": "any_time"}]'::jsonb::jsonb, '[]'::jsonb::jsonb, '045a20a2-4593-4536-810c-936429d293ac'::uuid, 0, '2020-06-30 12:25:19.389', '2021-02-11 11:02:28.029', '2021-02-11 11:02:28.029');
INSERT INTO audiences
(id, organization_id, "name", description, include_anonymous, events, traits, current_set_id, profile_count, created_at, updated_at, refreshed_at)
VALUES('80a45fce-2136-4954-a403-6894957a5ad4'::uuid, 'bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, 'Checkout Not Completed', 'The user has not completed the checkout process.', true, '[{"count": 0, "id": "5b762e7f-600f-4f6a-abab-fdc22f80f9ac", "internalId": "fe8354e0-ef2d-4875-9fe0-31ab26b1dd32", "operator": "more_than", "properties": [], "timeWindowOperator": "any_time"}, {"count": 0, "exclude": false, "id": "627f9036-dabd-4148-835e-fe0e284a7b23", "internalId": "c0391e1e-2bf7-427b-b77e-d11002324071", "operator": "more_than", "parentId": "fe8354e0-ef2d-4875-9fe0-31ab26b1dd32", "properties": [], "timeWindowOperator": "any_time"}, {"count": 0, "exclude": true, "id": "7eeec417-7f1d-47b5-8548-af25ba50143d", "internalId": "d1d5c7b1-4d3d-4027-9df8-9c5b7b04cba4", "operator": "more_than", "parentId": "c0391e1e-2bf7-427b-b77e-d11002324071", "properties": [], "timeWindowOperator": "any_time"}]'::jsonb::jsonb, '[]'::jsonb::jsonb, 'c6b8cdf9-cacd-4005-8fe8-12523ad54979'::uuid, 1, '2020-07-01 12:18:47.049', '2021-11-05 09:20:52.351', '2021-11-05 09:20:52.351');

INSERT INTO contexts
(id, organization_id, channel, "structure", json_schema, properties, created_at, updated_at)
VALUES('963eb52a-0ecc-4b09-8d34-9ec8b9c4553e'::uuid, 'bc70b33d-c77f-4fe3-813d-a2605c0915cb'::uuid, 'web', '[{"children": [{"id": "path", "name": "path", "properties": {"isRequired": true, "key": "path", "type": "string"}}], "id": "page", "name": "page", "properties": {"isRequired": true, "key": "page", "type": "object"}}, {"id": "ua", "name": "ua", "properties": {"isRequired": true, "key": "ua", "type": "string"}}, {"id": "lang", "name": "lang", "properties": {"isRequired": true, "key": "lang", "type": "string"}}, {"id": "conn", "name": "conn", "properties": {"key": "conn", "type": "string"}}, {"children": [{"id": "w", "name": "w", "properties": {"isRequired": true, "key": "w", "type": "integer"}}, {"id": "h", "name": "h", "properties": {"isRequired": true, "key": "h", "type": "integer"}}], "id": "screen", "name": "screen", "properties": {"isRequired": true, "key": "screen", "type": "object"}}, {"children": [{"id": "w", "name": "w", "properties": {"isRequired": true, "key": "w", "type": "integer"}}, {"id": "h", "name": "h", "properties": {"isRequired": true, "key": "h", "type": "integer"}}], "id": "window", "name": "window", "properties": {"isRequired": true, "key": "window", "type": "object"}}]'::jsonb::jsonb, '{"additionalItems": false, "additionalProperties": false, "properties": {"conn": {"additionalItems": false, "additionalProperties": false, "type": "string"}, "lang": {"additionalItems": false, "additionalProperties": false, "type": "string"}, "page": {"additionalItems": false, "additionalProperties": false, "properties": {"path": {"additionalItems": false, "additionalProperties": false, "type": "string"}}, "required": ["path"], "type": "object"}, "screen": {"additionalItems": false, "additionalProperties": false, "properties": {"h": {"additionalItems": false, "additionalProperties": false, "type": "integer"}, "w": {"additionalItems": false, "additionalProperties": false, "type": "integer"}}, "required": ["w", "h"], "type": "object"}, "ua": {"additionalItems": false, "additionalProperties": false, "type": "string"}, "window": {"additionalItems": false, "additionalProperties": false, "properties": {"h": {"additionalItems": false, "additionalProperties": false, "type": "integer"}, "w": {"additionalItems": false, "additionalProperties": false, "type": "integer"}}, "required": ["w", "h"], "type": "object"}}, "required": ["page", "ua", "lang", "screen", "window"], "type": "object"}'::jsonb::jsonb, '[{"key": "page.path", "type": "string"}, {"key": "ua", "type": "string"}, {"key": "lang", "type": "string"}, {"key": "conn", "type": "string"}, {"key": "screen.w", "type": "integer"}, {"key": "screen.h", "type": "integer"}, {"key": "window.w", "type": "integer"}, {"key": "window.h", "type": "integer"}]'::jsonb::jsonb, '2020-06-30 08:51:08.706', '2020-06-30 08:51:08.706');
