# GO BOT DOCUMENTATION

## GOAL

Provide tool to conveniently save information on a fly and then securely be able to access this information from everywhere

## TECHs

 - AWS CDK
 - Golang

### Golang
 - Main handler () => Response
    - Parses parameters
    - Validates them
    - Calls *logic handler*

    - Logic handler (params) => Error/Data
        - Checks status of user *Function* 
            -is active/there
            -what's current status
        - Call specific *Function* to perform action on DB [GET/PUT/UPDATE/DELETE]

        - Function (method, data) => error/data
            - calls specific function based on method
        
                - GET (id) => error/data
                - UPDATE (id, new data) => error/data


    




### Ex
#### Film
    - message: "fa"
    - message: "Home Alone"
Result: saved film 

    - message: "fl"
Result: show film list

    - message: "fd"
    - message: "2"
Result: remove film on index 2

    - message: "fu"
    - message: "2"
    - message: "Home Alone 2"
Result: update saved film
