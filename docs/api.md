### /categories
---
##### ***GET***
**Summary:** List menu categories.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 400 | MenuCategory | [ [MenuCategory](#menucategory) ] |
| 404 | GenericError | [GenericError](#genericerror) |

### /categories/{category_id}
---
##### ***GET***
**Summary:** List menu items with specified category.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | MenuItem | [ [MenuItem](#menuitem) ] |
| 500 | GenericError | [GenericError](#genericerror) |

### /categories/{id}
---
##### ***PUT***
**Summary:** Update menu category.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | MenuCategory | [MenuCategory](#menucategory) |
| 400 | GenericError | [GenericError](#genericerror) |

##### ***POST***
**Summary:** Create menu category.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | MenuCategory | [MenuCategory](#menucategory) |
| 400 | GenericError | [GenericError](#genericerror) |

### /confirm/{code}
---
##### ***GET***
**Summary:** Confirm reservation.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | State | [State](#state) |
| 400 | GenericError | [GenericError](#genericerror) |

### /menu
---
##### ***GET***
**Summary:** List all menu items.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | MenuItem | [ [MenuItem](#menuitem) ] |
| 500 | GenericError | [GenericError](#genericerror) |

##### ***POST***
**Summary:** Create menu item.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | MenuItem | [MenuItem](#menuitem) |
| 400 | GenericError | [GenericError](#genericerror) |

### /menu/{id}
---
##### ***GET***
**Summary:** Returns menu item.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | MenuItem | [MenuItem](#menuitem) |
| 404 | GenericError | [GenericError](#genericerror) |

##### ***PUT***
**Summary:** Update menu item.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | MenuItem | [MenuItem](#menuitem) |
| 400 | GenericError | [GenericError](#genericerror) |

##### ***DELETE***
**Summary:** Delete menu item.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 |  |  |
| 404 | GenericError | [GenericError](#genericerror) |

### /reservations
---
##### ***POST***
**Summary:** Creates reservation.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Reservation | [Reservation](#reservation) |
| 400 | GenericError | [GenericError](#genericerror) |

### /reservations/approve/{id}
---
##### ***POST***
**Summary:** Approve reservation.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | State | [State](#state) |
| 400 | GenericError | [GenericError](#genericerror) |

### /reservations/cancel/{id}
---
##### ***POST***
**Summary:** Cancel reservation.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | State | [State](#state) |
| 400 | GenericError | [GenericError](#genericerror) |

### /reservations/{id}
---
##### ***GET***
**Summary:** Returns reservation.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Reservation | [ [Reservation](#reservation) ] |
| 500 | Reservation | [ [Reservation](#reservation) ] |

### /tables
---
##### ***GET***
**Summary:** List all tables.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Table | [ [Table](#table) ] |
| 500 | GenericError | [GenericError](#genericerror) |

##### ***POST***
**Summary:** Creates table.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Table | [Table](#table) |
| 400 | GenericError | [GenericError](#genericerror) |

### /tables/{id}
---
##### ***GET***
**Summary:** Returns table.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Table | [Table](#table) |
| 400 | GenericError | [GenericError](#genericerror) |
| 404 | GenericError | [GenericError](#genericerror) |

##### ***PUT***
**Summary:** Updates table.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Table | [Table](#table) |
| 400 | GenericError | [GenericError](#genericerror) |
| 404 | GenericError | [GenericError](#genericerror) |

##### ***DELETE***
**Summary:** Deletes table.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 204 |  |  |
| 400 | GenericError | [GenericError](#genericerror) |

### Models
---

### Duration  

A Duration represents the elapsed time between two instants
as an int64 nanosecond count. The representation limits the
largest representable duration to approximately 290 years.

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| Duration | integer | A Duration represents the elapsed time between two instants
as an int64 nanosecond count. The representation limits the
largest representable duration to approximately 290 years. |  |

### GenericError  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| error | string | Error massage. | No |

### MenuCategory  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | integer (uint64) |  | No |
| name | string | Name of the menu category. | Yes |
| order | integer (uint64) | Order of this category in categories list. | Yes |

### MenuItem  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| active | boolean | Active flag for the menu item. | Yes |
| category_id | integer (uint64) | Category of the menu item. | Yes |
| description | string | Description of the menu item. | Yes |
| id | integer (uint64) |  | No |
| image_url | string | Image URL for the menu item. | Yes |
| name | string | Name of the menu item. | Yes |
| price | float | Price of the menu item in Bahrain Dinars. | Yes |

### Reservation  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| duration | [Duration](#duration) |  | Yes |
| email | string | Email of the client. | Yes |
| full_name | string | Full Name of the client. | Yes |
| guests | long | Number of people to seat for reservation. | Yes |
| id | integer (uint64) |  | No |
| phone | string | Phone of the client. | Yes |
| state | string (state) |  | No |
| table_id | integer (uint64) | ID of table, associated with reservation. | Yes |
| time | dateTime | Time of the reservation. | Yes |

### State  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| State | string |  |  |

### Table  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| active | boolean | Active flag for the table. | Yes |
| description | string | Description of the table. | Yes |
| id | integer (uint64) |  | No |
| places | long | Number of places to seat. | Yes |

### Token  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| Code | string |  | No |
| CreatedAt | dateTime |  | No |
| ID | integer (uint64) |  | No |
| ReservationID | integer (uint64) |  | No |
| Type | [TokenType](#tokentype) |  | No |
| UpdatedAt | dateTime |  | No |
| Used | boolean |  | No |

### TokenType  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| TokenType | string |  |  |