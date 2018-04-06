# GoQL supported functions 

## Structure functions 

| Name                                        | Description                                                                                                                                                                                                                                                       |
| ----                                        | -----------                                                                                                                                                                                                                                                       |
| is_struct(definition)                       | return true if the definition is struct, false if not and null if the parameter is not definition and is not cast-able to definition                                                                                                                              |
| field_def(definition, int/string)           | return definition of field at the position or with name of the second argument. if the first argument is not valid definition of a struct or the index is out of range (its from 1) or the name is not a field, t return null                                     |
| field_name(definition, int)                 | return the field name at the index (from 1), return null if the definition is not valid of a struct or the index is out of range                                                                                                                                  |
| field_count(definition)                     | return the number of fields in definition or null if its not valid definition or its not a struct                                                                                                                                                                 |
| filed_tag(definition, int/string, [string]) | return the struct tag of the field at the position or with name of the second argument add if the last argument is available, return the tag for that specific value, return null if definition is not a valid struct or the field is out of index or not correct |
| embed_def(definition, int)                  | return the embedded structure definition at the position, or null if the definition is not valid or the position is out of range                                                                                                                                  |
| embed_count(definition)                     | return the number of embedded item in struct, null if the definition is not a valid struct                                                                                                                                                                        |
| embed_tag(definition, int, [string])        | return the struct tag of the embedded item at the position, if the last arguments is available return specific tag, return null if the definition is not struct or index is out of range                                                                          |


## Map functions

| Name                | Description                                                                                                                         |
| ----                | -----------                                                                                                                         |
| is_map(definition)  | return true if the definition is a map, false if not and null if the parameter is not definition and is not cast-able to definition |
| map_key(definition) | return definition of the map key, and null if the definition is not a map or is not valid                                           |
| map_val(definition) | return definition of the map value, and null if the definition is not a map or is not valid                                         |


## Interface functions 

| Name                     | Description                                                                                                  |
| ----                     | -----------                                                                                                  |
| is_interface(definition) | return true if the definition is an interface or false if not. if the definition is not valid it return null |
