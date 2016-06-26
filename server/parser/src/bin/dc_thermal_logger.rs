extern crate dc_logger_parser;
extern crate dotenv;
extern crate postgres;

//use self::dc_logger_parser::*;
use self::dc_logger_parser::models::*;


fn parser_listen() -> std::result::Result<std::net::UdpSocket, std::io::Error> {
	let bind_str = std::env::var("DCLOGPARSER_BIND").expect("Environment variable DCLOGPARSER_BIND is required");
	return std::net::UdpSocket::bind(bind_str.as_str());
}

fn pgsql_connect() -> std::result::Result<postgres::Connection, postgres::error::ConnectError> {
	let pgsql_url = std::env::var("DATABASE_URL").expect("DATABASE_URL must be set");
	return postgres::Connection::connect(pgsql_url.as_str(), postgres::SslMode::None);
}

fn main() {
	dotenv::dotenv().ok();

	let pgsql_conn = pgsql_connect().unwrap();
	let socket     = parser_listen().unwrap();

	let mut msg = Message::default();
	let buf: &mut [u8] = unsafe { std::slice::from_raw_parts_mut(&mut msg as *mut Message as *mut u8, std::mem::size_of::<Message>()) };
	loop {
		msg = { Default::default() };
		let got = socket.recv(buf).unwrap();
		if got < 4 {
			println!("Too short message. Skipping.");
			continue;
		}

		if msg.channel_count as usize > MAX_CHANNELS {
			println!("Incorrect message (too many channels: {} > {}). Skipping.", msg.channel_count, MAX_CHANNELS);
			continue;
		}

		for channel_id in 0..(msg.channel_count as i16) {
			let value = msg.channels[channel_id as usize];
			pgsql_conn.execute("INSERT INTO log_records (sensor_id, channel_id, value) VALUES ($1, $2, $3)", &[&msg.sensor_id, &channel_id, &value]).unwrap();
		}
	}
}
