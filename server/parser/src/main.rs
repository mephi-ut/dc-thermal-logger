
extern crate sqlite;
extern crate time;

use std::net::UdpSocket;
use std::mem;
use std::slice;

use time::Timespec;

const MAX_CHANNELS: usize = 16;

struct Message {
	sensor_id:	 u16,
	channel_count:	 u16,
	channels:	[u16; MAX_CHANNELS]
}

impl Default for Message {
	#[inline]
	fn default() -> Message {
		Message { sensor_id: 0, channel_count: 0, channels: [0; MAX_CHANNELS] }
	}
}

struct HistoryRow {
	id:		i64,
	doc:		Timespec,
	sensor_id:	u16,
	channel_id:	u16,
	value:		u16,
}

fn main() {
	let socket = UdpSocket::bind("127.0.0.1:8888").unwrap();
	let mut msg: Message = { Default::default() };

	let buf: &mut [u8] = unsafe { slice::from_raw_parts_mut(&mut msg as *mut Message as *mut u8, mem::size_of::<Message>()) };

	loop {
		msg = { Default::default() };
		let got = socket.recv(buf).unwrap();
		if got < 4 {
			println!("Too short message. Skipping.");
			continue;
		}
	}
}
