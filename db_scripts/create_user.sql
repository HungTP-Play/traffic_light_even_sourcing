CREATE USER event_store_usr WITH ENCRYPTED PASSWORD 'event_store_pass';
GRANT ALL ON DATABASE event_store TO event_store_usr;
GRANT ALL ON SCHEMA public TO event_store_usr;
GRANT CREATE ON SCHEMA public TO event_store_usr; 
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO event_store_usr;
GRANT INSERT, SELECT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO event_store_usr;

CREATE USER projector_usr WITH ENCRYPTED PASSWORD 'projector_pass'; 
GRANT ALL ON DATABASE projections TO projector_usr;
GRANT ALL ON SCHEMA public TO projector_usr;
GRANT CREATE ON SCHEMA public TO projector_usr;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO projector_usr;
GRANT INSERT, SELECT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO projector_usr;

CREATE USER controller_usr WITH ENCRYPTED PASSWORD 'controller_pass';
GRANT ALL ON DATABASE controller TO controller_usr;
GRANT ALL ON SCHEMA public TO controller_usr;
GRANT CREATE ON SCHEMA public TO controller_usr;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO controller_usr;
GRANT INSERT, SELECT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO controller_usr;