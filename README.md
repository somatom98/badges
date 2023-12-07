# Badges

Simple Employees access control system capable of:
- Ingesting employee information about badge-in and badge-out events
- Letting single employees access their own ingested data
    - Badge events areee aggregated by year, month and day
    - The employee is able to receive the updated aggregation in case new badge events are received
- Letting managers access both their own ingested data and the data of her/his direct reports

## Events
Badge events are intended as simple structured data events containing:
- a reference to the employee,
- the kind of badge event (in or out)
- when it occurred by means of a timestamp

## Architecture

### API

GraphQL was used. There are two queries and one subscription available:
- query `events`: returns all the events related to a user id passed as parameter
    ```graphql
    query GetEventsList {
        events(id:"5f05a5adad8d849ccafee4b1") {
            years {
            year
                months {
                    month
                    days {
                    day
                        events {
                            id
                            uid
                            type
                            date
                        }
                    }
                }
            }
        }
    }
    ```
- query `reportsEvents`: returns all the events related to the users maneged by the manager id passed as parameter
    ```graphql
    query GetManagerEventsList {
        reportsEvents(mid:"5f05a5adad8d849ccafee4b1") {
            years {
                year
                months {
                    month
                    days {
                    day
                    events {
                        id
                        uid
                        type
                        date
                    }
                    }
                }
            }
        }
    }
    ```
- subscription `events`: returns a stream of the events with a userid equals to the one passed as parameter
    ```graphql
    subscription SubscribeToEvents{
        events(id: "5f05a5adad8d849ccafee4b3") {
            id
            uid
            type
            date
        }
    }
    ```

### Events ingestion

Events ingestion is implemented using Kafka. This allows to have a reliable system wich allows fault tolerance to both publisher and consumer sides.

One consumer group is always consuming events and saving them on DB and always starts reading from the last saved event, while other ones get instantiated at each new stream, starting to read from the last published message.

### Persistence

Mongo DB was chosen as database as there were not complicated relations.

Two collections are used:
- `users`, containing documents with the following fields:
  - _id: objectID
  - mid: objectID (nullable) - used to reference the _id field of the manager user
- `events`, containing documents with the following fields:
  - _id: objectID
  - id: string - used to reference the id wich identifies the event externally (a badge machine for example)
  - uid: objectID - reference to the _id field of the user
  - type: string (IN | OUT)
  - ts: date

## Setup

The project can be executed with `go run ./cmd`.

Configurations can be set using the `config.yaml` file. In particular, mongo connection string and kafka need to be setted up.

To simulate badge events, the script `scripts/kafka_publisher.go` can be used.
