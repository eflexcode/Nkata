package com.ifeanyi.nkataandroid

import android.app.ActivityManager
import android.app.ActivityManager.RunningAppProcessInfo
import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.util.Log
import androidx.core.content.ContextCompat


class MainReceiver : BroadcastReceiver() {

    override fun onReceive(context: Context, intent: Intent) {

        if (isAppInBackground(context)) {
            //BACKGROUND
        } else {
            //FOREGROUND
            if (intent.action == Intent.ACTION_BOOT_COMPLETED) {
                Log.d("BootReceiver", "Device finished booting!")

                val mService = Intent(context, MainService::class.java)
                // Start main service
                ContextCompat.startForegroundService(context, mService)
            }
        }
    }

    private fun isAppInBackground(context: Context?): Boolean {
        val appProcess = RunningAppProcessInfo()
        ActivityManager.getMyMemoryState(appProcess)
        return appProcess.importance != RunningAppProcessInfo.IMPORTANCE_FOREGROUND
    }
}