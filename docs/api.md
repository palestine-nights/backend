### /menu
---
##### ***GET***
**Summary:** List all menu items.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | MenuItem | [ [MenuItem](#menuitem) ] |

##### ***POST***
**Summary:** Create menu item.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | MenuItem | [MenuItem](#menuitem) |

### /menu/categories
---
##### ***GET***
**Summary:** List menu categories.

### /menu/{category}
---
##### ***GET***
**Summary:** List menu items with specified category.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | MenuItem | [ [MenuItem](#menuitem) ] |

### /menu/{id}
---
##### ***GET***
**Summary:** Returns menu item.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | MenuItem | [MenuItem](#menuitem) |

##### ***PUT***
**Summary:** Update menu item.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | MenuItem | [MenuItem](#menuitem) |

##### ***DELETE***
**Summary:** Delete menu item.

**Responses**

| Code | Description |
| ---- | ----------- |
| 204 |  |

### /reservations
---
##### ***POST***
**Summary:** Creates reservation.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Reservation | [Reservation](#reservation) |

### /reservations/{id}
---
##### ***GET***
**Summary:** Returns reservation.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Reservation | [ [Reservation](#reservation) ] |

##### ***POST***
**Summary:** Cancel reservation.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | State | [State](#state) |

### /tables
---
##### ***GET***
**Summary:** List all tables.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Table | [ [Table](#table) ] |

##### ***POST***
**Summary:** Creates table.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Table | [Table](#table) |

### /tables/{id}
---
##### ***GET***
**Summary:** Returns table.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Table | [Table](#table) |

##### ***PUT***
**Summary:** Updates table.

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Table | [Table](#table) |

##### ***DELETE***
**Summary:** Deletes table.

**Responses**

| Code | Description |
| ---- | ----------- |
| 204 |  |

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

### MenuItem  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| category | string | Category of the menu item. | Yes |
| description | string | Description of the menu item. | Yes |
| id | integer (uint64) |  | No |
| image_url | string | Image URL for the menu item. | Yes |
| name | string | Name of the menu item. | Yes |
| price | float | Price of the menu item in Bahrain Dinars. | Yes |

### Reservation  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| duration | [Duration](#duration) |  | No |
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
| description | string | Description of the table. | Yes |
| id | integer (uint64) |  | No |
| places | long | Number of places to seat. | Yes |