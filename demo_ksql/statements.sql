CREATE STREAM ST_TRACKING WITH (KAFKA_TOPIC='public.cc.sr.pb.demo',VALUE_FORMAT='PROTOBUF');
SELECT * FROM ST_TRACKING LIMIT 2;


-- {
--   "USER_ID": "4aa65dfe-b89c-4b1c-813a-409cdab05d16",
--   "FIRST_NAME": "Cool",
--   "LAST_NAME": "Guy",
--   "MESSAGE": "A normal message.",
--   "TIMESTAMP": "2022-04-23T03:38:54.114",
--   "ITEMS": [
--     {
--       "NAME": "Item-0",
--       "VALUE": "3274"
--     },
--     {
--       "NAME": "Item-1",
--       "VALUE": "1211"
--     },
--     {
--       "NAME": "Item-2",
--       "VALUE": "1445"
--     }
--   ],
--   "ADDRESS": {
--     "STREET": "Kafka Street",
--     "POSTALCODE": "000-0000",
--     "CITY": "Tokyo",
--     "COUNTRY": "Japan"
--   }
-- }


-- filter stream
CREATE STREAM ST_TRACKING_JAPAN_USER AS
SELECT * FROM ST_TRACKING WHERE ADDRESS->COUNTRY = 'Japan';


-- do aggregation, write to a table
CREATE TABLE TB_TRACKING_CITY_ITEMS  AS
SELECT
    from_unixtime(WINDOWSTART) as WINDOW_START,
    ADDRESS->CITY AS CITY,
    SUM(ARRAY_LENGTH(ITEMS)) TOTAL_ITEM_COUNT
FROM ST_TRACKING_JAPAN_USER
    WINDOW TUMBLING (SIZE 1 MINUTE)
GROUP BY ADDRESS->CITY;

