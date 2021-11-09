CREATE DATABASE attractify ON CLUSTER attractify;

CREATE TABLE attractify.events_local ON CLUSTER attractify (
  id UUID DEFAULT generateUUIDv4(),
  organization_id UUID,
  identity_id UUID,
  event_id UUID,
  channel String,
  context String,
  properties String,
  created_at DateTime('UTC')
) ENGINE = ReplicatedMergeTree('/clickhouse/tables/{shard}/events_local', '{replica}') PARTITION BY (toYYYYMM(created_at)) ORDER BY (organization_id, channel, event_id, created_at) SETTINGS index_granularity = 8192;
CREATE TABLE attractify.events ON CLUSTER 'attractify' AS attractify.events_local ENGINE = Distributed('attractify', 'attractify', 'events_local', cityHash64(id, organization_id));

CREATE TABLE attractify.reactions_local ON CLUSTER attractify (
  id UUID DEFAULT generateUUIDv4(),
  organization_id UUID,
  action_id UUID,
  identity_id UUID,
  channel String,
  event String,
  context String,
  properties String,
  result String,
  created_at DateTime('UTC')
) ENGINE = ReplicatedMergeTree('/clickhouse/tables/{shard}/reactions_local', '{replica}') PARTITION BY (toYYYYMM(created_at)) ORDER BY (organization_id, action_id, identity_id, channel, event) SETTINGS index_granularity = 8192;
CREATE TABLE attractify.reactions ON CLUSTER 'attractify' AS attractify.reactions_local ENGINE = Distributed('attractify', 'attractify', 'reactions_local', cityHash64(id, organization_id));

CREATE TABLE attractify.identities ON CLUSTER attractify (
  id UUID,
  organization_id UUID,
  profile_id UUID,
  channel String,
  type String,
  user_id String,
  is_anonymous Boolean,
  created_at DateTime
) ENGINE = ODBC('DSN=attractify', 'attractify', 'profile_identities');

CREATE TABLE attractify.full_identities ON CLUSTER attractify (
  id UUID,
  profile_id UUID,
  organization_id UUID,
  channel String,
  type String,
  user_id String,
  is_anonymous UInt8,
  custom_traits String,
  computed_traits String,
  created_at DateTime
) ENGINE = ODBC('DSN=attractify', 'attractify', 'full_identities');

