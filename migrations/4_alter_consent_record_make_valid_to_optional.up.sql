DROP INDEX uniq_record_version;

ALTER TABLE consent_record RENAME TO consent_record_tmp;

CREATE TABLE consent_record (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    patient_consent_id VARCHAR(255) REFERENCES patient_consent(id),
    valid_from DATE NOT NULL,
    valid_to DATE NULL,
    hash VARCHAR(255) NOT NULL UNIQUE,
    version INTEGER DEFAULT 1,
    uuid VARCHAR(255),
    previous_hash VARCHAR(255)
);

CREATE UNIQUE INDEX uniq_record_version ON consent_record(patient_consent_id, uuid, version);

INSERT INTO consent_record SELECT id, patient_consent_id, valid_from, valid_to, hash, version, uuid, previous_hash FROM consent_record_tmp;

DROP TABLE consent_record_tmp;
