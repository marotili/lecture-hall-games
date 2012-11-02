package GameJam.lectureroomgamescontroller;

import android.os.Bundle;
import android.preference.PreferenceManager;
import android.app.Activity;
import android.content.Intent;
import android.content.SharedPreferences;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;

public class ControllerMainActivity extends Activity {

	Activity activity = this;
	SharedPreferences preferences;
	
    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.controller_main);
        
        Button toControllerButton = (Button) findViewById(R.id.buttonConnect);
        
		preferences = PreferenceManager.getDefaultSharedPreferences(activity);
        String storedServer = preferences.getString("storedServerAddress", "");
        String storedNickname = preferences.getString("nickname", "Anonymous");
		
        EditText editTextServer = (EditText) findViewById(R.id.editTextServer);
        EditText editTextNicknae = (EditText) findViewById(R.id.editTextNickname);
        editTextServer.setText(storedServer);
        editTextNicknae.setText(storedNickname);
       
        
        toControllerButton.setOnClickListener(new View.OnClickListener() {
			
			public void onClick(View v) {
				Intent startActivityIntent = new Intent(ControllerMainActivity.this, controller.class);
				
				EditText editTextServer = (EditText) findViewById(R.id.editTextServer);
				EditText editTextNickname = (EditText) findViewById(R.id.editTextNickname);
					        
		        String serverAddress = editTextServer.getText().toString();
		        String nickname = editTextNickname.getText().toString();

		        SharedPreferences.Editor editor = preferences.edit();
		        editor.putString("storedServerAddress", serverAddress);
		        editor.putString("nickname", nickname);
		        editor.commit();
		       		        
				Bundle bundle = new Bundle();
				bundle.putString("serverAddress", serverAddress);
				bundle.putString("nickname", nickname);
				startActivityIntent.putExtras(bundle);
				startActivity(startActivityIntent);
			}
		});
    }
}
