ALTER TABLE consent_record ADD COLUMN version INTEGER DEFAULT 1;
ALTER TABLE consent_record ADD COLUMN uuid VARCHAR(255);

CREATE UNIQUE INDEX uniq_record_version ON consent_record(patient_consent_id, uuid, version);