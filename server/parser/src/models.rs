
extern crate time;
use self::time::Timespec;

pub struct RawRecord {
	pub date:		Timespec,
	pub raw_sensor_id:	i16,
	pub raw_channel_id:	i16,
	pub raw_value:		i16,
}

pub const MAX_CHANNELS: usize = 16;

pub struct Message {
	pub sensor_id:		 i16,
	pub channel_count:	 i16,
	pub channels:		[i16; MAX_CHANNELS]
}

impl Default for Message {
	#[inline]
	fn default() -> Message {
		Message { sensor_id: 0, channel_count: 0, channels: [0; MAX_CHANNELS] }
	}
}
