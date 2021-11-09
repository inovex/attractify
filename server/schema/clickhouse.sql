CREATE DATABASE IF NOT EXISTS attractify;

CREATE TABLE IF NOT EXISTS attractify.events (
  id UUID DEFAULT generateUUIDv4(),
  organization_id UUID,
  identity_id UUID,
  event_id UUID,
  channel String,
  context String,
  properties String,
  created_at DateTime('UTC')
) ENGINE = MergeTree ORDER BY (organization_id, channel, event_id, created_at) SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS attractify.reactions (
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
) ENGINE = MergeTree ORDER BY (organization_id, action_id, identity_id, channel, event) SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS attractify.identities (
  id UUID,
  organization_id UUID,
  profile_id UUID,
  channel String,
  type String,
  user_id String,
  is_anonymous Boolean,
  created_at DateTime
) ENGINE = ODBC('DSN=attractify', 'attractify', 'profile_identities');

CREATE TABLE IF NOT EXISTS attractify.full_identities (
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
