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
	private int buttonCode;
	private TextView boundingBoxA;
	private String serverAddress;
	
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
        
        
	    final TextView getCoord = (TextView) findViewById(R.id.textView1);

	    getCoord.setOnTouchListener(new View.OnTouchListener() {
			
			public synchronized boolean onTouch(View v, MotionEvent event) {
				
				Float x = event.getX() / getCoord.getWidth();
				Float y = event.getY() / getCoord.getHeight();
				
				try {
					outStream.writeFloat(x);
					outStream.writeFloat(y);
					outStream.writeInt(0);
					System.out.println("X: "  + x.toString());
					System.out.println("Y: "  + y.toString());
				} catch (IOException e) {

					e.printStackTrace();
				}
				
				
				return true;
			}
		});

	    boundingBoxA = (TextView) findViewById(R.id.boundingBoxA);
	    boundingBoxA.setOnTouchListener(new View.OnTouchListener() {
			
			public boolean onTouch(View v, MotionEvent event) {
				TextView middle = (TextView) findViewById(R.id.textView3);
				middle.setText("Button A");
				return true;
			}
		});
	    
	    TextView boundingBoxB = (TextView) findViewById(R.id.boundingBoxB);
	    boundingBoxB.setOnTouchListener(new View.OnTouchListener() {
			
			public boolean onTouch(View v, MotionEvent event) {
				TextView middle = (TextView) findViewById(R.id.textView3);
				middle.setText("Button B");
				return true;
			}
		});
	}
	

}
