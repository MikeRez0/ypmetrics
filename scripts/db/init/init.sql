CREATE USER metrics
    PASSWORD 'metrics';

CREATE DATABASE metrics_db
    OWNER 'metrics'
    ENCODING 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8';