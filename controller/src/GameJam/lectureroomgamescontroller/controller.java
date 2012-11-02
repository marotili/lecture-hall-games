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
        
        
        parentView = findViewById(R.id.entire_view);
       

	    parentView.setOnTouchListener(new View.OnTouchListener() {
			
			public synchronized boolean onTouch(View v, MotionEvent event) {			
				Float x = event.getX();
				Float y = event.getY();
				
				if(event.getAction() == MotionEvent.ACTION_MOVE) {
					int historySize = event.getHistorySize();
					
					for(int i = 0 ; i < historySize ; i++) {
						float hx = event.getHistoricalX(i);
						if(hx >  (0.5 * parentView.getWidth())) {
							System.out.println("2nd Half");
						}
					}
				}
				
				if(x < parentView.getWidth()/2)
				{
				writeToServer((float) ( x / (0.5 * parentView.getWidth())), y / parentView.getHeight(), 0);
					
				if(event.getAction() == MotionEvent.ACTION_UP) {
					writeToServer(0.0f, 0.0f, 0);
					}
				}

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
	
	private void writeToServer(float x_value, float y_value, int buttonCode) {
		try {
			outStream.writeFloat(x_value*2.0f-1.0f);
			outStream.writeFloat(y_value*2.0f-1.0f);
			outStream.writeInt(buttonCode);
		} catch (IOException e) {
			e.printStackTrace();
		}
	}
	
}
