# Tables

## files table 

list of files inside package 

| Field name | Type   | Nullable | Description                                          | Plugin |
| -----      | ---    | -----    | ------                                               | ----   |
| name       | String | false    | the file name                                        | No     |
| pkg_name   | String | false    | the package name (the current package)               | No     |
| pkg_path   | String | false    | the package path (like github.com/fzerorubigd/onion) | No     |
| docs       | String | true     | documents at the beginning of the file               | No     |
| goimport   | Bool   | false    | is this file is go imported                          | Yes    |

## imports table 

list of imported package inside a package

| Field name | Type   | Nullable | Description                                          | Plugin |
| ---------- | ----   | -------- | -----------                                          | ------ |
| name       | String | False    | the file name with this import is defined            | No     |
| pkg_name   | String | false    | the package name (the current package)               | No     |
| pkg_path   | String | false    | the package path (like github.com/fzerorubigd/onion) | No     |
| docs       | String | true     | document of this import if any                       | No     |
| canonical  | String | true     | the canonical name of this import (if available)     | No     |
| path       | String | false    | the imported package path                            | No     |
| package    | String | false    | the imported package name                            | No     |

## funcs table

list of function inside the package

| Field name       | Type       | Nullable | Description                                              | Plugin |
| ----------       | ----       | ---      | ----------                                               | ---    |
| name             | String     | false    | the function name                                        | No     |
| pkg_name         | String     | false    | the package name (normally the current package name)     | No     |
| pkg_path         | String     | false    | the package path                                         | No     |
| file             | String     | false    | the file of this function, where the function is defined | No     |
| receiver         | String     | true     | the receiver, if its not a method, its null              | No     |
| pointer_receiver | Bool       | true     | is this receiver is pointer, null if it is not method    | No     |
| exported         | Bool       | false    | if the function name is exported (started with A to Z)   | No     |
| docs             | String     | true     | documents of this function                               | No     |
| body             | String     | false    | body of function (experimental)                          | No     |
| def              | Definition | false    | the type definition, special type                        | No     |

## vars table

list of variables inside package (only global one)

| Field name | Type       | Nullable | Description                                              | Plugin |
| ---------- | ----       | ---      | ----------                                               | ---    |
| name       | String     | false    | the variable name                                        | No     |
| pkg_name   | String     | false    | the package name (normally the current package name)     | No     |
| pkg_path   | String     | false    | the package path                                         | No     |
| file       | String     | false    | the file of this variable, where the variable is defined | No     |
| exported   | Bool       | false    | if the variable name is exported (started with A to Z)   | No     |
| docs       | String     | true     | documents of this variable                               | No     |
| def        | Definition | false    | the type definition, special type                        | No     |


## types table

list of global types defined in the package

| Field name | Type       | Nullable | Description                                          | Plugin |
| ---------- | ----       | ---      | ----------                                           | ---    |
| name       | String     | false    | the type name                                        | No     |
| pkg_name   | String     | false    | the package name (normally the current package name) | No     |
| pkg_path   | String     | false    | the package path                                     | No     |
| file       | String     | false    | the file of this type, where the type is defined     | No     |
| exported   | Bool       | false    | if the type name is exported (started with A to Z)   | No     |
| docs       | String     | true     | documents of this type                               | No     |
| def        | Definition | false    | the type definition, special type                    | No     |


## consts table

list of global constant inside package 

| Field name | Type       | Nullable | Description                                              | Plugin |
| ---------- | ----       | ---      | ----------                                               | ---    |
| name       | String     | false    | the constant name                                        | No     |
| pkg_name   | String     | false    | the package name (normally the current package name)     | No     |
| pkg_path   | String     | false    | the package path                                         | No     |
| file       | String     | false    | the file of this constant, where the constant is defined | No     |
| exported   | Bool       | false    | if the constant name is exported (started with A to Z)   | No     |
| docs       | String     | true     | documents of this constant                               | No     |
| def        | Definition | false    | the type definition, special type                        | No     |


