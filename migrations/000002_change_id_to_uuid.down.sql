BEGIN;

ALTER TABLE event
DROP CONSTRAINT event_pkey;

ALTER TABLE event
ALTER COLUMN id TYPE INTEGER USING (NULL);

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT FROM pg_sequences WHERE sequencename = 'event_id_seq') THEN
        CREATE SEQUENCE event_id_seq;
    END IF;
END $$;

ALTER TABLE event
ALTER COLUMN id SET DEFAULT nextval('event_id_seq');

SELECT setval('event_id_seq', COALESCE((SELECT MAX(id) FROM event), 1));

ALTER TABLE event
ADD CONSTRAINT event_pkey PRIMARY KEY (id);

COMMIT;