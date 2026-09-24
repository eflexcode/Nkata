package com.ifeanyi.nkataandroid

import android.app.Application
import android.app.NotificationChannel
import android.app.NotificationManager
import android.os.Build
import com.ifeanyi.nkataandroid.logic.Util

class MainApplication : Application() {
    override fun onCreate() {
        super.onCreate()
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val notificationChannel = NotificationChannel(
                Util.CHANNEL_ID,
                Util.CHANNEL_NAME,
                NotificationManager.IMPORTANCE_LOW
            )

            val manager =  getSystemService(NotificationManager::class.java)
            manager.createNotificationChannel(notificationChannel)

        }
    }
}