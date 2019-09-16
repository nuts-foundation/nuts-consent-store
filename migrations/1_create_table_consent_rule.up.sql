CREATE TABLE patient_consent (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    subject VARCHAR(255) NOT NULL,
    custodian VARCHAR(255) NOT NULL,
    actor VARCHAR(255) NOT NULL
);

CREATE UNIQUE INDEX uniq_patient_consent ON patient_consent(subject, custodian, actor);

CREATE TABLE consent_record (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    patient_consent_id INTEGER REFERENCES patient_consent(id),
    valid_from DATE NOT NULL,
    valid_to DATE NOT NULL,
    proof_hash VARCHAR(255) NOT NULL UNIQUE
);

CREATE UNIQUE INDEX uniq_record ON consent_record(proof_hash);

CREATE TABLE resource (
    consent_record_id INTEGER,
    resource_type VARCHAR(255) NOT NULL,

    FOREIGN KEY(consent_record_id)
    REFERENCES consent_record (id)
    ON DELETE CASCADE
);

CREATE UNIQUE INDEX uniq_resource ON resource(consent_record_id, resource_type);
