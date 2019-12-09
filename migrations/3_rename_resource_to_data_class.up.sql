CREATE TABLE data_class (
      consent_record_id INTEGER,
      code VARCHAR(255) NOT NULL,

      FOREIGN KEY(consent_record_id)
          REFERENCES consent_record (id)
          ON DELETE CASCADE
);

CREATE UNIQUE INDEX uniq_data_class ON data_class(consent_record_id, code);

INSERT INTO data_class (consent_record_id, code) SELECT resource.consent_record_id, resource.resource_type FROM resource;

DROP INDEX uniq_resource;

DROP TABLE resource;