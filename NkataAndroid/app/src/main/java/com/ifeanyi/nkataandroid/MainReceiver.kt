package com.ifeanyi.nkataandroid

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.util.Log
import androidx.core.content.ContextCompat

class MainReceiver : BroadcastReceiver() {

    override fun onReceive(context: Context, intent: Intent) {
        if (intent.action == Intent.ACTION_BOOT_COMPLETED) {
            Log.d("BootReceiver", "Device finished booting!")

            val mService = Intent(context, MainService::class.java)

            ContextCompat.startForegroundService(context,mService)
            // Start main service
        }
    }
}