CREATE TABLE consent_rule (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    subject VARCHAR(255) NOT NULL,
    custodian VARCHAR(255) NOT NULL,
    actor VARCHAR(255) NOT NULL
);

CREATE UNIQUE INDEX uniq_consent_rule ON consent_rule(subject, custodian, actor);

CREATE TABLE resource (
    consent_rule_id INTEGER REFERENCES consent_rule(id),
    resource_type VARCHAR(255) NOT NULL
);

CREATE UNIQUE INDEX uniq_resource ON resource(consent_rule_id, resource_type);
