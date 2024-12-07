#@get_postal@
POST https://uniqlo-customerapi.brainergy.io/reqft/get_postal
{
    "subdistrict_id": "{{subdistrict_id}}"
}
HTTP 200
[Asserts]
header "Content-Type" contains  "application/json"
jsonpath  "$.status" == 200
jsonpath  "$.msg" == "SUCCESS"
jsonpath  "$.data.postal_code" != null
[Captures]
postal_code : jsonpath "$.data.postal_code"


##
#@get_postal_fail@
POST https://uniqlo-customerapi.brainergy.io/reqft/get_postal
{
    "subdistrict_id": "0"
}
HTTP 200
[Asserts]
header "Content-Type" contains  "application/json"
jsonpath  "$.status_code" == 400

#
https://uniqlo-customerapi.brainergy.io/reqft/get_postal post
{
    "subdistrict_id": "0"
}