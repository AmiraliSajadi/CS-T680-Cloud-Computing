# Final Project

## <mark style="background: #5332A0A6;">Description</mark>
This project implements a number of containerized APIs that interact with one another and a Redis cache container:
```
               ┏━━━━━━━━━━━┓               
      ┌───────▶┃ Votes API ┃◀────────┐     
      ▼        ┗━━━━━━━━━━━┛         ▼     
┌───────────┐        │         ┌──────────┐
│ Voter API │        │         │ Poll API │
└───────────┘        │         └──────────┘
      │              │                │    
      ▼              ▼                ▼    
┌─────────────────────────────────────────┐
│              Cache (Redis)              │
└─────────────────────────────────────────┘
```

- The Poll API allows us to register new polls.
- The Voter API allows us to create new voters without prior votes.
- The Votes API allows us to create new votes when provided with existing voters and polls

## <mark style="background: #5332A0A6;">Overall Structure</mark>

**project-directory/**</br>
|-- poll-api/</br>
|         &emsp;&emsp;|-- Dockerfile</br>
|         &emsp;&emsp;|-- build-docker.sh</br>
|         &emsp;&emsp;|-- ...</br>
|-- voter-api/</br>
|         &emsp;&emsp;|-- Dockerfile</br>
|         &emsp;&emsp;|-- build-docker.sh</br>
|         &emsp;&emsp;|-- ...</br>
|-- votes-api/ </br>
|         &emsp;&emsp;|-- Dockerfile</br>
|         &emsp;&emsp;|-- build-docker.sh</br>
|         &emsp;&emsp;|-- ...</br>
|-- docker-compose.yml</br>
|-- do_the_thing.sh</br>
|-- ...


The three APIs are available in the three following folders:
- poll-api
- voter-api
- votes-api


In each API folder there is a *Dockerfile* and a *build-docker.sh* script which builds the API and then the container. There is also a docker-compose.yaml file in the root directory, used for configuring and running all the containers.

## <mark style="background: #5332A0A6;">Run It:</mark>
The easiest way to run the containerized APIs along with the Redis container is to use the provided script *do_the_thing.sh*. The script will use the *curl* tool (already added to the containers via Dockerfile) to send *http* requests to insert some sample data.


## <mark style="background: #5332A0A6;">Make Changes</mark>
If you need to make changes to any of the three APIs all you need to do afterward is to run:
```
./build-docker.sh
```
inside the corresponding API directory. As mentioned before, this script will rebuild the go project and our container to make sure all the changes will be reflected in the container.