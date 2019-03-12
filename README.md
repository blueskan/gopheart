# gopheart
Simple healthchecking solution for many providers

Supported Providers:
  - Url Provider [ COMPLETED ]
  - Redis Provider [ COMPLETED ]
  - Memcache Provider [ COMPLETED ]
  - MongoDB Provider [ COMPLETED ]
  - MySQL Provider [ COMPLETED ]
  - PostgreSQL Provider [ COMPLETED ]
  - MSSQL Provider [ COMPLETED ]
  - Cassandra Provider [ COMPLETED ]
  - CouchBase Provider [ COMPLETED ]
  - HBase Provider [ COMPLETED - NEED SOME CONTROL ]
  - Neo4J Provider [ TODO ]
  - ElasticSearch Provider [ TODO ]
  - Solr Provider [ TODO ]
  - RabbitMQ Provider [ TODO ]
  - Kafka Provider [ TODO ]

Supported Notifiers:
  - Slack Notifier [ COMPLETED ]
  - Mail Notifier [ TODO ]
  
Features:
  - Easy configuration JSON or YAML
  - Smart scheduling between health checks
  - Audit Logging [ Optionally with stats persistence option ]
  - Creates Json Response automatically with provider infos ( If any error happened in any service you would return specific http status code, you can bind this url to your health check system )
  - If you open `stats option` in configuration, /stats endpoint available for downtime stats

TODOs:
  - We should divide Scheduling and Persistent logic :(
  - We should provide tests
  - Notifier Supports [ IN PROGRESS ]

This project in under development, more user friendly instructions coming soon..
