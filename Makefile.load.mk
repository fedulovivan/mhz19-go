#
# Here are all tasks related to the load testing
# which are extracted into separate file and included into main Makefile.
#

API_LOAD_COUNT ?= 10
API_LOAD_THREADS ?= 10
API_RPS ?= 50
REST_API_DOMAIN ?= 127.0.0.1
REST_API_URL ?= http://$(REST_API_DOMAIN):$(REST_API_PORT)$(REST_API_PATH)

api-load-rules-read:
	oha -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) $(REST_API_URL)/rules

api-load-rule-read:
	oha -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) $(REST_API_URL)/rules/200

api-load-rule-write:
	curl -H 'Content-Type: application/json' -X PUT -d '{"deviceId":"DeviceId(0x00158d0004244bda)","deviceClass":1}' $(REST_API_URL)/devices
	curl -H 'Content-Type: application/json' -X PUT -d '{"deviceId":"DeviceId(10011cec96)","deviceClass":6}' $(REST_API_URL)/devices
	oha --method PUT -H 'Content-Type: application/json' -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) -D ./assets/load/create-rule.json --rand-regex-url $(REST_API_URL)/rules/name-[a-z0-9]{16}

api-load-messages-read:
	oha -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) $(REST_API_URL)/messages/temperature/0x00158d00067cb0c9

api-load-push-message-write-limited:
	oha --method PUT -H 'Content-Type: application/json' -D ./assets/load/push-message.json -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) -q $(API_RPS) $(REST_API_URL)/push-message

api-load-push-message-write-nolimit:
	oha --method PUT -H 'Content-Type: application/json' -D ./assets/load/push-message.json -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) $(REST_API_URL)/push-message

api-load-stats-read:
	oha -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) $(REST_API_URL)/stats

api-load-once:
	wget -O /dev/null $(REST_API_URL)/rules

# oha -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) $(REST_API_URL)/messages/device/0x00158d00067cb0c9?tocsv=1
# oha -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) $(REST_API_URL)/messages/device/0x00158d00067cb0c9
# oha --method PUT -H 'Content-Type: application/json' -D ./assets/load/push-message.json -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) -q $(API_RPS) $(REST_API_URL)/push-message
# ab -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) $(REST_API_URL)/rules
# ab -T application/json -u ./assets/load/create-rule.json -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) $(REST_API_URL)/rules
# oha --method PUT -H 'Content-Type: application/json' -d "{\"name\":\"name-`uuidgen`\"}" -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) $(REST_API_URL)/rules
