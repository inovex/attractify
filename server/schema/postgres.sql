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

CREATE TABLE IF NOT EXISTS invalid_events (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	organization_id uuid NOT NULL REFERENCES organizations(id) ON DELETE CASCADE ON UPDATE CASCADE,
	name text NOT NULL,
	properties jsonb NOT NULL DEFAULT '[]'::jsonb,
	context jsonb NOT NULL DEFAULT '[]'::jsonb,
	error text NOT NULL,
	type text NOT NULL,
	created_at timestamp NOT NULL DEFAULT now()
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
