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
        }
        			
        try {
        	client = new Socket(serverAddress, 8001);
        	outStream = new DataOutputStream(client.getOutputStream());
        	
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
				
				try {
					if(x < parentView.getWidth()/2)
					{
					outStream.writeFloat((float) ( x/ (0.5 * parentView.getWidth())));
					outStream.writeFloat(y/ parentView.getHeight());
					outStream.writeInt(0);
					
					if(event.getAction() == MotionEvent.ACTION_UP) {
						outStream.writeFloat(0.5f);
						outStream.writeFloat(0.5f);
						outStream.writeInt(0);
					}
					
					}
				} catch (IOException e) {
					e.printStackTrace();
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
}
