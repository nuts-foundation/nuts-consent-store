CREATE TABLE patient_consent (
    id VARCHAR(255) PRIMARY KEY,
    subject VARCHAR(255) NOT NULL,
    custodian VARCHAR(255) NOT NULL,
    actor VARCHAR(255) NOT NULL
);

CREATE UNIQUE INDEX uniq_patient_consent ON patient_consent(actor, subject, custodian);
CREATE INDEX idx_patient_consent_custodian ON patient_consent(custodian);

CREATE TABLE consent_record (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    patient_consent_id VARCHAR(255) REFERENCES patient_consent(id),
    valid_from DATE NOT NULL,
    valid_to DATE NOT NULL,
    hash VARCHAR(255) NOT NULL UNIQUE
);

CREATE UNIQUE INDEX uniq_record ON consent_record(hash);

CREATE TABLE resource (
    consent_record_id INTEGER,
    resource_type VARCHAR(255) NOT NULL,

    FOREIGN KEY(consent_record_id)
    REFERENCES consent_record (id)
    ON DELETE CASCADE
);

CREATE UNIQUE INDEX uniq_resource ON resource(consent_record_id, resource_type);
