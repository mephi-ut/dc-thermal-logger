
CREATE TYPE history_aggregation_period AS ENUM ('second', 'minute', 'hour', 'day', 'week');

CREATE TABLE history_records (
	id		SERIAL PRIMARY KEY,

	date		timestamp without time zone NOT NULL,

	aggregation_period history_aggregation_period NOT NULL,

	sensor_id	SMALLINT NOT NULL,
	raw_value	REAL NOT NULL,
	converted_value	REAL NOT NULL,
	counter		SMALLINT DEFAULT 1

);

CREATE UNIQUE INDEX history_records_uniidx ON history_records (aggregation_period, date, sensor_id);
CREATE INDEX history_records_date ON history_records (date);
CREATE INDEX history_records_aggregation_period ON history_records (aggregation_period);
CREATE INDEX history_records_timeperiod ON history_records (aggregation_period, date);
CREATE INDEX history_records_converted_value ON history_records (converted_value);

