
CREATE TABLE raw_records (
	id		SERIAL PRIMARY KEY,

	date		timestamp without time zone default (now() at time zone 'utc'),

	sensor_id	SMALLINT NOT NULL,
	channel_id	SMALLINT NOT NULL,
	raw_value	SMALLINT NOT NULL
);

CREATE INDEX raw_records_date ON raw_records (date);
CREATE UNIQUE INDEX raw_records_uniidx ON raw_records (date, sensor_id, channel_id);

