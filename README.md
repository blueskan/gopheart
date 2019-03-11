# gopheart
Simple healthchecking solution for many providers

Supported Providers:
  - Url Provider [ COMPLETED ]
  - Redis Provider [ COMPLETED ]
  - Memcache Provider
  - MongoDB Provider
  - MySQL Provider
  - PostgreSQL Provider
  - MSSQL Provider
  - Cassandra Provider
  - CouchBase Provider
  - HBase Provider
  - Neo4J Provider
  - ElasticSearch Provider
  - Solr Provider
  - RabbitMQ Provider
  - Kafka Provider

Supported Notifiers:
  - Slack Notifier
  - Mail Notifier
  
Features:
  - Easy configuration JSON or YAML
  - Smart scheduling between health checks
  - Audit Logging [ Optionally with stats persistence option ]
  - Creates Json Response automatically with provider infos ( If any error happened in any service you would return specific http status code, you can bind this url your health check system )
  - If you open stats option in configuration, /stats endpoint available for downtime stats
