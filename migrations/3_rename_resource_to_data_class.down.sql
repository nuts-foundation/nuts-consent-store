CREATE TABLE resource (
      consent_record_id INTEGER,
      resource_type VARCHAR(255) NOT NULL,

      FOREIGN KEY(consent_record_id)
          REFERENCES consent_record (id)
          ON DELETE CASCADE
);

CREATE UNIQUE INDEX uniq_resource ON resource(consent_record_id, resource_type);

INSERT INTO resource (consent_record_id, resource_type) SELECT consent_record_id, code FROM data_class;

DROP INDEX uniq_data_class;

DROP TABLE data_class;