CREATE USER event_store_usr WITH ENCRYPTED PASSWORD 'event_store_pass';
GRANT ALL ON DATABASE event_store TO event_store_usr;

CREATE USER projector_usr WITH ENCRYPTED PASSWORD 'projector_pass'; 
GRANT ALL ON DATABASE projections TO projector_usr;

CREATE USER light_metadata_usr WITH ENCRYPTED PASSWORD 'light_metadata_pass';
GRANT ALL ON DATABASE light_metadata TO light_metadata_usr;