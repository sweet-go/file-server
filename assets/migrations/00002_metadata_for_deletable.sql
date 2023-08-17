-- +migrate Up notransaction

ALTER TYPE delete_rule ADD VALUE 'SCHEDULED_DELETE';

ALTER TABLE "deletable_media" ADD COLUMN IF NOT EXISTS metadata JSONB default '{}'::jsonb;

-- +migrate Down

ALTER TABLE "deletable_media" DROP COLUMN IF EXISTS metadata;