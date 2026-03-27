ALTER TABLE groups
  ADD COLUMN IF NOT EXISTS default_proxy_id BIGINT;

CREATE INDEX IF NOT EXISTS idx_groups_default_proxy_id
  ON groups (default_proxy_id);

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'groups_default_proxy_id_fkey'
  ) THEN
    ALTER TABLE groups
      ADD CONSTRAINT groups_default_proxy_id_fkey
      FOREIGN KEY (default_proxy_id)
      REFERENCES proxies(id)
      ON DELETE SET NULL;
  END IF;
END $$;
