
CREATE TABLE raw_records (
	date		timestamp without time zone default (now() at time zone 'utc'),

	sensor_id	SMALLINT NOT NULL,
	channel_id	SMALLINT NOT NULL,
	raw_value	SMALLINT NOT NULL,

	PRIMARY KEY(date, sensor_id, channel_id)
);

CREATE INDEX raw_records_date ON raw_records (date);

