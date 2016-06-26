
CREATE TABLE log_records (
	date		timestamp without time zone default (now() at time zone 'utc'),
	sensor_id	SMALLINT NOT NULL,
	channel_id	SMALLINT NOT NULL,
	value		SMALLINT NOT NULL,
	converted_value	SMALLINT,

	PRIMARY KEY(date, sensor_id, channel_id)
);

CREATE INDEX log_records_date ON log_records (date);

