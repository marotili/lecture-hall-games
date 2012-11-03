package com.lecturehallgames.controller;

import com.lecturehallgames.controller.ConnectActivity;

import android.app.Activity;
import android.os.AsyncTask;
import android.os.Vibrator;
import android.app.AlertDialog;
import android.os.Bundle;
import android.content.Context;
import android.content.DialogInterface;
import android.content.Intent;
import android.view.MotionEvent;
import android.view.View;
import android.view.View.OnTouchListener;
import android.widget.ImageView;

import java.io.BufferedOutputStream;
import java.io.DataOutputStream;
import java.io.DataInputStream;
import java.io.IOException;
import java.net.Socket;
import java.net.UnknownHostException;


public class ControlActivity extends Activity implements OnTouchListener
{
	private ImageView imageController;
	private TransmitInputTask transmitter;
	private String address, nickname;
	protected Activity activity = this;

    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.control);

        imageController = (ImageView)findViewById(R.id.image_controller);
		imageController.setOnTouchListener(this);

		Intent intent = getIntent();
		address = intent.getStringExtra(ConnectActivity.SERVER_ADDRESS);
		nickname = intent.getStringExtra(ConnectActivity.NICKNAME);
    }

    @Override
    public void onPause() {
    	super.onPause();
    	if (transmitter != null) {
    		transmitter.cancel(true);
    	}
    }

    @Override
    public void onResume() {
    	super.onResume();
    	transmitter = new TransmitInputTask(address, nickname);
		transmitter.execute();
    }

    public boolean onTouch(View view, MotionEvent event) {
    	float joystickX = 0.0f;
		float joystickY = 0.0f;
		boolean buttonA = false;
		boolean buttonB = false;

		int action = event.getAction() & MotionEvent.ACTION_MASK;
		int pointer = (event.getAction() & MotionEvent.ACTION_POINTER_INDEX_MASK) >> MotionEvent.ACTION_POINTER_INDEX_SHIFT;

		if (action != MotionEvent.ACTION_UP) {
			for(int i = 0 ; i < event.getPointerCount() ; i++) {
				if(action != MotionEvent.ACTION_POINTER_UP || pointer != i) {
					float x = event.getX(i) / view.getWidth();
					float y = event.getY(i) / view.getHeight();

					if(x < 0.5 ) {
						joystickX = x * 4.0f - 1.0f;
						joystickY = x * 2.0f - 1.0f;
					} else if(x < 0.75 && y > 0.5) {
						buttonB = true;
					} else if(y > 0.5) {
						buttonA = true;
					}
				}
			}
		}
		transmitter.setState(joystickX, joystickY, buttonA, buttonB);
		return true;
	}

	private class TransmitInputTask extends AsyncTask<Void, Void, Void>
	{
		public final String address;
		public final String nickname;
		private String errorMsg;

		private float joystickX, joystickY;
		private boolean buttonA, buttonB;
		private boolean idle;

		public TransmitInputTask(String address, String nickname) {
			super();
			this.address = address;
			this.nickname = nickname;
		}

		public synchronized void setState(float joystickX, float joystickY, boolean buttonA, boolean buttonB) {
			this.joystickX = joystickX;
			this.joystickY = joystickY;
			this.buttonA = buttonA;
			this.buttonB = buttonB;
		}

		@Override
		protected Void doInBackground(Void... params) {
			Socket socket = null;
			VibrateTask vibrator = null;
			errorMsg = null;
			try {
				socket = new Socket(address, 8001);
				BufferedOutputStream buffer = new BufferedOutputStream(socket.getOutputStream());
				DataOutputStream output = new DataOutputStream(buffer);

				output.writeInt(nickname.length());
				output.write(nickname.getBytes());
				buffer.flush();

				vibrator = new VibrateTask(socket);
				vibrator.execute();

				while (!isCancelled() && socket.isConnected()) {
					float joystickX, joystickY;
					boolean buttonA, buttonB;
					synchronized (this) {
						joystickX = this.joystickX;
						joystickY = this.joystickY;
						buttonA = this.buttonA;
						buttonB = this.buttonB;
					}

					int buttonCode = 0;
					if (buttonA) {
						buttonCode |= 1;
					}
					if (buttonB) {
						buttonCode |= 2;
					}

					if (joystickX != 0 || joystickY != 0 || buttonCode != 0) {
						idle = false;
					}
					if (!idle) {
						try {
							output.writeFloat(joystickX);
							output.writeFloat(joystickY);
							output.writeInt(buttonCode);
							buffer.flush();
						} catch (IOException exception) {
						}
					}
					if (joystickX == 0 && joystickY == 0 && buttonCode == 0) {
						idle = true;
					}

					try {
						Thread.sleep(20);
					} catch(InterruptedException e) {
					}
				}
			} catch(UnknownHostException exception) {
				errorMsg = "Unknown Host: " + address;
			} catch(IOException exception) {
				errorMsg = "Failed to connect to \"" + address + "\".";
			} finally {
				if (vibrator != null) {
					vibrator.cancel(true);
				}
				if (socket != null) {
					try {
						socket.close();
					} catch (IOException exception) {
					}
				}
			}
			return null;
		}

		@Override
		protected void onPostExecute(Void result) {
			if (errorMsg != "") {
				AlertDialog ad = new AlertDialog.Builder(activity).create();
				ad.setCancelable(false);
				ad.setMessage(errorMsg);
				ad.setButton("OK", new DialogInterface.OnClickListener() {
					@Override
					public void onClick(DialogInterface dialog, int which) {
						dialog.dismiss();
					}
				});
				ad.show();
			}
		}
	}

	private class VibrateTask extends AsyncTask<Void, Integer, Void>
	{
		Socket socket = null;
		Vibrator vibrator = null;

		public VibrateTask(Socket socket) {
			this.socket = socket;
			vibrator = (Vibrator)getSystemService(Context.VIBRATOR_SERVICE);
			if (vibrator != null) {
				vibrator.vibrate(300);
			}
		}

		@Override
		protected Void doInBackground(Void... params) {
			try {
				DataInputStream input = new DataInputStream(socket.getInputStream());
				while (!isCancelled() && socket.isConnected()) {
					try {
						int action = input.readInt();
						publishProgress(action);
						if (vibrator != null) {
							vibrator.vibrate(300);
						}
					} catch (IOException exception) {
					}
				}
			} catch (IOException exception) {
			}
			return null;
		}
	}
}
