Tasks

User

POST:

User Creation

curl -i -d "\$BODY" localhost:4000/v1/users

BODY='{"username":"lopezvictor","password":"belize12345","fullname":"Victor Lopez", "email":"lopezvic@example.com", "phone":"501-607-2462", "address":"George St San Ignacio","district_id":1, "user_type_id":1}'

Upload user profile image

PUT:

    BODY='{"username":"lopezvictor","fullname":"Victor J Lopez", "email":"lopezvic@example.com", "phone":"501-607-2463", "address":"George Second St San Ignacio","district_id":"Cayo", "user_type_id":"agent"}'

    curl -i -X PUT -d "$BODY" localhost:4000/v1/users/updated/:id

Activation

curl -X PUT -d '{"token": "BVL2MDEWKWBHP7PR5G2LMTJ7PI"}' localhost:4000/v1/users/activated

Authentication

BODY='{"username":"lopezvictor","password":"belize12345"}'

curl -i -X POST -d "\$BODY" localhost:4000/v1/tokens/authentication

Listing

Create Listing

GET listing by id

POST

      BODY='{"property_title": "Land for Sale in the Area of San Ignacio ", "property_status_id":1, "property_type_id":1,"price":70000, "description":"Land Size is 200 ft by 200 ft", "address":"27 Street, San Ignacio Town", "district_id": 1, "google_map_url": "google.com/3wdfdyf9"}'

      curl -i -d "$BODY" localhost:4000/v1/listings

PUT

BODY='{"property_title": "Land for Sale in San Ignacio Town", "property_status_id":"Available", "property_type_id":"Land","price":80000, "description":"Land Size is 100 ft by 200 ft", "address":"27 Street, San Ignacio Town", "district_id": "Cayo", "google_map_url": "google.com/3wdfdsf9"}'

      curl -X PUT -d "$BODY" localhost:4000/v1/listings/update/id

//Task to Complete

add user agent to listing/get - done
add listing image - done
update property status - leased, sold - done
views for top agents - done
views for report - done
func(w http.ResponseWriter, r *http.Request)

//ADD USER PROPERTIES

BODY='{"username": "trumpvictor", "listing_id": 3}'

curl -i -d "$BODY" localhost:4000/v1/users/listings


//Update property status -
BODY='{ "property_status_id":"Sold"}'
curl -X PUT -d "$BODY" localhost:4000/v1/listings/update/3