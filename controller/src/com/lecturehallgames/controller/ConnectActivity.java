package com.lecturehallgames.controller;

import android.app.Activity;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.Bundle;
import android.preference.PreferenceManager;
import android.view.View;
import android.widget.EditText;


public class ConnectActivity extends Activity
{
	public final static String SERVER_ADDRESS = "com.lecturehallgames.controller.SERVER_ADDRESS";
	public final static String NICKNAME = "com.lecturehallgames.controller.NICKNAME";

	private SharedPreferences preferences;
	private EditText          editServer;
	private EditText          editNick;

	/** Called when the activity is first created. */
	@Override
	public void onCreate(Bundle savedInstanceState)
	{
		super.onCreate(savedInstanceState);
		setContentView(R.layout.connect);

		preferences = PreferenceManager.getDefaultSharedPreferences(this);
		editServer = (EditText)findViewById(R.id.edit_server);
		editNick = (EditText)findViewById(R.id.edit_nick);

		editServer.setText(preferences.getString(SERVER_ADDRESS, ""));
		editNick.setText(preferences.getString(NICKNAME, "Anonymous"));
	}

	public void connect(View view)
	{
		String server = editServer.getText().toString();
		String nick = editNick.getText().toString();

		SharedPreferences.Editor editor = preferences.edit();
		editor.putString(SERVER_ADDRESS, server);
		editor.putString(NICKNAME, nick);
		editor.commit();

		Intent intent = new Intent(this, ControlActivity.class);
		intent.putExtra(SERVER_ADDRESS, server);
		intent.putExtra(NICKNAME, nick);
		startActivity(intent);
	}
}
