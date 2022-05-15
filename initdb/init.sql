SELECT 'CREATE DATABASE store'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'store')\gexec