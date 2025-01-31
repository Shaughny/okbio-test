# **Agents Monitoring API** 

## **Overview**
The **Agents Monitoring API** is a Go-based REST API that allows agents to send their **IP addresses**, retrieve IP-related details (such as **ASN** and **ISP**), and manage agent information. The API uses **SQLite** for storage and integrates with an external IP lookup service (`ip-api.com`).

---

## ** Features**
- Accepts IP address submissions from agents.  
- Retrieves and stores **Autonomous System Number (ASN)** and **Internet Service Provider (ISP)** details.  
- Provides endpoints to **list all agents** and **fetch details of a specific agent**.   
- Supports **Docker** and `.env` configuration.

---
## API Endpoints

| Method | Endpoint        | Description |
|--------|----------------|-------------|
| `POST` | `/agents`      | Register an agent's IP address |
| `GET`  | `/agents`      | Get a list of all registered agents |
| `GET`  | `/agents/{id}` | Get details (ASN, ISP) of a specific agent |




### 🔹 **Example: Register an Agent**
#### **Request:**
```sh
curl -X POST "http://localhost:8080/agents" \
     -H "Content-Type: application/json" \
     -d '{"ip_address": "8.8.8.8"}'

```
#### **Response:**
```json
{
  "id": 1,
  "ip_address": "8.8.8.8",
  "asn": "15169",
  "isp": "Google LLC"
}
```

### 🔹 **Example: Get a List of Agents**
#### **Request:**
```sh
curl -X GET "http://localhost:8080/agents"
```
#### **Response:**
```json
[
  {
    "id": 1,
    "ip_address": "8.8.8.8",
},
  {
    "id": 2,
    "ip_address": "1.1,1.1",
}
]
```


### 🔹 **Example: Get Details of a Specific Agent**
#### **Request:**
```sh
curl -X GET "http://localhost:8080/agents/1"
```
#### **Response:**
```json
{
  "id": 1,
  "ip_address": "8.8.8.8",
    "asn": "15169",
    "isp": "Google LLC"
}
```

---
## ** Running the API**
### **🔹 Using Docker**
1. **Clone the repository:**
```sh
git clone https://github.com/Shaughny/okbio-test.git
```

2. **Start with docker:**
```sh
make docker-run
```
3. **Read the logs:**
```sh
make logs
```

4.**Stop the container:**
```sh
make docker-stop
```

### **🔹 Using Local Go Installation**
1. **Clone the repository:**
```sh
git clone https://github.com/Shaughny/okbio-test.git
```

2. **Make sure go is installed on your machine.**
3. **Run the API:**
```sh
make run
```

---
### **🔹 Running Tests**
```sh
make test
```


### **🔹 All make commands**
```sh
make help
```

### **🔹 Docker Without using make**
1. **Clone the repository:**
```sh 
docker compose up -d --build
```
2. **Stop the container:**
```sh
docker compose down
```

