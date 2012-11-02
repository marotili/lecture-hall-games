package GameJam.lectureroomgamescontroller;

import java.io.DataOutputStream;
import java.io.IOException;
import java.net.Socket;
import java.net.UnknownHostException;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.graphics.Color;
import android.os.Bundle;
import android.os.StrictMode;
import android.util.Log;
import android.view.MotionEvent;
import android.view.TextureView;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;

@SuppressLint("NewApi")
public class controller extends Activity {
	private Socket client;
	private DataOutputStream outStream;
	private String serverAddress;
	private String nickname;
	View parentView;
	
	@SuppressLint({ "NewApi", "NewApi" })
	@Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.controller);
        
        StrictMode.ThreadPolicy policy = new StrictMode.ThreadPolicy.Builder().permitAll().build();
        StrictMode.setThreadPolicy(policy); 
        
        Bundle bundle = getIntent().getExtras();
       	if (bundle != null) {
       		serverAddress = bundle.getString("serverAddress");
       		nickname = bundle.getString("nickname");
        }
        			
        
        connectoToServer();

        parentView = findViewById(R.id.entire_view);
       

	    parentView.setOnTouchListener(new View.OnTouchListener() {
			
			public synchronized boolean onTouch(View v, MotionEvent event) {							
				Float joystickX = 0.0f;
				Float joystickY = 0.0f;
				boolean buttonA = false;
				boolean buttonB = false;
				
				int action = event.getAction() & MotionEvent.ACTION_MASK;
				
				int pointer = (event.getAction() & MotionEvent.ACTION_POINTER_INDEX_MASK) >> MotionEvent.ACTION_POINTER_INDEX_SHIFT;
				
				
				if(action != MotionEvent.ACTION_UP) {
					for(int i = 0 ; i < event.getPointerCount() ; i++) {
						if(action == MotionEvent.ACTION_POINTER_UP && pointer == i)
						{
							continue;
						}
						float currentX = event.getX(i) / parentView.getWidth();
						float currentY = event.getY(i) / parentView.getHeight();
						
						if(currentX < 0.5 ) {
							joystickX = currentX * 4.0f - 1.0f;
							joystickY = currentY * 2.0f - 1.0f;
						}
						else if(currentX < 0.75 && currentY > 0.5) {
							buttonB = true;
						}
						else if(currentY > 0.5) {
							buttonA = true;
						}
						
					} 
				}
				//System.out.println("Input: " + joystickX + " " + joystickY + " " + buttonB + " " + buttonA);
				writeToServer(joystickX,joystickY,buttonA,buttonB);
				return true;
			}
		});
	}
	
	@Override
	public void onPause() {
		super.onPause();		
		try {
			client.close();
		} catch (IOException e) {
			e.printStackTrace();
		}
	}
	
	private void writeToServer(float joystickX, float joystickY, boolean buttonA, boolean buttonB) {
		try {
			
			int buttonCode = 0;
			if(buttonA) {
				buttonCode |= 1;
			}
			if(buttonB) {
				buttonCode |= 2;
			}
			
			
			outStream.writeFloat(joystickX);
			outStream.writeFloat(joystickY);
			outStream.writeInt(buttonCode);
		} catch (IOException e) {
			e.printStackTrace();
		}
	}
	
	private void connectoToServer() {
		try {
	    	client = new Socket(serverAddress, 8001);
	    	outStream = new DataOutputStream(client.getOutputStream());
	    	
	    	int nicknameLength = nickname.length();
	    	
	    	outStream.writeInt(nicknameLength);
	    	outStream.write(nickname.getBytes());
        } catch(UnknownHostException exception) {
        	exception.printStackTrace();
        } catch(IOException exception) {
        	exception.printStackTrace();
        }
	}
}
