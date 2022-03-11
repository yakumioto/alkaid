erDiagram
    USER {
        string  resourceId
        string  userId
        string  name
        string  email
        string  password "基于 pbkdf2 和 hkdf 算法扩展密码，用于加密用户唯一的对称密钥"
        boolean root
        string  protectedSymmetricKey "使用用户扩展密码进行加密的对称密钥（系统生成）"
        string  protectedSignPrivateKey "使用用户的对称密钥加密签名私钥（系统生成）"
        string  signPublicKey
        string  protectedTlsPrivateKey "使用用户的对称密钥加密通讯私钥（系统生成）"
        string  tlsPublicKey
        string  protectedRSAPrivateKey "使用用户的对称密钥加密RSA私钥（系统生成，用于组织间对称密钥共享）"
        string  rsaPublicKey
        string  deactivate
        string  status
        int     createAt
        int     updateAt
    }
    ORGANIZATION {
        string resourceId
        string organizationId
        string name
        string domain
        string description
        string country
        string province
        string locality
        string organizationalUnit
        string streetAddress
        string postalCode
        string protectedSignPrivateKey "使用组织对称密钥加密的根签名证书（系统生成）"
        string signPublicKey
        string protectedTlsPrivateKey "使用组织对称密钥加密的根通讯证书（系统生成）"
        string tlsPublicKey
        string signCACertificate "组织签名根CA证书"
        string tlsCACertificate "组织通讯根CA证书"
        int    createAt
        int    updateAt
    }
    USER_ORGANIZATION {
        string  resourceId
        string  userId
        string  organizationId
        string  protectedSymmetricKey "使用用户RSA公钥进行加密的组织对称密钥"
        int     role
        string  status
        boolean deactivate
        int     createAt
        int     updateAt
        int     deactivateAt
    }
    IDENTITY {
        string  resourceId
        string  identityId
        string  organizationId
        string  use "用户使用还是节点使用"
        string  type "fabric中规定的 admin client orderer peer"
        string  description
        boolean nodeOUs
        array   sans
        string  signCertificate "组织签发的签名证书"
        string  tlsCertificate "组织签发的通讯证书"
        int     createAt
        int     updateAt
    }
    NETWORK {
        string resourceId
        string networkId
        string name
        string description
        string type
        string createdAt
        string updatedAt
    }
    NODE {
        string  resourceId
        string  nodeId
        string  organizationId
        string  networkId
        string  name
        string  type
        string  status
        string  protectedSignPrivateKey "使用组织对称密钥加密的签名私钥（系统生成）"
        string  signPublicKey
        string  protectedTlsPrivateKey "使用组织对称密钥加密的通讯私钥（系统生成）"
        string  tlsPublicKey
        int     createdAt
        int     updatedAt
    }

    USER }|--|{ USER_ORGANIZATION: "用户和组织一对多关系"
    ORGANIZATION }|--|{ USER_ORGANIZATION: "用户和组织一对多关系"
    IDENTITY ||--|{ NETWORK : "用户通过特定身份操作网络"
    USER ||--|{ IDENTITY : "用户拥有多个身份"
    NODE ||--|| IDENTITY : "节点拥有一个身份"
    NODE ||--|| NETWORK : "节点可以加入一个网络"
    ORGANIZATION }|--|{ NETWORK : "组织可以创建或加入多个网络"
    ORGANIZATION ||--|{ NODE : "组织拥有多个节点"