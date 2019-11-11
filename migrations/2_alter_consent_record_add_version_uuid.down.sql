DROP INDEX uniq_record_version;

ALTER TABLE consent_record RENAME TO consent_record_tmp;

CREATE TABLE consent_record (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    patient_consent_id VARCHAR(255) REFERENCES patient_consent(id),
    valid_from DATE NOT NULL,
    valid_to DATE NOT NULL,
    hash VARCHAR(255) NOT NULL UNIQUE
);

INSERT INTO consent_record SELECT (id, patient_consent_id, valid_from, valid_to, hash) FROM consent_record_tmp;

DROP TABLE consent_record_tmp;
