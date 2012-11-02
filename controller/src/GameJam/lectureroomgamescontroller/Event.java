package GameJam.lectureroomgamescontroller;

import java.net.Socket;

public class Event {
	public Socket client;
	public float joy_x;
	public float joy_y;
	public int buttonCode;
	
	public Event(float joy_x, float joy_y, int buttonCode, Socket client) {
		this.joy_x = joy_x;
		this.joy_y = joy_y;
		this.buttonCode = buttonCode;
		this.client = client;
	}
}