
CREATE TYPE history_aggregation_period AS ENUM ('minute', 'hour', 'day', 'week');

CREATE TABLE history_records (
	date		timestamp without time zone default (now() at time zone 'utc'),

	aggregation_period history_aggregation_period NOT NULL,

	sensor_id	SMALLINT NOT NULL,
	channel_id	SMALLINT NOT NULL,
	raw_value	SMALLINT NOT NULL,
	converted_value	SMALLINT NOT NULL,

	PRIMARY KEY(aggregation_period, date, sensor_id, channel_id)
);

CREATE INDEX history_records_date ON history_records (date);
CREATE INDEX history_records_aggregation_period ON history_records (aggregation_period);
CREATE INDEX history_records_timeperiod ON history_records (aggregation_period, date);
CREATE INDEX history_records_converted_value ON history_records (converted_value);

