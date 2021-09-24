```mermaid
erDiagram
    USER {
        string resourceId
        string id
        string name
        string email
        string password
        string role
        string protectedSigPrivateKey
        string protectedTlsPrivateKey
        string status
        int64  createAt
        int64  updateAt
    }
    ORGANIZATION {
        string id
        string name
        string domain
        string description
        string type
        string country
        string province
        string locality
        string organizationalUnit
        string streetAddress
        string postalCode
        string signCAPrivateKey
        string tlsCAPrivateKey
        string signCACertificate
        string tlsCACertificate
        int64  createAt
        int64  updateAt
    }
    IDENTITY {
        string  id
        string  organizationId
        string  use
        string  type
        string  description
        boolean nodeOUs
        array   sans
        string  signCertificate
        string  tlsCertificate
        int64   createAt
        int64   updateAt
    }
    NETWORK {
        string id
        string name
        string description
        string type
        string createdAt
        string updatedAt
    }
    NODE {
        string  id
        string  organizationId
        string  networkId
        string  name
        boolean enableCouchDB
        string  status
        string signPrivateKey
        string tlsPrivateKey
        string createdAt
        string updatedAt
    }

    IDENTITY ||--|{ NETWORK : "操作"
    USER }|--|{ ORGANIZATION : "属于"
    USER ||--|{ IDENTITY : "拥有"
    NODE ||--|| IDENTITY : "拥有"
    NODE ||--|| NETWORK : "加入"
    ORGANIZATION }|--|{ NETWORK : "加入"
    ORGANIZATION ||--|{ NODE : "拥有"
```