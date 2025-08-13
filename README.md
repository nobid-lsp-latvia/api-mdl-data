# api-mdl-data

API for receiving mobile Drivers licence data from "Ceļu Satiksmes Drošības direkcijas"

## Abbreviations

| Abbreviation | Description                       |
| ------------ | --------------------------------- |
| mdl          | Mobile Drivers license            |
| CSDD         | Ceļu Satiksmes Drošības direkcija |

## Request for mdl data from CSDD

### Request

```bash
GET {host}/mdl
```

Personas kods atnāks header parametrā iekodēts Bear
Tālāk idAuth to mācēs izņemt laukā.

### Nepieciešami šādi ENV parametri

```bash
    ENVIRONMENT: "production"
    BASE_PATH: "/mdl"
    REVERSE_PROXY_TRUSTED_IPS: "*"
    REVERSE_PROXY_LIMIT: "3"

    IDAUTH_URL: ""
    IDAUTH_CLIENT_ID: ""
    IDAUTH_CLIENT_SECRET_FILE: /secret/edim-idauth-client-secret-api-mdl-data

    VAULT_LOGIN_URL: "https://vault.example.lv/v1/auth/lvrtc-edim/login"
    VAULT_DATA_URL: "https://vault.example.lv/v1/secrets-v2/data/lvrtc/edim/csdd/dev/edim-csdd-service-password"
    VAULT_ROLE_ID: ""
    VAULT_SECRET_ID_FILE: /secret/edim-api-mdl-data-vault-secret

    CSDD_URL: "https://example.lv"
    CSDD_USERNAME: "test-user"
    CSDD_CHANGE_PASSWORD_DAYS: "10"
    CSDD_SKIP_TLS_VERIFY: "true"
    CSDD_SYSTEM_GUID: "AAA-BBBB-CCCCC-DDDDDDDD"
    CSDD_SYSTEM_NAME: "TEST"
```

| Variable | Value | Description |
|----------|-------|-------------|
| **Application Configuration** | | |
| `ENVIRONMENT` | "production" | Specifies the deployment environment (production) |
| `BASE_PATH` | "/mdl" | Base path for the application routing |
| `REVERSE_PROXY_TRUSTED_IPS` | "*" | Defines trusted IP addresses for reverse proxy (all IPs allowed) |
| `REVERSE_PROXY_LIMIT` | "3" | Maximum number of reverse proxy hops to trust |
| **Authentication (IDAUTH)** | | |
| `IDAUTH_URL` | "" | URL for IDAuth service (empty/not configured) |
| `IDAUTH_CLIENT_ID` | "" | `api-mdl-data` id registrated in idAuth service |
| `IDAUTH_CLIENT_SECRET_FILE` | "/secret/edim-idauth-client-secret-api-mdl-data" | Path to the file containing the client secret for authentication |
| **Vault Configuration** | | |
| `VAULT_LOGIN_URL` | "https://vault.example.lv/v1/auth/lvrtc-edim/login" | URL for Vault authentication login |
| `VAULT_DATA_URL` | "https://vault.example.lv/v1/secrets-v2/data/lvrtc/edim/csdd/dev/edim-csdd-service-password" | URL for retrieving secret from Vault |
| `VAULT_ROLE_ID` | "" | Vault role ID |
| `VAULT_SECRET_ID_FILE` | "/secret/edim-api-mdl-data-vault-secret" | Path to the file containing Vault secret ID |
| **CSDD (Central Traffic Register) Configuration** | | |
| `CSDD_URL` | "" | Endpoint URL for CSDD api. SHALL BE FQDN (register internal) |
| `CSDD_USERNAME` | "" | Username for CSDD api access |
| `CSDD_CHANGE_PASSWORD_DAYS` | "10" | Number of days after which password should be changed |
| `CSDD_SKIP_TLS_VERIFY` | "true" | Indicates whether to skip TLS certificate verification |
| `CSDD_SYSTEM_GUID` | "" | Unique identifier issued by CSDD. Check password change documentation. |
| `CSDD_SYSTEM_NAME` | "" | System name for CSDD integration. Check password change documentation. |

### Response

JSON object

```json
{
  "document_number": "tstr",
  "birth_date": "full-date",
  "given_name": "tstr",
  "family_name": "tstr",
  "issue_date": "tdate or full-date",
  "expiry_date": "tdate or full-date",
  "issuing_country": "tstr",
  "issuing_authority": "tstr",
  "un_distinguishing_sign": "tstr",
  "portrait": "bstr",
    "driving_privileges": [
      {
        "vehcile_category_code": "tstr",
        "issue_date": "full-date",
        "expiry_date": "full-date",
        "code": [
                {
                  "sign": "tstr",
                  "value": "tstr"
                }
            ]
       }
    ]
}
```

#### Description of atributes

#### 2.3.2 Overview

| **Attribute identifier** | **Definition** | **Presence** | **Encoding format** |
| --- | --- | --- | --- |
| `family_name` | Current last name(s) or surname(s) of the mdl holder. | `M` | tstr |
| `given_name` | Current first name(s), including middle name(s), of the mdl holder. | `M` | tstr |
| `birth_date` | Day, month, and year on which the mdl holder was born. | `M` | full-date |
| `birth_place` | The country, state (or where applicable province, district, or local area), and city where the mdl holder was born | `M` | full-date |
| `issue_date` | Date when the mdl was issued. | `M` | tdate or full-date |
| `expiry_date` | Date when the mdl will expire. | `M` | tdate or full-date |
| `issuing_country` | Alpha-2 country code, as defined in ISO 3166-1, of the mdl Provider's country or territory. | `M` | tstr |
| `issuing_authority` | Name of the administrative authority that has issued this mdl instance, or the ISO 3166 Alpha-2 country code of the respective Member State if there is no separate authority authorized to issue mdls. | `M` | tstr |
| `document_number` | A number for the mdl, assigned by the mdl Provider. | `O` | tstr |
| `portrait` | Portrait of mdl holder | `O` | tstr |
| `signature_usual_mark` | Image of signature of the mdl holder | `O` | bstr |
| `personal_administrative_number` | Personas kods. A value assigned to the natural person that is unique among all personal administrative numbers issued by the provider of person identification data. | `M` | tstr |
| `driving_privileges` | The country where the mdl holder currently resides, as an Alpha-2 country code as specified in ISO 3166-1. | `O` | tstr |

##### Encoding reqirements

- `tstr`, `uint`, `bstr`, `bool` and `tdate` are CDDL representation types defined in [RFC 8610](https://www.rfc-editor.org/rfc/rfc8610.html).
- All attributes having encoding format tstr SHALL have a maximum length of 150 characters
- This document specifies `full-date` as `full-date` = #6.1004(tstr), where tag 1004 is specified in [RFC 8943](https://datatracker.ietf.org/doc/html/rfc8943)
- In accordance with [RFC 8949], Section 3.4.1, a `tdate` attribute shall contain a `date-time` string as specified in [RFC 3339]. In accordance with [RFC 8943], a `full-date` attribute shall contain a `full-date` string as specified in [RFC 3339].
- The following requirements SHALL apply to the representation of dates in attributes, unless otherwise indicated:
  - Fractions of seconds **SHALL NOT** be used;
  - A local offset from UTC SHALL NOT be used; the time-offset defined in [RFC 3339] SHALL be to "Z".
